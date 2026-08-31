// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package scan walks a directory tree with lstat only, turns every entry into
// a contract record, and feeds it to the decision layer.
//
// The walk never follows symbolic links (they are recorded as Ссылка and not
// descended into), never leaves the starting filesystem unless asked to, and
// treats an unreadable entry as a fact to report (доступен = ложь) rather than
// a reason to stop.
//
// Size accounting matches du(1) -sb exactly: the apparent st_size of every
// distinct inode, with hard-linked inodes counted once, symlinks counted by
// the length of their target string, and directories charged nothing for
// themselves — GNU du with --apparent-size reports 0 for a directory's own
// size, so counting it here would inflate the total by roughly 4 KiB per
// directory.  A directory's own st_size is still reported separately as
// DirBytes.  That makes TotalBytes directly checkable against an outside tool.
package scan

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"digitdisk/internal/core"
)

// Options configures a walk.
type Options struct {
	Root        string
	CrossDevice bool // descend into mounts of other filesystems
	Decider     core.Decider
	Now         time.Time
	Top         int  // how many entries to keep in each ranking
	MaxDepth    int  // 0 = unlimited
	FollowRoot  bool // resolve the root itself if it is a symlink
}

// Entry is one path kept for a ranking, with the verdict it received.
type Entry struct {
	Path    string       `json:"путь"`
	Size    int64        `json:"размер"`
	AgeDays float64      `json:"возраст_дней"`
	Kind    core.Kind    `json:"вид"`
	Class   core.Class   `json:"разряд"`
	Verdict core.Verdict `json:"приговор"`
	Weight  float64      `json:"вес"`
}

// Bucket is a count/size pair.
type Bucket struct {
	Count int   `json:"count"`
	Bytes int64 `json:"bytes"`
}

// Skips counts what the walk could not look at.
type Skips struct {
	PermissionDenied int      `json:"permission_denied"`
	Vanished         int      `json:"vanished"`
	OtherErrors      int      `json:"other_errors"`
	DeviceBoundaries int      `json:"device_boundaries"`
	DepthLimited     int      `json:"depth_limited"`
	Examples         []string `json:"examples,omitempty"`
}

// Total is every entry the walk refused or failed to descend into.
func (s Skips) Total() int {
	return s.PermissionDenied + s.Vanished + s.OtherErrors + s.DeviceBoundaries + s.DepthLimited
}

// Result is the whole outcome of a walk.
type Result struct {
	Root            string `json:"root"`
	CrossDevice     bool   `json:"cross_device"`
	Decider         string `json:"decider"`
	DeciderReady    bool   `json:"decider_ready"`
	ContractVersion int    `json:"contract_version"`
	Entries         int    `json:"entries"`
	Files           int    `json:"files"`
	Dirs            int    `json:"dirs"`
	Links           int    `json:"links"`
	Others          int    `json:"others"`
	Unreadable      int    `json:"unreadable"`
	// TotalBytes is files + symlinks + other inodes, du -sb compatible.
	TotalBytes int64 `json:"total_bytes"`
	FileBytes  int64 `json:"file_bytes"`
	// DirBytes is the sum of the directories' own st_size.  It is reported
	// for information and is NOT part of TotalBytes, matching du.
	DirBytes        int64                   `json:"dir_bytes"`
	LinkBytes       int64                   `json:"link_bytes"`
	HardlinkDupes   int                     `json:"hardlink_duplicates"`
	HardlinkBytes   int64                   `json:"hardlink_duplicate_bytes"`
	Skipped         Skips                   `json:"skipped"`
	ByClass         map[core.Class]Bucket   `json:"by_class"`
	ByVerdict       map[core.Verdict]Bucket `json:"by_verdict"`
	Removable       []Entry                 `json:"removable_top"`
	Largest         []Entry                 `json:"largest_top"`
	DurationSeconds float64                 `json:"duration_seconds"`
}

type topList struct {
	limit int
	items []Entry
}

func (t *topList) add(e Entry) {
	if t.limit <= 0 {
		return
	}
	t.items = append(t.items, e)
	if len(t.items) > t.limit*4 {
		t.trim()
	}
}

func (t *topList) trim() {
	sort.Slice(t.items, func(i, j int) bool {
		if t.items[i].Size != t.items[j].Size {
			return t.items[i].Size > t.items[j].Size
		}
		return t.items[i].Path < t.items[j].Path
	})
	if len(t.items) > t.limit {
		t.items = t.items[:t.limit]
	}
}

func (t *topList) result() []Entry {
	t.trim()
	if t.items == nil {
		return []Entry{}
	}
	return t.items
}

type inode struct {
	dev uint64
	ino uint64
}

// Walk runs the traversal described in the package comment.
func Walk(opt Options) (Result, error) {
	if opt.Decider == nil {
		opt.Decider = core.Default()
	}
	if opt.Now.IsZero() {
		opt.Now = time.Now()
	}
	if opt.Top <= 0 {
		opt.Top = 15
	}
	started := time.Now()

	res := Result{
		Root:            opt.Root,
		CrossDevice:     opt.CrossDevice,
		Decider:         opt.Decider.Name(),
		DeciderReady:    opt.Decider.Ready(),
		ContractVersion: core.ContractVersion,
		ByClass:         map[core.Class]Bucket{},
		ByVerdict:       map[core.Verdict]Bucket{},
	}

	rootInfo, err := os.Lstat(opt.Root)
	if err != nil {
		return res, err
	}
	if opt.FollowRoot && rootInfo.Mode()&fs.ModeSymlink != 0 {
		if rootInfo, err = os.Stat(opt.Root); err != nil {
			return res, err
		}
	}
	rootDev := devOf(rootInfo)

	seen := make(map[inode]struct{}, 1024)
	removable := &topList{limit: opt.Top}
	largest := &topList{limit: opt.Top}

	// account turns one lstat result into a record, asks the decision layer
	// about it, and folds the answer into the totals.
	account := func(path string, info fs.FileInfo, accessible bool) {
		res.Entries++
		rec := core.Record{Path: path, Accessible: accessible}
		var size int64
		k := KindOther
		if info != nil {
			k = kindOf(info)
			size = info.Size()
			rec.AgeDays = opt.Now.Sub(info.ModTime()).Hours() / 24
			if rec.AgeDays < 0 {
				rec.AgeDays = 0
			}
			// Hard-linked inodes contribute their bytes once, the way
			// du(1) counts them.
			if nlink(info) > 1 {
				id := inode{dev: devOf(info), ino: inoOf(info)}
				if _, dup := seen[id]; dup {
					res.HardlinkDupes++
					res.HardlinkBytes += size
					size = 0
				} else {
					seen[id] = struct{}{}
				}
			}
			rec.Size = info.Size()
		}
		// The contract knows exactly three виды; anything else on disk
		// (socket, fifo, device node, or an entry we could not lstat) is
		// handed over as Файл while our own totals keep it separate.
		rec.Kind = k
		if k == KindOther {
			rec.Kind = core.KindFile
		}

		charged := size
		switch k {
		case core.KindDir:
			res.Dirs++
			res.DirBytes += size
			charged = 0 // du charges a directory nothing for itself
		case core.KindLink:
			res.Links++
			res.LinkBytes += size
		case core.KindFile:
			res.Files++
			res.FileBytes += size
		default:
			res.Others++
		}
		if !accessible {
			res.Unreadable++
		}
		res.TotalBytes += charged

		d := opt.Decider.Decide(rec)
		// Buckets are charged the same bytes as TotalBytes, so they sum
		// back to it; an Entry still shows the path's own st_size.
		bc := res.ByClass[d.Class]
		bc.Count++
		bc.Bytes += charged
		res.ByClass[d.Class] = bc
		bv := res.ByVerdict[d.Verdict]
		bv.Count++
		bv.Bytes += charged
		res.ByVerdict[d.Verdict] = bv

		e := Entry{Path: path, Size: rec.Size, AgeDays: rec.AgeDays, Kind: k,
			Class: d.Class, Verdict: d.Verdict, Weight: d.Weight}
		largest.add(e)
		if d.Verdict == core.VerdictRemovable {
			removable.add(e)
		}
	}

	noteSkip := func(path string, err error) {
		switch {
		case errors.Is(err, fs.ErrPermission):
			res.Skipped.PermissionDenied++
		case errors.Is(err, fs.ErrNotExist):
			res.Skipped.Vanished++
		default:
			res.Skipped.OtherErrors++
		}
		if len(res.Skipped.Examples) < 10 {
			res.Skipped.Examples = append(res.Skipped.Examples, path+": "+err.Error())
		}
	}

	account(opt.Root, rootInfo, true)

	type job struct {
		path  string
		depth int
	}
	var stack []job
	if rootInfo.IsDir() {
		stack = append(stack, job{opt.Root, 0})
	}

	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		entries, err := os.ReadDir(cur.path)
		if err != nil {
			noteSkip(cur.path, err)
			continue
		}
		for _, de := range entries {
			child := filepath.Join(cur.path, de.Name())
			info, err := de.Info()
			if err != nil {
				noteSkip(child, err)
				account(child, nil, false)
				continue
			}
			account(child, info, true)
			if !info.IsDir() {
				continue // symlinks are never followed
			}
			if !opt.CrossDevice && devOf(info) != rootDev {
				res.Skipped.DeviceBoundaries++
				continue
			}
			if opt.MaxDepth > 0 && cur.depth+1 >= opt.MaxDepth {
				res.Skipped.DepthLimited++
				continue
			}
			stack = append(stack, job{child, cur.depth + 1})
		}
	}

	res.Removable = removable.result()
	res.Largest = largest.result()
	res.DurationSeconds = time.Since(started).Seconds()
	return res, nil
}

// KindOther labels entries that are neither file, directory nor symlink
// (sockets, fifos, device nodes) plus entries that could not be lstat'ed.  It
// exists only in this package's own totals: the contract record always carries
// one of the three виды it defines.
const KindOther = core.Kind("Прочее")

func kindOf(info fs.FileInfo) core.Kind {
	switch {
	case info.Mode()&fs.ModeSymlink != 0:
		return core.KindLink
	case info.IsDir():
		return core.KindDir
	case info.Mode().IsRegular():
		return core.KindFile
	default:
		return KindOther
	}
}

func statOf(info fs.FileInfo) *syscall.Stat_t {
	st, _ := info.Sys().(*syscall.Stat_t)
	return st
}

func devOf(info fs.FileInfo) uint64 {
	if st := statOf(info); st != nil {
		return uint64(st.Dev)
	}
	return 0
}

func inoOf(info fs.FileInfo) uint64 {
	if st := statOf(info); st != nil {
		return st.Ino
	}
	return 0
}

func nlink(info fs.FileInfo) uint64 {
	if st := statOf(info); st != nil {
		return uint64(st.Nlink)
	}
	return 1
}
