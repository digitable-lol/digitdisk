// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package clean removes files — and is built so that removing them is hard to
// regret.
//
// # What may be removed
//
// Exactly one thing: a path the decision layer gave the приговор
// «МожноУбрать».  Not a path that looks like one, not a path matching a
// pattern, not a directory whose name is on a list of well-known caches.  The
// host never forms an opinion about what a path is; it asks core.Decider and
// obeys.  That is the whole difference between this and a cleaner that ships a
// list of paths: the rule that decides is proved once, in flang, and printed
// into three languages that agree.
//
// The host does keep a veto.  Every item the layer marked «МожноУбрать» is put
// through the checks in guard() before it is touched, and an item that fails
// one is refused, counted and reported — never removed.  The checks restate,
// in Go, what the flang rules already prove (a каталог is never «МожноУбрать»,
// nor is a ссылка, nor is anything недоступное).  Two layers agreeing is not
// redundancy here: if they ever disagree, the disagreement is a fact about the
// rules that a person must see, and `clean` prints it instead of acting on it.
//
// # Three steps, not one
//
//	digitdisk clean <путь>            план: что, сколько, почему.  Ничего не тронуто.
//	digitdisk clean <путь> --apply    перенос в корзину.  Обратимо.
//	digitdisk restore <корзина>       возврат на прежние места.
//	digitdisk purge <корзина> --confirm N   стирание.  Необратимо.
//
// Only the last step calls os.Remove, and only on files inside a корзина this
// tool itself filled, named in a journal it itself wrote.  Moving to the
// корзина frees no space at all — the bytes are still on the disk under a
// different name — and the plan says so rather than reporting a saving that
// has not happened.
//
// # Why the корзина lives inside the корень
//
// The корзина is `<корень>/.digitdisk-trash/<метка времени>/` and may be moved
// with --trash only to another place inside the same корень.  Two properties
// pay for that restriction:
//
//   - The move is rename(2): atomic, O(1), no copying, no window in which the
//     file exists twice or not at all.  A корзина on another filesystem would
//     turn every "move" into copy-then-delete — the cost of reversibility
//     would become the size of what you are cleaning, and a crash mid-copy
//     would leave a half-file.  Refusing that is cheaper than explaining it.
//   - Every path operation goes through os.Root, opened on the корень.  A
//     *os.Root resolves each component itself and refuses to leave the
//     directory it was opened on, symbolic links included.  "Does not escape
//     the корень" is therefore a property of the syscalls used, not a check
//     that could be forgotten — and it holds even if a directory is swapped
//     for a symlink after the walk has passed it.
//
// # Crash safety
//
// The journal is written in full BEFORE the first file moves, listing every
// intended путь → в_корзине pair, and rewritten with the outcome afterwards.
// A crash in between therefore leaves a корзина that `restore` can still empty
// back: an entry whose file never made it is simply skipped.  The opposite
// order — move first, record after — would produce exactly the state nobody
// can undo, a корзина full of files nobody can put back.
package clean

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"digitdisk/internal/core"
	"digitdisk/internal/lang"
	"digitdisk/internal/places"
	"digitdisk/internal/protect"
	"digitdisk/internal/scan"
)

// TrashName is the directory, always directly under the корень, that holds
// every корзина this tool makes.  It is excluded from the walk by name at any
// depth: without that, a second run would find the first run's корзина, decide
// its contents are still кэш, and move them into a корзина inside the корзина.
const TrashName = ".digitdisk-trash"

// JournalName is the machine-readable record inside one корзина.
const JournalName = "journal.json"

// FilesDir is where a корзина keeps the files themselves, under their path
// relative to the корень.  Keeping the shape means restore needs no cleverness
// to know where a file came from, and a person reading the корзина can see it.
const FilesDir = "files"

// JournalVersion is the shape of journal.json.  It is written into the file so
// a later version can refuse a journal it does not understand rather than
// misread one.
const JournalVersion = 1

// Options configures a plan.
type Options struct {
	Root        string
	Trash       string // корзина root; empty means <корень>/.digitdisk-trash
	CrossDevice bool
	MaxDepth    int
	Decider     core.Decider
	Now         time.Time
	Version     string // digitdisk version, recorded in the journal

	// Protect is the operator standing veto.  It is asked only about items
	// the decision layer already marked «МожноУбрать», and it can only
	// subtract from the plan.  See internal/protect for why it lives here
	// and not in the rules.
	Protect *protect.List

	// Places names the справочник the decision layer was given, so the plan
	// can say which place claimed an item.  It decides nothing: the разряд
	// on every item came from the layer.
	Places *places.Directory
}

// Identity is what the host remembers about a file between the moment it was
// looked at and the moment it is touched.  It is compared field for field
// before a file is moved or erased: a difference means the file is not the one
// that was judged, and an item that is not the one that was judged is refused.
//
// What it catches: replacement (dev/ino), a write (size, mtime), a chmod
// (mode), a change of type.  What it does not: a hard link added or removed
// elsewhere, which changes neither the bytes nor the name being moved, and a
// write that restores the previous size and mtime to the nanosecond.
type Identity struct {
	Dev           uint64 `json:"устройство"`
	Ino           uint64 `json:"узел"`
	Size          int64  `json:"размер"`
	MtimeUnixNano int64  `json:"изменён_наносекунд"`
	Mode          uint32 `json:"режим"`
	Nlink         uint64 `json:"ссылок"`
}

// Same reports whether two identities describe the same file in the same
// state.  Nlink is deliberately not compared: see the type comment.
func (i Identity) Same(o Identity) bool {
	return i.Dev == o.Dev && i.Ino == o.Ino && i.Size == o.Size &&
		i.MtimeUnixNano == o.MtimeUnixNano && i.Mode == o.Mode
}

// Differs names the first field that changed, for a refusal a person can act
// on.  "изменился" without saying what changed is not a report.
func (i Identity) Differs(o Identity) lang.Phrase {
	switch {
	case i.Dev != o.Dev || i.Ino != o.Ino:
		return lang.Say("это уже другой файл (был узел %d:%d, стал %d:%d)", i.Dev, i.Ino, o.Dev, o.Ino)
	case i.Size != o.Size:
		return lang.Say("размер изменился (был %d Б, стал %d Б)", i.Size, o.Size)
	case i.MtimeUnixNano != o.MtimeUnixNano:
		return lang.Say("в файл писали после обхода (время изменения другое)")
	case i.Mode != o.Mode:
		return lang.Say("права изменились (были %v, стали %v)", fs.FileMode(i.Mode), fs.FileMode(o.Mode))
	}
	return lang.Phrase{}
}

func identityOf(info fs.FileInfo) Identity {
	id := Identity{
		Size:          info.Size(),
		MtimeUnixNano: info.ModTime().UnixNano(),
		Mode:          uint32(info.Mode()),
		Nlink:         1,
	}
	if st, ok := info.Sys().(*syscall.Stat_t); ok && st != nil {
		id.Dev = uint64(st.Dev)
		id.Ino = st.Ino
		id.Nlink = uint64(st.Nlink)
	}
	return id
}

// Item is one file the plan proposes to move, with the layer's verdict and the
// numbers the verdict was made from.
type Item struct {
	Path          string       `json:"путь"`
	Rel           string       `json:"путь_от_корня"`
	Size          int64        `json:"размер"`
	AgeDays       float64      `json:"возраст_дней"`
	Class         core.Class   `json:"разряд"`
	Verdict       core.Verdict `json:"приговор"`
	Weight        float64      `json:"вес"`
	ThresholdDays float64      `json:"порог_дней,omitempty"`
	HasThreshold  bool         `json:"порог_известен"`

	// Place names the row of the справочник that claims this path, when one
	// does.  It is written for the report only: the разряд above came from
	// the decision layer, and this says which line of the file the layer was
	// given that could account for it.
	Place      string   `json:"место,omitempty"`
	Hardlinked bool     `json:"жёсткая_ссылка"`
	Before     Identity `json:"отпечаток"`

	// Filled in by Apply.
	TrashRel string    `json:"в_корзине,omitempty"`
	MovedAt  string    `json:"перенесено,omitempty"`
	After    *Identity `json:"отпечаток_после,omitempty"`

	// Filled in by Restore and Purge.
	RestoredAt string `json:"возвращено,omitempty"`
	PurgedAt   string `json:"стёрто,omitempty"`

	// Failed is why this item was NOT acted on, when the file turned out to
	// have changed between the walk and the moment of touching it.  An empty
	// MovedAt with a filled Failed is the whole record of a refusal.
	Failed lang.Phrase `json:"не_сделано,omitzero"`
}

// Why states, in the layer's own terms, why this item is on the list.  Every
// number in it is one the decision layer used: the разряд it assigned, the
// порог it applies to that разряд, and the возраст the host measured and
// handed it.  The host adds no rule of its own here — it has none.
func (i Item) Why(l lang.Lang) string {
	if i.HasThreshold {
		return l.F("разряд %s: возраст %s ≥ порога %s", l.Word(string(i.Class)), l.Days(i.AgeDays), l.Days(i.ThresholdDays))
	}
	return l.F("разряд %s, приговор %s", l.Word(string(i.Class)), l.Word(string(i.Verdict)))
}

// Where names the известное место this path lies in, or an empty string when
// the разряд came from the general приметы instead.  Two different sentences,
// and a person deciding whether to trust a line wants to know which one it is.
func (i Item) Where() string { return i.Place }

// Refusal is a path the decision layer marked «МожноУбрать» that the host will
// not touch, and the check that stopped it.  A refusal is never a silent skip:
// it means the two layers disagree, and that is a fact about the rules.
type Refusal struct {
	Path    string       `json:"путь"`
	Class   core.Class   `json:"разряд"`
	Verdict core.Verdict `json:"приговор_ядра"`
	Reason  lang.Phrase  `json:"отказ"`
}

// Protected is a path the decision layer marked «МожноУбрать» and the operator
// told the tool to leave alone.  It is kept apart from Refusal on purpose: a
// refusal means the two layers disagree and somebody should look at the rules,
// while this means the rules worked and the answer was overruled by the person
// whose disk it is.  Mixing them would make the first invisible among the
// second on any machine with a защитный список.
type Protected struct {
	Path  string       `json:"путь"`
	Class core.Class   `json:"разряд"`
	Size  int64        `json:"размер"`
	Rule  protect.Rule `json:"правило"`
}

// Plan is the whole answer to "what would `clean` remove here".
type Plan struct {
	Root            string      `json:"корень"`
	Trash           string      `json:"корзина"`
	Decider         string      `json:"решающий_слой"`
	DeciderReady    bool        `json:"решающий_слой_готов"`
	ContractVersion int         `json:"версия_договора"`
	Walk            scan.Result `json:"обход"`
	Items           []Item      `json:"к_переносу"`
	Refused         []Refusal   `json:"отказано"`
	Protected       []Protected `json:"защищено"`
	ProtectOrigins  []string    `json:"защитный_список,omitempty"`
	ProtectedBytes  int64       `json:"защищено_байт"`
	PlacesOrigin    string      `json:"справочник,omitempty"`
	PlacesCount     int         `json:"мест_в_справочнике"`
	ByClass         []ClassSum  `json:"по_разрядам"`
	Bytes           int64       `json:"байт"`
	FreeableBytes   int64       `json:"освободится_стиранием"`
	HardlinkItems   int         `json:"из_них_жёстких_ссылок"`

	// PrunedTrash counts корзины the walk refused to descend into, not the
	// files inside them: a pruned directory is never opened, so nobody here
	// knows how much was in it.
	PrunedTrash int `json:"своих_корзин_пропущено"`
}

// ClassSum is the plan broken down by разряд.  It is computed over EVERY item,
// never over the shortened list a report prints: a summary that changed with
// --top would be a summary of the screen, not of the disk.
type ClassSum struct {
	Class core.Class `json:"разряд"`
	Count int        `json:"файлов"`
	Bytes int64      `json:"байт"`
}

// Make walks the tree and returns what `clean --apply` would move.  It opens
// no file for writing and creates no корзина: a plan is a question, and asking
// it must be free of consequences or nobody will ask it first.
func Make(opt Options) (Plan, error) {
	if opt.Decider == nil {
		opt.Decider = core.Default()
	}
	if opt.Now.IsZero() {
		opt.Now = time.Now()
	}

	rootAbs, err := filepath.Abs(opt.Root)
	if err != nil {
		return Plan{}, err
	}
	rootAbs = filepath.Clean(rootAbs)

	trashAbs, err := trashRoot(rootAbs, opt.Trash)
	if err != nil {
		return Plan{}, err
	}

	p := Plan{
		Root:            rootAbs,
		Trash:           trashAbs,
		Decider:         opt.Decider.Name(),
		DeciderReady:    opt.Decider.Ready(),
		ContractVersion: core.ContractVersion,
		Items:           []Item{},
		Refused:         []Refusal{},
		Protected:       []Protected{},
	}
	if opt.Places != nil {
		p.PlacesOrigin = opt.Places.Origin
		p.PlacesCount = len(opt.Places.Applicable())
	}
	if opt.Protect != nil {
		p.ProtectOrigins = opt.Protect.Origins
	}

	thresholds, _ := opt.Decider.(core.Thresholder)

	res, err := scan.Walk(scan.Options{
		Root:        rootAbs,
		CrossDevice: opt.CrossDevice,
		MaxDepth:    opt.MaxDepth,
		Top:         1, // rankings are a report; the work list is built below
		Decider:     opt.Decider,
		Now:         opt.Now,
		Prune: func(path string) bool {
			if filepath.Base(path) == TrashName {
				p.PrunedTrash++
				return true
			}
			return false
		},
		Observe: func(e scan.Entry, info fs.FileInfo) {
			if e.Verdict != core.VerdictRemovable {
				return
			}
			// The operator's veto comes before the host's own checks:
			// "I said do not touch this" is an answer to the whole
			// question, and running the other checks on a path nobody
			// may touch would only produce noise about it.
			if rule, ok := opt.Protect.Covers(e.Path, e.Class); ok {
				p.Protected = append(p.Protected, Protected{
					Path: e.Path, Class: e.Class, Size: e.Size, Rule: rule,
				})
				return
			}
			rel, reason := guard(rootAbs, e, info)
			if !reason.Empty() {
				p.Refused = append(p.Refused, Refusal{
					Path: e.Path, Class: e.Class, Verdict: e.Verdict, Reason: reason,
				})
				return
			}
			id := identityOf(info)
			it := Item{
				Path: e.Path, Rel: rel, Size: e.Size, AgeDays: e.AgeDays,
				Class: e.Class, Verdict: e.Verdict, Weight: e.Weight,
				Hardlinked: id.Nlink > 1, Before: id,
			}
			if thresholds != nil {
				it.ThresholdDays, it.HasThreshold = thresholds.Threshold(e.Class)
			}
			// Only for the items that made the plan — a few hundred paths,
			// not the millions the walk visited.
			if opt.Places != nil {
				if place, ok := opt.Places.Match(e.Path); ok {
					it.Place = place.Name
				}
			}
			p.Items = append(p.Items, it)
		},
	})
	if err != nil {
		return p, err
	}
	p.Walk = res

	sort.Slice(p.Items, func(i, j int) bool {
		if p.Items[i].Size != p.Items[j].Size {
			return p.Items[i].Size > p.Items[j].Size
		}
		return p.Items[i].Path < p.Items[j].Path
	})
	byClass := map[core.Class]ClassSum{}
	for _, it := range p.Items {
		p.Bytes += it.Size
		sum := byClass[it.Class]
		sum.Class, sum.Count, sum.Bytes = it.Class, sum.Count+1, sum.Bytes+it.Size
		byClass[it.Class] = sum
		if it.Hardlinked {
			p.HardlinkItems++
			continue // another name keeps the bytes: erasing this one frees nothing
		}
		p.FreeableBytes += it.Size
	}
	// Report order, not map order: the summary is read by a person, and a
	// summary whose lines move between runs is read twice every time.
	p.ByClass = []ClassSum{}
	for _, c := range core.Classes {
		if sum, ok := byClass[c]; ok {
			p.ByClass = append(p.ByClass, sum)
		}
	}
	sort.Slice(p.Refused, func(i, j int) bool { return p.Refused[i].Path < p.Refused[j].Path })
	sort.Slice(p.Protected, func(i, j int) bool {
		if p.Protected[i].Size != p.Protected[j].Size {
			return p.Protected[i].Size > p.Protected[j].Size
		}
		return p.Protected[i].Path < p.Protected[j].Path
	})
	for _, pr := range p.Protected {
		p.ProtectedBytes += pr.Size
	}
	return p, nil
}

// guard is the host's own veto over the decision layer.  It returns the path
// relative to the корень and an empty reason when the item may be moved, or a
// reason why it may not.
//
// Every check here duplicates something the flang rules already prove.  That
// is the point: this is the assertion that the proof and the running program
// are about the same file.
func guard(rootAbs string, e scan.Entry, info fs.FileInfo) (rel string, reason lang.Phrase) {
	if e.Verdict != core.VerdictRemovable {
		return "", lang.Say("приговор %s — убирается только «%s»", e.Verdict, core.VerdictRemovable)
	}
	if info == nil {
		return "", lang.Say("файл не удалось прочитать при обходе — недоступное не трогаем")
	}
	switch e.Kind {
	case core.KindDir:
		return "", lang.Say("это каталог: снос каталога рекурсивен, и ядро его «МожноУбрать» не выдаёт (правило П3)")
	case core.KindLink:
		return "", lang.Say("это символическая ссылка: ядро её не убирает (правило П2)")
	case core.KindFile:
	default:
		return "", lang.Say("это %s, а не обычный файл", e.Kind)
	}
	if !info.Mode().IsRegular() {
		return "", lang.Say("не обычный файл (%v)", info.Mode())
	}
	rel, err := filepath.Rel(rootAbs, e.Path)
	if err != nil {
		return "", lang.Say("путь не считается от корня: %v", err)
	}
	if rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", lang.Say("путь вне указанного корня")
	}
	if filepath.IsAbs(rel) {
		return "", lang.Say("путь вне указанного корня")
	}
	for _, part := range strings.Split(rel, string(filepath.Separator)) {
		if part == TrashName {
			return "", lang.Say("путь внутри корзины digitdisk — своё же убирает `purge`, не `clean`")
		}
	}
	return filepath.ToSlash(rel), lang.Phrase{}
}

// trashRoot resolves where корзины go and refuses a place outside the корень.
func trashRoot(rootAbs, given string) (string, error) {
	if given == "" {
		return filepath.Join(rootAbs, TrashName), nil
	}
	abs, err := filepath.Abs(given)
	if err != nil {
		return "", err
	}
	abs = filepath.Clean(abs)
	rel, err := filepath.Rel(rootAbs, abs)
	if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", lang.Errorf(`корзина %s лежит вне корня %s.
Корзина обязана быть внутри корня: тогда перенос — это rename(2), то есть
мгновенно и без копирования, и все обращения идут через os.Root, который
из корня не выпускает даже по символической ссылке. Корзина на другой
файловой системе превратила бы перенос в копирование: цена обратимости
стала бы равна объёму уборки, а обрыв на середине оставил бы полфайла`,
			abs, rootAbs)
	}
	return abs, nil
}

// phraseOf keeps a refusal this package wrote as the wording it was written
// as, so that «не сделано» reads in the language of whoever is looking, and
// leaves a message from the system alone — nobody here wrote it and nobody
// here can translate it honestly.
//
// The Russian rendering is the same either way: an Error renders its Phrase,
// and the Phrase is what the Error was built from.  The журнал therefore keeps
// the byte it always kept.
func phraseOf(err error) lang.Phrase {
	var ours *lang.Error
	if errors.As(err, &ours) {
		return ours.P
	}
	return lang.FromError(err)
}

// stamp names one корзина.  Colons are legal in a POSIX filename and illegal
// on the FAT volumes people keep on USB sticks; dashes cost nothing.
func stamp(t time.Time) string {
	return t.UTC().Format("2006-01-02T15-04-05Z")
}

// relUnder returns target's path relative to rootAbs in os.Root form (slashes,
// no leading separator), refusing anything that is not below rootAbs.
func relUnder(rootAbs, target string) (string, error) {
	rel, err := filepath.Rel(rootAbs, target)
	if err != nil {
		return "", err
	}
	if rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return "", lang.Errorf("%s лежит вне корня %s", target, rootAbs)
	}
	return filepath.ToSlash(rel), nil
}

// trashRelFor returns the path of one корзина entry relative to the корень.
func trashRelFor(rootAbs, boxAbs, itemRel string) (string, error) {
	rel, err := relUnder(rootAbs, boxAbs)
	if err != nil {
		return "", err
	}
	return path.Join(rel, FilesDir, itemRel), nil
}

// openRoot opens the корень for traversal-safe access.  Every later operation
// is relative to it and cannot leave it.
func openRoot(rootAbs string) (*os.Root, error) {
	r, err := os.OpenRoot(rootAbs)
	if err != nil {
		return nil, lang.Errorf("корень %s не открывается: %s", rootAbs, err)
	}
	return r, nil
}
