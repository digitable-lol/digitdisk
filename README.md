**English** · [Русский](README.ru.md)

# digitdisk — where the disk went and how the machine feels, read-only

digitdisk prints two readings of a machine and changes nothing while it does:
**where the disk space went** — directories by size and the largest files — and
**how the machine is feeling right now** — CPU, memory, disk and network. It
reads, prints, and exits. There is no cleaning mode, no uninstaller, and no flag
that removes a file.

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

**Linux only, x86-64 and arm64.** The host reads `/proc` and `/sys` and calls
`uname(2)` and `statfs(2)` in the shape the Linux kernel gives them. macOS has
no `/proc`, and the host does not compile for it at all — there is no macOS
build to install, and that is a design decision, not a gap.

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

Both readings take `--json`. Both commands read. Neither writes.

## What it does not do

- **It never deletes anything.** No cleaning, no uninstalling, no `--force`, no
  quarantine directory. digitdisk prints what it found; what to remove is the
  reader's decision and somebody else's command. A patch that adds a delete path
  is out of scope, not a feature request.
- **macOS is not supported.** The host reads `/proc` and `/sys`; macOS has
  neither, and `uname(2)` in the form the host calls it does not exist there,
  so the host does not even compile for macOS. Not a build target, not a test
  target, and not a bug report we can act on.
- **It is not a fork of mole**, and carries none of its GPL-3.0 code — the idea
  came from there, the code did not. See [`NOTICE`](NOTICE).

## Documents

| Document | What is in it |
|---|---|
| [`LICENSE`](LICENSE) | the binding text: BSD-2-Clause, verbatim |
| [`LICENSE-RU.md`](LICENSE-RU.md) | what that licence means, in plain Russian |
| [`NOTICE`](NOTICE) | where the idea came from, what was deliberately not taken, and why |
| [`AGENTS.md`](AGENTS.md) | the rules of this tree: write boundary, no GPL, the order of the checks |

## Checks

```bash
python3 tools/check-licensing.py    # no copyleft in the tree; our files carry SPDX BSD-2-Clause
cd host && go vet ./... && go test ./...
scripts/build-release.sh            # release archives, sums, formula; verifies the build repeats
```

## State

The tree is complete and installable: the licences and the gate, the flang core
printed into `core/out-go`, the Go host with `status` and `analyze`, and the
release path — [`scripts/build-release.sh`](scripts/build-release.sh), the
Homebrew formula, and the tag-driven workflow in
[`.github/workflows/release.yml`](.github/workflows/release.yml). The version
lives in one place, [`VERSION`](VERSION); the build stamps it into the binary,
and the workflow refuses a tag that disagrees with it.
