**English** · [Русский](README.ru.md)

# digitdisk — where the disk went, how the machine feels, and what may go

digitdisk prints two readings of a machine — **where the disk space went**,
directories by size and the largest files, and **how the machine is feeling
right now**, the processor and each of its cores, memory, disk, network and
video cards — and it can act on the first of them: `clean` removes files, in
three steps, none of which is a surprise.

`status` and `analyze` read and write nothing. `clean` shows a plan and needs
`--apply` to move anything; what it moves goes to a корзина inside the tree you
named and comes back with `restore`; erasing is a separate command with a
separate confirmation. **What may be removed is not a list of well-known paths
and not a pattern**: it is exactly what the decision layer in `core/` gives the
verdict «МожноУбрать», which is proved to be nothing outside Кэш, Журнал and
Сборка.

The concrete places — npm's cache, Go's build cache, Xcode's derived data —
digitdisk knows from a **справочник**, a separate data file edited without a
rebuild. The справочник names a разряд; the verdict is still the core's, and
thresholds, directories, symlinks and content-addressed stores are judged
exactly as before. To say "leave this alone" there is the
[защитный список](#the-protection-list); what past cleanups did is
`digitdisk history`.

## What it is made of

```
core/       the readings as a flang specification, plus the Go printed from it into core/out-go
host/       the Go host: system calls, the command line, the output
host/internal/lang/                the dictionary: what a person reads, in both languages
host/internal/places/places.conf   the справочник of known places: data, not code
packaging/  the Homebrew formula
scripts/    the release build
tools/      the licence gate
docs/       notes that are not this page
```

The split is the whole design. Everything that can be decided without touching
the operating system is decided in `core/`: a specification in flang, checked by
flang's own runs and **printed** into Go. Everything that must touch the
operating system — walking directories, `statfs`, reading counters — is
hand-written Go in `host/`. The core never opens a file; the host never decides
what a number means.

`core/out-go` is printed, not written. Hand edits there are lost at the next
print — see [`AGENTS.md`](AGENTS.md).

## Install

**Released binaries: Linux (x86-64, arm64) and macOS (Apple Silicon, Intel).**
All four are produced by one cross-compilation on Linux, with CGO off and a
repeatable fingerprint — see [macOS](#macos) below.

### Homebrew

```bash
brew install digitable-lol/tap/digitdisk
digitdisk --version
```

The formula installs the released binary. It does not compile anything on your
machine and does not need a Go toolchain. Its source is
[`packaging/homebrew/digitdisk.rb`](packaging/homebrew/digitdisk.rb); the copy
Homebrew reads lives in
[digitable-lol/homebrew-tap](https://github.com/digitable-lol/homebrew-tap).

It also installs both manual pages, and into the two places `man` looks: the
English one as `share/man/man1/digitdisk.1`, the Russian one as
`share/man/ru/man1/digitdisk.1`. So `man digitdisk` answers in English and
`LANG=ru_RU.UTF-8 man digitdisk` answers in Russian, with nothing to configure
— `man` picks the translation by locale on its own.

### A released binary

```bash
V=0.1.0; A=amd64          # or A=arm64
base=https://github.com/digitable-lol/digitdisk/releases/download/v$V
curl -fsSLO $base/digitdisk-$V-linux-$A.tar.gz
curl -fsSLO $base/SHA256SUMS
sha256sum --check --ignore-missing SHA256SUMS
tar -xzf digitdisk-$V-linux-$A.tar.gz
sudo install -m 0755 digitdisk-$V-linux-$A/digitdisk /usr/local/bin/digitdisk
```

Check the sums before unpacking, not after. Every release also carries the
formula, so the two ways install the same bytes.

### From source

A Go toolchain is enough from a clean checkout, because the printed Go is
committed. The build runs inside `host/`: that is where the module lives, and
there is no module at the root.

```bash
cd host && go build -tags flangcore -o ../digitdisk .
```

`-tags flangcore` is what puts the flang core inside the binary. Without it the
host builds against a placeholder that counts but decides nothing —
`digitdisk --version` names which one is inside, so a build is never in doubt.

There are two tags, and the second is not the default yet. `flangui` also puts
the **screen layout in flang** inside the binary — the
[flang-tui](https://github.com/digitable-lol/flang-tui) library, wired in as a
submodule under [`ui-flang/`](ui-flang/README.md):

```bash
cd host && go build -tags flangcore,flangui -o ../digitdisk .
```

Both builds print the same bytes — that is a run, not a promise
(`tools/sverka-ui.sh`). But the flang layout recomputes its postconditions on
every return and costs **13–21 ms per keystroke** against 0.13–0.19 ms for the
hand-written Go. A timer redraw hides behind the ≈1.5 s of collection; a
keystroke is waited for — so the default stays hand-written Go for now. The
numbers, and what would change the default: [`ui-flang/README.md`](ui-flang/README.md).

<a id="macos"></a>

### macOS: the same readings, a different source

The host builds and runs on macOS, arm64 and x86-64. It takes its facts from
`sysctl(3)`, `getfsstat(2)`, the routing socket and the documented functions of
`libSystem` instead of `/proc` and `/sys`; the flang core is untouched by any of
that, because the core has never known what a system call is.

```bash
cd host && GOOS=darwin go build -tags flangcore -o ../digitdisk .
```

**Why there is no cgo.** The obvious way to call a C function from Go is cgo,
and it would have ended the release: the four binaries are cross-compiled on one
Linux machine and checked byte for byte against a second build of themselves,
and cgo ends both properties at once. So the calls are made the way the Go
standard library itself makes them on macOS — the symbol is recorded as a
dynamic import, a two-instruction assembly stub jumps to it, and the call goes
out through `syscall.syscall6`. The Go linker writes the import into the Mach-O
file and the system loader binds it to `libSystem` at start-up, exactly as it
binds the imports the runtime already needs. No Mac is needed to build this; one
is needed only to check it.

**How the layouts are proved.** The decoders are written from Apple's headers,
not from anybody else's source, and no number is printed until its provenance
has been confirmed on the machine itself:

| What is read | What proves it was read correctly |
|---|---|
| process record (`kinfo_proc`) | our own pid, parent and user turn up where we expect them |
| process memory and threads (`proc_taskinfo`) | the kernel says how many bytes it wrote; our own process holds pages and has at least one thread |
| command line (`KERN_PROCARGS2`) | our own arguments match `os.Args` word for word, which the runtime got by another road |
| memory breakdown (`vm_statistics64`) | no page count exceeds the machine's pages; the read-ahead pages do not outnumber the free ones the kernel folds them into; and the disjoint buckets sum to `hw.memsize` within a third of a percent |
| CPU busy share | it is a ratio of two differences, so it depends on no tick rate at all |
| per-processor shares (`processor_cpu_load_info`) | the kernel says how many processors it wrote about, and that is the number `hw.logicalcpu` gives; their mean comes out as the machine-wide share, which is the sum of exactly those counters |
| interface counters (`if_data64`) | the MTU matches what the standard library reports |

If a check does not agree, the field stays empty rather than being printed on a
guess. On top of that, every push runs those same self-checks on live GitHub
macOS runners, on Apple Silicon and on Intel:
`.github/workflows/check.yml` does not only build — it takes a snapshot and
looks for numbers in it.

What macOS measures, and the call each number comes from:

| Reading | Source |
|---|---|
| host, kernel, release, model | `sysctl` `kern.hostname`, `kern.osrelease`, `kern.version`, `kern.osproductversion`, `kern.osversion`, `hw.machine`, `hw.model` |
| uptime | `sysctl kern.boottime` (`struct timeval`) |
| load average, cores | `sysctl vm.loadavg` (`struct loadavg`), `hw.logicalcpu` |
| CPU busy share | `host_statistics(HOST_CPU_LOAD_INFO)` |
| the share of each core | `host_processor_info(PROCESSOR_CPU_LOAD_INFO)` |
| memory total, page size, swap | `sysctl hw.memsize`, `hw.pagesize`, `vm.swapusage` (`struct xsw_usage`) |
| memory free, cache, available, used, wired, compressed | `host_statistics64(HOST_VM_INFO64)` (`struct vm_statistics64`) |
| processes: the list, and how many | `sysctl kern.proc.all` (`struct kinfo_proc`) |
| per-process memory, threads, threads on a processor, CPU time | `proc_pidinfo(PROC_PIDTASKINFO)` (`struct proc_taskinfo`) |
| per-process command lines | `sysctl {CTL_KERN, KERN_PROCARGS2, pid}` |
| disks | `getfsstat(2)` (`struct statfs`) |
| interfaces, addresses | `net.Interfaces` |
| interface counters | `sysctl NET_RT_IFLIST2` (`struct if_data64`) |

"Used" and "available" are the two sides of one statement: memory that is
neither free nor file cache is in use. It is the same reading `free(1)` gives on
Linux, and the report prints the arithmetic next to the number.

**What is still missing on a Mac, and why.** Two kinds, and they are not the
same kind:

- **Closed by permission, not by the language.** The memory, threads and command
  line of a process belonging to *another user* are refused to anybody but the
  administrator: the kernel checks the owner. Running under `sudo` fills those
  rows in; nothing else will.
- **Not published by the system.** What a Mac knows about its video cards lives
  in the IORegistry, and the documented way in is IOKit: Core Foundation
  objects rather than numbers.  We do not read those without cgo and will not
  guess, so the ВИДЕОКАРТЫ section is empty on a Mac.
- **Not published by the system.** Die temperature comes from the SMC through
  IOKit, and Apple documents no interface to it — what circulates is a
  reverse-engineered structure. A number read that way would be a guess wearing
  a unit, so there is none.

The report names them in one line and leaves it at that. The reasons live behind
`digitdisk status --why` and in `--json`.

Two more macOS facts worth knowing before reading a report: a walk of `/`
stops at `/System/Volumes/Data` unless `--cross-device` is given, because the
system and the data volume are two filesystems; and a directory the privacy
machinery refuses is counted as "нет доступа", the same as an unreadable
directory on Linux.

Reprinting the core additionally needs
[flang](https://github.com/digitable-lol/flang), and only when the
specification changed:

```bash
make -C core            # check, print into core/out-go and core/out-c, compare
```

Release archives are built by
[`scripts/build-release.sh`](scripts/build-release.sh), which builds every
target twice and refuses to package if the two builds differ: the same commit
and the same Go toolchain give the same archive, byte for byte.

## Run

```bash
./digitdisk                  # no subcommand: the same as `status`, with its defaults
./digitdisk analyze <path>   # where the space went: directories by size, the largest files
./digitdisk status           # how the machine feels: CPU, memory, disk, network
./digitdisk places           # what the tool knows about concrete caches, and what of it is here
./digitdisk history <path>   # what past cleanups under this root did
./digitdisk --version        # version, build hash, toolchain, decision layer
./digitdisk --help           # subcommands and flags, one line each
./digitdisk status --lang ru # this run in Russian; every subcommand takes --lang ru|en
man digitdisk                # the reference: subcommands, flags, files, examples, exit codes
LANG=ru_RU.UTF-8 man digitdisk   # the same page in Russian
```

All four readings take `--json`. None of them writes anything. A word that is
not a subcommand is refused with code 2, never guessed at; a flag in place of a
subcommand belongs to `status`, so `digitdisk --json` and `digitdisk status
--json` are one command.

### The language of the output

digitdisk writes in Russian and in English, and everything a person reads is in
both: the sections of the report, the labels, the units, the разряды and
приговоры on the screen, `--help`, the refusals, `--why`, the list of commands
on the live screen, and both manual pages.

Who chooses, in the order they are asked:

| | |
|---|---|
| `--lang ru\|en` | this run; every subcommand takes it |
| `DIGITDISK_LANG=ru\|en` | this session |
| `~/.digitable/digitdisk/settings.conf` | what was chosen before |
| the question | asked once, on a first run at a terminal |
| `LC_ALL`, `LC_MESSAGES`, `LANG` | the machine's locale, in that order |
| nothing said | English |

**The default is English, and the reason is POSIX rather than taste.** An unset
locale, `C` and `POSIX` all name the portable locale, whose messages are
English by definition. A machine that has said nothing about its language has
not said "Russian" — it has said "the portable one" — and answering it in
Russian would be a guess about the reader. Somebody who wants Russian either
has a `ru` locale, or is asked once and says so.

**Where there is no terminal, nothing is asked and nothing is written.** A
pipe, a file, a script, a CI job, `--json`: no question, no settings file
brought into being, and the language comes from the locale. Both ends of the
conversation have to be a terminal for the question to happen at all — stdin
and stderr — because a question written to a terminal whose answer would come
from a pipe hangs forever, and a tool that hangs in somebody's build is worse
than a tool in the wrong language.

**Writing in a home directory is an action, and it is announced.** digitdisk
stores two things there: the language — and only after a person answered the
question with their own hands — and a mark that the move of the settings has
already been mentioned, so that it is not mentioned at every run.  Either way
it says what it wrote and where, in one line on stderr — «язык сохранён:
~/.digitable/digitdisk/settings.conf».  If the directory cannot be written — a
read-only mount, a directory owned by somebody else — the run goes on in the
language that was chosen and says plainly that it was not saved: refusing to
look at a disk because a preference could not be stored would be answering a
small problem with a big one.  `digitdisk --version` names the language of the
run and which of the six lines above decided it.

Numbers and dates are written the way each language writes them: «12,3 ГиБ»
against `12.3 GiB`, a non-breaking space against a comma between the thousands,
`02.09.2026` against `2026-09-02`, `Б`/`КиБ`/`МиБ` against `B`/`KiB`/`MiB`,
`дн` against `d`. That is not decoration. «12,3» read as English is twelve and
three, and a report whose numbers change meaning with the reader is worse than
a report in the wrong language, because the wrong language is obvious and a
wrong number is not.

### `--json` is not translated, byte for byte

The keys and the machine values are the same in either language: a script that
parses `digitdisk clean --json` must not care what language the person who ran
it reads. Russian words do travel in that JSON as VALUES, and they stay exactly
where they are — 33 fields carry one. `grep -rn 'json:"' host/internal
host/*.go` lists every field there is; these are the two kinds among them.

**20 of them are identifiers of the договор**, and are not text at all: разряд
(`Кэш`, `Журнал`, `Сборка`, `Загрузка`, `Крупное`, `Неизвестное`), приговор
(`МожноУбрать`, `Спросить`, `НеТрогать`), вид (`Файл`, `Каталог`, `Ссылка`),
якорь (`ОтКорня`, `ГдеУгодно`, and the anchors a справочник row is written
with), the `система` column of the справочник, and the kind of a protection
rule (`путь`, `разряд`) — counting the places where they are the keys of
`by_class` and `by_verdict` rather than a value. They may not change and need
not: they are the names the layer in `core/` proves things about.

**13 of them carry human text**: the refusals `отказ` and `не_сделано`, the
notes `замечание` and `беда`, the map `missing` — whose KEYS are Russian too,
being the names of the readings — `uptime_human` («5д 03:14»), the name of the
decision layer (`решающий_слой` in three records, `decider` in a fourth), and
where the справочник came from (`справочник`, `откуда`). Those are records of
what happened, written once and read back later by `restore`, `purge` and
`history`; rewriting a журнал to suit whoever opens it next would make it a
worse record. `lang.Phrase` is what holds both properties at once — the Russian
wording into the file, the reader's language onto the screen — so translating a
refusal moves no byte of the JSON.

For that second group there is a way forward that breaks nothing, and it is
proposed here rather than done: a machine code beside the Russian value —
`отказ_код` next to `отказ` — with the old field kept for good. A reader that
has always matched on the Russian sentence goes on working, a new one matches
the code, and nothing has to be guessed about which. None of it is written yet.

**The names inside the flang core are not translated and will not be.**
`МожноУбрать` and `Кэш` are identifiers of the layer in `core/`, proved there
and named there. What is translated is the WORD THE SCREEN SHOWS for them, and
that happens in the host, in `host/internal/lang`: the value that arrived is
never touched, so the identifier goes on travelling in the JSON unchanged. That
is where the border runs — `core/` does not know that a language exists, and
not a letter of it moves when the output changes language.

### Cleaning, in three steps

```bash
./digitdisk clean <path>                    # the plan: what, how much, why. Nothing is touched.
./digitdisk clean <path> --apply            # move into <path>/.digitdisk-trash/<stamp>/
./digitdisk restore <trash>                 # put it all back
./digitdisk purge <trash> --confirm N       # erase. This one cannot be undone.
```

The default is the harmless one: `clean` without `--apply` opens no file for
writing and does not even create the корзина, so finding out what it would do
never means having it done.

`--apply` is a `rename(2)` into a корзина inside the same tree, which is why it
is instant and reversible — and why it **frees no space at all**: the bytes are
still there under another name. Only `purge` frees space, it needs `--confirm N`
with N the exact number of files in the корзина, and the failure message does
not tell you N — you get it by running `purge` with no flag and reading the
plan. A confirmation you can satisfy without looking confirms nothing.

Every корзина carries a `journal.json`: what was moved, from where, how many
bytes, when, and where it went. It is written *before* the first file moves, so
a crash in the middle still leaves something `restore` can empty back, and it
survives `purge` as the record of what is gone.

A file that changed between the walk and the move is not moved. digitdisk
remembers each file's dev/ino, size, mtime and mode, checks them again before
touching it, and refuses by name — "размер изменился (был 25 Б, стал 30 Б)" —
rather than removing something it no longer recognises.

**The plan is meant to be read.** Every list stops at `--top` (15 by default,
the same as `analyze`; `--top 0` prints all of it) and ends with a line saying
how many files and bytes were left out. The counts do not move: the total, the
bytes and the breakdown by разряд are computed over the whole plan and are
independent of `--top` — a summary that shrank with the screen would be a
summary of the screen. `--json` is never cut: `--json` is how scripts call this
tool, and a shortened work list would make `clean --json | jq` quietly wrong.

<a id="the-directory-of-known-places"></a>

### The справочник of known places

The core's rules know what a cache IS in general: a path component called
`.cache`, `Caches`, `cache`. That is enough to recognise a cache and not enough
to recognise **npm's**, whose store is `~/.npm/_cacache` — no component with the
word `cache` anywhere in it. The missing knowledge is a LIST OF PLACES, and a
list is data.

```bash
./digitdisk places                      # the whole справочник and what of it is here
./digitdisk places --json               # the same for a machine
./digitdisk clean <path> --places FILE  # your own справочник instead of the built-in one
./digitdisk clean <path> --no-places    # judge by приметы alone, as before 0.4.0
```

It lives in
[`host/internal/places/places.conf`](host/internal/places/places.conf), travels
inside the binary as the default, and is replaced whole — by `--places` or by
`~/.digitable/digitdisk/places.conf`. A row looks like this:

```
разряд | якорь | система | путь | переменная | имя | источник | имя_en
кэш|дом|все|.npm//_cacache|npm_config_cache|npm: кэш загрузок|https://docs.npmjs.com/...|npm: download cache
```

The SOURCE — the seventh field — is mandatory, and not for decoration: **every
row comes from the tool's own documentation**, and a place is listed only when
that documentation calls it a cache, a log, or derived data the tool will
rebuild by itself.  The eighth field is the English name and is optional: a row
without it is read exactly as before and shows its Russian name in either
language.  All 102 rows of the built-in directory have one — `cd host && go
test ./internal/places/` says how many.

The double slash splits the path into a base and a tail: the base is what the
environment variable relocates (`npm_config_cache`, `GOCACHE`, `CARGO_HOME`,
`GRADLE_USER_HOME` and others), so a developer who moved a cache gets their real
place instead of one that is no longer there.

**How it reaches the verdict.** A row becomes a "цепь" — the place's path with a
slash at both ends, `/home/u/.npm/_cacache/`. The host assembles the chain; the
core matches it, and the slashes are what makes the match a match of whole
COMPONENTS: `/home/u/x.npm/_cacache/` does not contain `/home/u/.npm/_cacache/`,
because there is no slash before `.npm`. A справочник whose chains are not
bounded is refused whole («Справочник ограничен») — matching one as a bare
substring would bring back the bug fixed on 1 September.

**What it may not do.** It names a разряд and nothing else, and only four of
them: Кэш, Журнал, Сборка, Загрузка. «Крупное» is decided by size and
«Неизвестное» means "no place matched", and the core refuses to let a file
assert either (постусловие «Место обосновано»). It softens no threshold, removes
no directory and no symlink, and a content-addressed store stays untouched:
постусловие И3 outranks any line of the file. Invariant И1 — «МожноУбрать»
never leaves Кэш, Журнал and Сборка — holds exactly as before, however many
places the file knows.

<a id="the-protection-list"></a>

### The защитный список

How to say "do not touch this", by path and by разряд:

```bash
./digitdisk clean <path> --protect ~/projects        # the path and everything under it
./digitdisk clean <path> --protect разряд:Загрузка   # a whole разряд
./digitdisk clean <path> --protect-file FILE         # a list from a file
```

Without a flag, `~/.digitable/digitdisk/protect.conf` is read; a row there is
`путь|~/projects|why` or `разряд|Журнал|why`. A path written without a leading
slash protects that chain of components at any depth.

The защитный список lives **in the host and not in the rules**, and that is not
an implementation detail. The core answers one question — what this path IS —
and every answer it gives is proved; "do not touch my `~/projects`" is not an
answer to that question: the path may very well be a cache, and writing the
opposite into the справочник would be putting a falsehood into the layer to get
an effect. An instruction from the person who owns the machine belongs where the
host already keeps its veto — next to the checks in `internal/clean`. That is
why it weakens nothing: the list can only subtract from a plan, and no
постусловие of the core moves because of it. What it protected is printed in its
own ЗАЩИЩЕНО section, with the rule and the file line, rather than quietly
missing from the plan — and it is kept apart from ОТКАЗАНО, because a refusal
means the two layers disagree and somebody should look at the rules, while a
protection means the rules worked and a person overruled the answer.

### One home for the settings

```
~/.digitable/digitdisk/settings.conf   the language, and nothing else
~/.digitable/digitdisk/places.conf     a справочник of one's own
~/.digitable/digitdisk/protect.conf    the защитный список
```

There was one home already — `~/.config/digitdisk/` — and the language would
have made a second. Two homes are two places to look for one answer, and every
document would then have to say which of them holds what. So there is one, and
it is not digitdisk's alone: `~/.digitable/` is the family's, and the tools
beside this one keep their settings beside it.

The old home is not broken and not deleted. `~/.config/digitdisk/places.conf`
and `~/.config/digitdisk/protect.conf` are still READ where they are; a run
that takes a file from there says so once and not at every start; and nothing
is copied on anybody's behalf, because a tool that writes into a person's home
unasked is the thing this tool exists to clean up after. A file lying in both
homes is read from the new one: it was moved, and the copy left behind is not
the one that was meant.

### The cleanup journal

```bash
./digitdisk history <path>          # a cleaned root, the корзина store, or one корзина
./digitdisk history <path> --json
```

What was removed, when, how many bytes are sitting in корзины, how much went
back, how much was erased — and what puts the last one back. **digitdisk
remembers nothing between runs**: every number is read out of the same
`journal.json` files that `restore` and `purge` obey. A separate history
database would be a second account of the same events, and the two would
disagree the first time somebody moved a корзина with `mv`.

"Freed" in that summary counts only what was erased: moving into a корзина frees
no bytes at all, and a number claiming otherwise would be a lie about the disk.

<a id="why-not-the-system-trash"></a>

### Why not the system Trash

What is cleaned goes into digitdisk's own корзина inside the tree, not into the
desktop Trash, for three reasons, one of which is a number.

**The number.** A корзина inside the tree means `rename(2)`: on this machine
moving a gibibyte does not register on the timer at all (0.00 s, three runs),
because no bytes move. A корзина across a filesystem boundary turns the move
into a copy: the same gibibyte written to disk with `fsync` takes **0.91 s per
GiB** (best of three: 0.91 / 1.01 / 1.12). The cost of reversibility would
become the size of the cleanup, the file would exist twice while the copy runs
(so the space has to be free beforehand), and a crash halfway would leave half a
file. Both `~/.local/share/Trash` on Linux and `~/.Trash` on macOS live in the
home directory, and cleaning usually happens on other volumes.

**The write boundary.** "Does not leave the tree you named" is a property of the
system calls digitdisk uses: everything goes through an `os.Root` opened on the
root, which cannot be walked out of even through a symlink. A cleanup of
`/var/tmp` that writes into `~/.local/share/Trash` cancels that property.

**The two systems have no common behaviour.** On Linux the Trash is defined by
the [freedesktop.org specification](https://specifications.freedesktop.org/trash-spec/latest/):
`files/` and `info/` with `.trashinfo` records — and the same specification says
the home Trash only accepts files from its own filesystem, while another volume
needs a `.Trash-$uid` at its top level, which digitdisk would have to create. On
macOS the layout is different and Finder's "Put Back" lives in an unpublished
store: the documented way in is `NSFileManager trashItemAtURL:`, which is Cocoa,
which is cgo, which is the end of cross-building four targets from one machine
with a reproducible digest. There is no common behaviour to implement here —
there are two different Trashes and one way to lie about freed space.

What there is instead: digitdisk's корзина is an ordinary directory. Whoever
wants to hand it to the system Trash hands it over themselves, in one gesture,
and knows they did.

### Hardware: the mark of the system, the cores one by one, and the cards

**СИСТЕМА — the mark, and what a machine is recognised by.** A drawing on the
left, and on the right what a person wants to know about their own machine:
node (`user@host`), distribution, the machine's model, kernel and word size,
shell, desktop, terminal, uptime, the processor on one line, the memory on one
line, the video cards on one line. The model is what the firmware calls the
machine — `/sys/class/dmi/id/sys_vendor` and `product_name` on Linux,
`hw.model` on macOS; the processor is the `model name` line of `/proc/cpuinfo`,
or `machdep.cpu.brand_string` on macOS. The marks were drawn in this tree and
nowhere else: somebody else's collection is somebody else's work under
somebody else's licence. They are drawn in printable ASCII, so they hold
together in a font without our glyphs and under `LANG=C`; none is wider than
fourteen columns, so a wide terminal puts the mark beside the fields and a
narrow one above them. A distribution nobody drew gets the general mark rather
than an empty space, and a family counts as a family: Rocky, Alma and CentOS
take the RHEL mark.

**ЗАГРУЗКА — every core of it.** "Занято ЦП" is one number for the whole
machine: on a machine with 256 cores it says "8%" both when the load is spread
and when one core is on fire and the rest are asleep. So the cores are drawn
underneath it, and the screen picks how: while the gauges fit the height, every
core gets its own gauge in columns; when they stop fitting, a map where one
cell is one core, plus the list of the busiest. On 256 cores the map is 4 rows
of 64 cells at 80 columns, and at 200 columns all 256 gauges fit instead, in 26
rows. The printed report gets one line of it: minimum, median, maximum, the
number of the busiest core, and how many cores are busy more than half the
time. The source is `/proc/stat` line by line on Linux and
`host_processor_info(PROCESSOR_CPU_LOAD_INFO)` on macOS — and on both the list
is published only if the mean of the cores comes out as the machine-wide share,
which is the sum of those same counters.

**ВИДЕОКАРТЫ — a section of its own, and there may be several cards.** Name,
busy share, memory used out of total, temperature, power and clock, for each
card. The cards come from files: `/sys/class/drm`, the display-class devices of
the PCI bus (a card with no driver is still a card), and
`/proc/driver/nvidia`. The name comes from the driver, and where the driver is
silent, from the `pci.ids` database the distribution ships; nothing of it is
copied into this tree. What is shown is what the driver published:

| driver | what it gives in files |
|---|---|
| `amdgpu` | load, memory, temperature, clock, power |
| `i915`, `xe` | temperature and power on the newer chips; no busy share |
| `nvidia` | the name, the bus and the firmware versions — and not one counter |
| `mgag200` and its kind | the name and nothing else |

`--gpu-tool` allows asking `nvidia-smi` — **somebody else's program**, not a
file. Without the key it is never run; with it, every card says underneath
where its numbers came from: "числа из /sys/class/drm/card1/device" or "числа
от чужой программы nvidia-smi". A row about a card the files never saw is
thrown away: a program cannot add hardware to a machine. Power is read as the
hwmon documentation defines it, in microwatts, and printed only if the result
is at least half a watt: some drivers count in something else, and a number
without a unit is not a number.

**On macOS there are no video cards in the snapshot.** What a Mac knows about
its graphics lives in the IORegistry, and the documented way in is IOKit —
Core Foundation objects rather than numbers. We do not read those without cgo,
and we will not guess. The reason is behind `--why`; the field is empty.

### The live screen

In a terminal, `digitdisk` and `digitdisk status` open a live screen in the
Digitable Focus palette: the sections of the printed report as pages that keep
measuring themselves. `← →` and `Tab` move between them, `1`…`9` go straight to
one, `↑ ↓` scroll a long one, `p` holds, `r` measures now, `l` switches the
language, `?` lists the subcommands, `q` leaves. There are ten sections; the
digits reach the first nine, and the tenth — НЕ ПРОЧИТАНО — sits to the left of
the first, one `←` away from ОБЗОР.


`l` is the one key here that touches anything outside the screen: it turns the
whole report into the other language where the reader is looking at it, and
puts the new choice into `settings.conf`, so that the next run — and `digitdisk
clean` tomorrow — speaks the same language. It says which file it wrote, on the
screen, for the six seconds after; a program that silently rewrites a file in a
home directory is the thing this tool is for cleaning up after.

The `?` list only names the commands; it runs none of them. This screen is
`status`, which reads and writes nothing, and `clean` moves files — a command
that changes the disk is not put one keystroke away from a reading.

Everywhere else it prints, exactly as it always has. A pipe, a file,
`/dev/null`, `--json`, `TERM=dumb` and an empty `TERM` all receive the text
report: the screen is never drawn into something that is not a terminal, so
scripts see what they have always seen.

| | |
|---|---|
| `--plain` | print the snapshot once, even in a terminal |
| `--live` | demand the screen; fail rather than print if there is none |
| `--interval MS` | how often the screen measures again (default 2000) |
| `DIGITDISK_PALETTE` | `carbon` (default), `paper`, `signal` — the palettes of the stack |

`NO_COLOR` is honoured: the screen still runs, it is simply drawn without
colour.

## What it does not do

- **It does not delete by pattern, and it does not delete from a list of
  paths.** digitdisk does have a list of known places, and that list removes
  nothing: the справочник names a РАЗРЯД, and what `clean` removes is still
  exactly what the decision layer gave the verdict «МожноУбрать» — with the same
  thresholds, the same refusal to touch directories and symlinks, and the same
  refusal to touch content-addressed stores. A place is an argument, not an
  order. The host keeps a veto on top of that and refuses a directory, a symlink
  or anything unreadable even if the layer were to ask — and when the two
  disagree it prints the disagreement instead of acting on it.
- **It never deletes in one step, and never without being asked.** There is no
  flag that erases without a plan first and a separate confirmation after, and
  `clean` on its own touches nothing at all.
- **It does not leave the tree you named.** Every path operation goes through
  `os.Root` opened on that directory: it resolves each component itself and
  cannot be walked out of, symbolic links included. The корзина must live inside
  the same tree — a корзина elsewhere would make every move a cross-filesystem
  copy, and the cost of reversibility would become the size of the cleanup.
- **It does not delete recursively.** `os.RemoveAll` appears nowhere in this
  tree and `tools/licensing.flang` fails the build if it ever does. Files go
  one at a time, from a list in a journal; empty directories go through the
  call that refuses a directory with anything in it.
- **It does not explain instead of measuring.** Where there is a number, it is
  printed; where there is none, a dash, and the name of the reading on one line
  at the end. Why it is missing lives behind `digitdisk status --why`, a flag of
  its own, and not in the middle of the report: a reader wants a number, not an
  essay about kernel calls.
- **It is not a fork of mole**, and carries none of its GPL-3.0 code — the idea
  came from there, the code did not. See [`NOTICE`](NOTICE).

## Documents

| Document | What is in it |
|---|---|
| [`LICENSE`](LICENSE) | the binding text: BSD-2-Clause, verbatim |
| [`LICENSE-RU.md`](LICENSE-RU.md) | what that licence means, in plain Russian |
| [`NOTICE`](NOTICE) | where the idea came from, what was deliberately not taken, and why |
| [`AGENTS.md`](AGENTS.md) | the rules of this tree: write boundary, no GPL, where removal may live, the order of the checks |
| [`digitdisk.en.1`](digitdisk.en.1) | the manual page: subcommands, every flag, files, examples, exit codes |
| [`digitdisk.1`](digitdisk.1) | the same page in Russian; the formula puts it where `man` looks for a translation |

## Checks

```bash
flang io tools/licensing.flang      # no copyleft; SPDX headers; removal only in host/internal/clean
flang check core/disk-inventory.flang && flang test core/disk-inventory.flang
make -C core                        # check, emit to Go and C, cross-check the two emissions
cd host && go vet ./... && go test -count=1 ./...
cd host && go test ./internal/lang/  # every line a person reads has a pair
cd host && GOOS=darwin GOARCH=arm64 go build ./... && GOOS=darwin GOARCH=amd64 go build ./...
cd host && GOOS=darwin go vet ./...  # the macOS host, checked from a machine that is not one
scripts/build-release.sh            # release archives, sums, formula; verifies the build repeats
```

The translation is checked by a run and not promised in a document.
`go test ./internal/lang/` reads the host's own source, finds every line that
reaches a person, and fails on one with no pair in the other language; it fails
separately on Cyrillic printed past the dictionary out of `main`, `report`,
`ui` or `cli`, on a dictionary entry nobody ever asks for, and on
`%`-placeholders that disagree between the two halves of one entry. What it
covered is printed by the run itself —

```bash
cd host && go test ./internal/lang/ -v -run 'Пары|Заполнители|Договор'
```

— which at the moment reports 419 lines in the source, 488 entries in the
dictionary, and 29 names of the договор translated as words.

The licensing guard and the emission cross-check are written in flang, not in
Python or JavaScript: neither is present in this tree. The flang compiler is a
single binary that needs only a C compiler (`brew install flang`, `asdf`, or
`make -C bootstrap` in a clone of the language); it does not require Node.

## State

The tree is complete and installable: the licences and the gate, the flang core
printed into `core/out-go`, the Go host with `status`, `analyze`, `places`,
`history` and the three steps of `clean` / `restore` / `purge`, in Russian and
in English, and the release path — [`scripts/build-release.sh`](scripts/build-release.sh),
the Homebrew formula with both manual pages, and the tag-driven workflow in
[`.github/workflows/release.yml`](.github/workflows/release.yml). The version
lives in one place, [`VERSION`](VERSION); the build stamps it into the binary,
and the workflow refuses a tag that disagrees with it.
