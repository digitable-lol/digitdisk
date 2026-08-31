**English** · [Русский](README.ru.md)

# digitdisk — where the disk went and how the machine feels, read-only

digitdisk prints two readings of a machine and changes nothing while it does:
**where the disk space went** — directories by size and the largest files — and
**how the machine is feeling right now** — CPU, memory, disk and network. It
reads, prints, and exits. There is no cleaning mode, no uninstaller, and no flag
that removes a file.

## What it is made of

```
core/     the readings as a flang specification, plus the Go printed from it into core/out-go
host/     the Go host: system calls, the command line, the output
tools/    the licence gate
docs/     notes that are not this page
```

The split is the whole design. Everything that can be decided without touching
the operating system is decided in `core/`: a specification in flang, checked by
flang's own runs and **printed** into Go. Everything that must touch the
operating system — walking directories, `statfs`, reading counters — is
hand-written Go in `host/`. The core never opens a file; the host never decides
what a number means.

`core/out-go` is printed, not written. Hand edits there are lost at the next
print — see [`AGENTS.md`](AGENTS.md).

## Build

A Go toolchain is enough for a build from a clean checkout, because the printed
Go is committed. Reprinting the core additionally needs
[flang](https://github.com/digitable-lol/flang).

```bash
flang emit core/<module>.flang --target go --out core/out-go   # only when the specification changed
go build -o digitdisk ./host
```

## Run

```bash
./digitdisk disk [path]    # where the space went: directories by size, the largest files
./digitdisk status         # how the machine feels: CPU, memory, disk, network
```

Both commands read. Neither writes.

## What it does not do

- **It never deletes anything.** No cleaning, no uninstalling, no `--force`, no
  quarantine directory. digitdisk prints what it found; what to remove is the
  reader's decision and somebody else's command. A patch that adds a delete path
  is out of scope, not a feature request.
- **macOS is not supported.** The host speaks Linux and the BSDs. macOS is not a
  build target, not a test target, and not a bug report we can act on.
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
python3 tools/check-licensing.py    # no GPL in the tree; our files carry SPDX BSD-2-Clause
```

## State

The root of the tree — licences, rules, the gate — is complete. `core/` and
`host/` are written separately and are not here yet.
