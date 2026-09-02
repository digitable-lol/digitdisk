**English** · [Русский](README.ru.md)

# digitdisk — where the disk went, how the machine feels, and what may go

digitdisk prints two readings of a machine — **where the disk space went**,
directories by size and the largest files, and **how the machine is feeling
right now**, CPU, memory, disk and network — and it can act on the first of
them: `clean` removes files, in three steps, none of which is a surprise.

`status` and `analyze` read and write nothing. `clean` shows a plan and needs
`--apply` to move anything; what it moves goes to a корзина inside the tree you
named and comes back with `restore`; erasing is a separate command with a
separate confirmation. **What may be removed is not a list of well-known paths
and not a pattern**: it is exactly what the decision layer in `core/` gives the
verdict «МожноУбрать», which is proved to be nothing outside Кэш, Журнал and
Сборка.

## What it is made of

```
core/       the readings as a flang specification, plus the Go printed from it into core/out-go
host/       the Go host: system calls, the command line, the output
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
./digitdisk analyze <path>   # where the space went: directories by size, the largest files
./digitdisk status           # how the machine feels: CPU, memory, disk, network
./digitdisk --version        # version, build hash, toolchain, decision layer
```

Both readings take `--json`. Neither writes anything.

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

### The live screen

In a terminal, `digitdisk status` opens a live screen in the Digitable Focus
palette: the sections of the printed report as pages that keep measuring
themselves. `← →` and `Tab` move between them, `1`…`9` go straight to one,
`↑ ↓` scroll a long one, `p` holds, `r` measures now, `q` leaves.

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

- **It does not delete by pattern, by name, or by a list of known paths.** The
  only thing `clean` will remove is a path the decision layer gave the verdict
  «МожноУбрать». The host keeps a veto on top of that and refuses a directory, a
  symlink or anything unreadable even if the layer were to ask — and when the
  two disagree it prints the disagreement instead of acting on it.
- **It never deletes in one step, and never without being asked.** There is no
  flag that erases without a plan first and a separate confirmation after, and
  `clean` on its own touches nothing at all.
- **It does not leave the tree you named.** Every path operation goes through
  `os.Root` opened on that directory: it resolves each component itself and
  cannot be walked out of, symbolic links included. The корзина must live inside
  the same tree — a корзина elsewhere would make every move a cross-filesystem
  copy, and the cost of reversibility would become the size of the cleanup.
- **It does not delete recursively.** `os.RemoveAll` appears nowhere in this
  tree and `tools/check-licensing.py` fails the build if it ever does. Files go
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

## Checks

```bash
python3 tools/check-licensing.py    # no copyleft; SPDX headers; removal only in host/internal/clean
cd host && go vet ./... && go test ./...
cd host && GOOS=darwin GOARCH=arm64 go build ./... && GOOS=darwin GOARCH=amd64 go build ./...
cd host && GOOS=darwin go vet ./...  # the macOS host, checked from a machine that is not one
scripts/build-release.sh            # release archives, sums, formula; verifies the build repeats
```

## State

The tree is complete and installable: the licences and the gate, the flang core
printed into `core/out-go`, the Go host with `status`, `analyze` and the three
steps of `clean` / `restore` / `purge`, and the release path — [`scripts/build-release.sh`](scripts/build-release.sh), the
Homebrew formula, and the tag-driven workflow in
[`.github/workflows/release.yml`](.github/workflows/release.yml). The version
lives in one place, [`VERSION`](VERSION); the build stamps it into the binary,
and the workflow refuses a tag that disagrees with it.
