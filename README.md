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

**Released binaries: Linux, x86-64 and arm64.** On macOS the host builds from
source and runs — see [macOS](#macos) below — but nothing is released for it
yet, because nobody has run it on a Mac.

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

### macOS: builds from source, not released

The host builds and runs on macOS, arm64 and x86-64. It takes its facts from
`sysctl(3)`, `getfsstat(2)` and the routing socket instead of `/proc` and
`/sys`; the flang core is untouched by any of that, because the core has never
known what a system call is.

```bash
cd host && GOOS=darwin go build -tags flangcore -o ../digitdisk .
```

**No Mac has ever run this.** The authors have none. What is checked is that it
compiles and passes vet for both macOS architectures, and that every kernel
structure it decodes is covered by tests over samples the tests build
themselves. A snapshot of a real Mac has not been seen, which is why there is no
macOS archive in the release and the Homebrew formula still says
`depends_on :linux`.

What macOS measures, and the call each number comes from:

| Reading | Source |
|---|---|
| host, kernel, release, model | `sysctl` `kern.hostname`, `kern.osrelease`, `kern.version`, `kern.osproductversion`, `kern.osversion`, `hw.machine`, `hw.model` |
| uptime | `sysctl kern.boottime` (`struct timeval`) |
| load average, cores | `sysctl vm.loadavg` (`struct loadavg`), `hw.logicalcpu` |
| memory total, swap | `sysctl hw.memsize`, `vm.swapusage` (`struct xsw_usage`) |
| processes: total, running | `sysctl kern.proc.all` (`struct kinfo_proc`) |
| disks | `getfsstat(2)` (`struct statfs`) |
| interfaces, addresses | `net.Interfaces` |
| interface counters | `sysctl NET_RT_IFLIST2` (`struct if_data64`) |

What macOS does **not** give a program without cgo, and what digitdisk
therefore prints as "—" while naming it under `НЕ ИЗМЕРЕНО`: the CPU busy share
and per-process CPU time, the memory breakdown (free, active, inactive, wired,
compressed), per-process resident memory and command line, the thread count,
and temperatures. Those live in Mach calls (`host_statistics64`), in libproc
and in IOKit — reachable from Go only through cgo, and cgo would end the
repeatable cross-compiled release build. An empty field is honest; a zero would
not be.

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
- **It does not pretend macOS is finished.** The host builds and runs there,
  and a good part of a snapshot is out of reach without cgo and comes out
  empty; every empty field says which call it would have needed. No Mac has run
  it yet, so there is no macOS archive in the release and the formula installs
  on Linux only.
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
