// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package ui

import (
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"digitdisk/internal/core"
	"digitdisk/internal/scan"
)

// This file builds the picture the walk screen draws: what the walk has passed
// so far, and the shape of the tree it leaves behind.  It runs inside the walk
// itself — scan.Options.Watch is called on every one of several million
// entries — so the whole of it is written for the cost of one entry:
//
//   - Nothing is looked up per entry.  Entries arrive grouped by the directory
//     they sit in, so the node they belong to is remembered and reused; the
//     lookup happens once per directory instead of once per entry.
//   - Nothing is summed up the tree per entry.  A directory keeps only what
//     lies DIRECTLY in it, and the subtree totals are added up once, at the
//     end, in a single pass over the directories.  Adding a file's bytes to
//     each of its ancestors would be depth times the work, per entry.
//   - Nothing is locked per entry.  The walking goroutine owns everything
//     here; the screen asks for a snapshot by raising a flag, and the walk
//     hands one over at the next thousandth entry.  The screen never reads
//     these fields.
//
// The one exception to "per entry" is the top-level list — the children of the
// root, which is what the screen shows growing — and it is two pointer
// additions, because the child of the root an entry belongs to is remembered
// alongside its directory.

// maxNodes caps the directories kept for walking around afterwards.
//
// WHAT THE TREE COSTS, measured rather than guessed.  A walk of /srv here —
// 5 446 842 entries, 574 005 directories, 434.8 GiB — finishes with a live
// heap of 186 MB and a peak RSS of 291 MB, against 31 MB live and 22 MB RSS
// for the same walk printing its report and keeping no tree.  That is about
// 320 bytes per directory: the node itself, its children, and — the half that
// goes away when the walk ends — the lookup table of path to node, which
// settle drops.  `analyze --plain` and `analyze --json` build no tree at all,
// so a script keeps the footprint it always had.
//
// A million directories is therefore the point past which a tool that measures
// disks would become a memory problem itself.  Past the cap the walk still
// counts everything — the totals stay right — but new directories are folded
// into the nearest kept ancestor and the screen says so out loud.
const maxNodes = 1 << 20

// wnode is one directory of the walked tree.
//
// name is the last component, copied rather than sliced out of the full path:
// a Go substring holds its whole original alive, and half a million directories
// each pinning a hundred-byte path is tens of megabytes for nothing.
// own/ownEntries are what lies directly in this directory; bytes/entries are the
// whole subtree and are filled in by settle after the walk.
type wnode struct {
	name   string
	parent *wnode
	kids   []*wnode
	depth  int

	own        int64 // bytes charged to entries directly inside
	ownEntries int32
	ownFiles   int32
	topName    string // the largest single file directly inside
	topBytes   int64

	bytes   int64 // the whole subtree, after settle
	entries int32

	live       int64 // running subtree total, kept only for children of the root
	liveEntrie int32

	sorted bool
}

// path spells the node out from the root down.
func (n *wnode) path() string {
	if n.parent == nil {
		return n.name
	}
	var parts []string
	for c := n; c != nil; c = c.parent {
		parts = append(parts, c.name)
	}
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	head := parts[0]
	if strings.HasSuffix(head, "/") {
		return head + strings.Join(parts[1:], "/")
	}
	return head + "/" + strings.Join(parts[1:], "/")
}

// children returns the subdirectories largest first.  The order is settled
// once, when a directory is first opened, and not on every redraw.
func (n *wnode) children() []*wnode {
	if !n.sorted {
		sort.Slice(n.kids, func(i, j int) bool {
			if n.kids[i].bytes != n.kids[j].bytes {
				return n.kids[i].bytes > n.kids[j].bytes
			}
			return n.kids[i].name < n.kids[j].name
		})
		n.sorted = true
	}
	return n.kids
}

// topRow is one line of the "what is filling up" list.
type topRow struct {
	Name    string
	Bytes   int64
	Entries int32
	// Own marks the row that stands for the files lying directly in the
	// root rather than for a subdirectory.  It has no name of its own: the
	// screen writes one in the language it is drawn in.
	Own bool
}

// walkSnap is what the screen draws: the counters as they stood at one moment
// of the walk.  It is a copy — the screen never touches the collector.
type walkSnap struct {
	At        time.Time
	Entries   int64
	Files     int64
	Dirs      int64
	Links     int64
	Bytes     int64
	Cur       string
	Depth     int
	Tops      []topRow
	Truncated bool
}

// walkFeed collects the picture.  Everything above `want` belongs to the
// walking goroutine alone.
type walkFeed struct {
	root string

	entries, files, dirs, links int64
	bytes                       int64

	nodes map[string]*wnode
	tree  *wnode

	cur    *wnode // the directory whose entries are arriving now
	top    *wnode // the child of the root that directory sits under
	prefix string // that directory's path with a separator, for the cheap test
	depth  int

	count     int
	truncated bool
	tick      uint64

	want atomic.Bool
	out  chan walkSnap
}

func newWalkFeed(root string) *walkFeed {
	f := &walkFeed{
		root:  root,
		nodes: make(map[string]*wnode, 1024),
		out:   make(chan walkSnap, 1),
	}
	f.tree = &wnode{name: root}
	f.nodes[root] = f.tree
	f.count = 1
	f.cur = f.tree
	f.prefix = withSlash(root)
	f.depth = 0
	return f
}

func withSlash(p string) string {
	if strings.HasSuffix(p, "/") {
		return p
	}
	return p + "/"
}

// step is the hot path: one accounted entry.
func (f *walkFeed) step(st scan.Step) {
	p := st.Path
	n := len(f.prefix)
	// The entry belongs to the remembered directory when it starts with that
	// directory's path and has no separator left after it.  Both tests run
	// over a handful of bytes; the alternative — asking a map for the parent
	// of every entry — is a hash of the whole path several million times.
	if len(p) <= n || p[:n] != f.prefix || strings.IndexByte(p[n:], '/') >= 0 {
		f.turn(p)
	}

	f.entries++
	f.bytes += st.Charged
	d := f.cur
	d.own += st.Charged
	d.ownEntries++
	switch st.Kind {
	case core.KindDir:
		f.dirs++
		f.attach(d, p)
	case core.KindLink:
		f.links++
	default:
		f.files++
		d.ownFiles++
		if st.Charged > d.topBytes {
			// The name is copied, not sliced out of the path: a Go substring
			// keeps the whole string it came from alive, and a directory
			// holding the full path of its largest file would keep a hundred
			// bytes where fifteen are wanted — half a million times over.
			d.topBytes, d.topName = st.Charged, strings.Clone(base(p))
		}
	}
	if t := f.top; t != nil {
		t.live += st.Charged
		t.liveEntrie++
	}

	f.tick++
	if f.tick&1023 == 0 && f.want.Load() {
		f.emit()
	}
}

// turn moves the remembered directory to the parent of p.  It runs once per
// directory the walk opens, not once per entry.
func (f *walkFeed) turn(p string) {
	dir := parentOf(p, f.root)
	f.prefix = withSlash(dir)
	// A directory that was not kept (the cap) leaves its children looking for
	// a parent that does not exist.  They are charged to the nearest ancestor
	// that does, so the bytes stay in the tree even where its shape stops.
	for probe := dir; ; {
		if n, ok := f.nodes[probe]; ok {
			f.cur = n
			break
		}
		up := parentOf(probe, f.root)
		if up == probe {
			f.cur = f.tree
			break
		}
		probe = up
	}
	f.depth = f.cur.depth

	f.top = nil
	for c := f.cur; c != nil && c.parent != nil; c = c.parent {
		if c.parent == f.tree {
			f.top = c
			break
		}
	}
}

// attach makes a node for a subdirectory just met, unless the cap is reached.
// The root is not a child of itself: it is accounted like any other entry, and
// without this it would appear inside its own listing at nought bytes.
func (f *walkFeed) attach(parent *wnode, p string) {
	if p == f.root {
		return
	}
	if f.count >= maxNodes {
		f.truncated = true
		return
	}
	n := &wnode{name: strings.Clone(base(p)), parent: parent, depth: parent.depth + 1}
	parent.kids = append(parent.kids, n)
	f.nodes[p] = n
	f.count++
}

// base is the last component of a path.  The result shares memory with the
// path the walk already allocated, so a name costs nothing but its header.
func base(p string) string {
	if i := strings.LastIndexByte(p, '/'); i >= 0 {
		return p[i+1:]
	}
	return p
}

// parentOf is the directory holding p, expressed the way the walk spells its
// paths: filepath.Join is what built them, so "." and "/" are the two roots
// that have to answer for themselves.
func parentOf(p, root string) string {
	if p == root {
		return root
	}
	i := strings.LastIndexByte(p, '/')
	switch {
	case i < 0:
		return "."
	case i == 0:
		return "/"
	}
	return p[:i]
}

// emit hands the screen a snapshot.  It never blocks: a screen that has not
// picked up the last one does not get to hold the walk still.
func (f *walkFeed) emit() {
	f.want.Store(false)
	select {
	case f.out <- f.snapshot():
	default:
	}
}

func (f *walkFeed) snapshot() walkSnap {
	s := walkSnap{
		At: time.Now(), Entries: f.entries, Files: f.files, Dirs: f.dirs,
		Links: f.links, Bytes: f.bytes, Depth: f.depth, Truncated: f.truncated,
	}
	if f.cur != nil {
		s.Cur = f.cur.path()
	}
	s.Tops = make([]topRow, 0, len(f.tree.kids)+1)
	for _, k := range f.tree.kids {
		s.Tops = append(s.Tops, topRow{Name: k.name + "/", Bytes: k.live, Entries: k.liveEntrie})
	}
	// What lies directly in the root is part of what fills it, and it is
	// shown as one row rather than left out of the picture.
	if f.tree.own > 0 {
		s.Tops = append(s.Tops, topRow{Own: true, Bytes: f.tree.own, Entries: f.tree.ownFiles})
	}
	sort.Slice(s.Tops, func(i, j int) bool {
		if s.Tops[i].Bytes != s.Tops[j].Bytes {
			return s.Tops[i].Bytes > s.Tops[j].Bytes
		}
		return s.Tops[i].Name < s.Tops[j].Name
	})
	return s
}

// settle turns what lies directly in each directory into subtree totals, and
// lets the lookup table go.  It runs once, between the walk and the browsing.
func (f *walkFeed) settle() *wnode {
	var sum func(*wnode) (int64, int32)
	sum = func(n *wnode) (int64, int32) {
		b, e := n.own, n.ownEntries
		for _, k := range n.kids {
			kb, ke := sum(k)
			b += kb
			e += ke
		}
		n.bytes, n.entries = b, e
		return b, e
	}
	sum(f.tree)
	f.nodes = nil
	return f.tree
}
