// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package darwinsys

// KinfoProcSize is sizeof(struct kinfo_proc) on 64-bit macOS: struct
// extern_proc (296 bytes, <sys/proc.h>) followed by struct eproc (352 bytes,
// <sys/sysctl.h>).  kern.proc.all answers with an array of these, and the
// array is the only way to enumerate processes without libproc, which needs
// cgo.
const KinfoProcSize = 648

// Byte offsets inside one kinfo_proc record.  Each is the sum of the field
// sizes above it in the two structures named at KinfoProcSize; the comment
// gives the field and the structure it lives in.
//
// This is the one place in the tree that depends on a C structure layout we
// cannot compile against, which is why ParseProcs hands its result to Verify
// before anybody believes it.
const (
	offPStat  = 36  // extern_proc.p_stat,   char
	offPID    = 40  // extern_proc.p_pid,    pid_t
	offPctCPU = 88  // extern_proc.p_pctcpu, fixpt_t
	offComm   = 243 // extern_proc.p_comm,   char[MAXCOMLEN+1]
	lenComm   = 17
	offUID    = 420 // eproc.e_ucred.cr_uid, uid_t  (296 + 120 + 4)
	offPPID   = 560 // eproc.e_ppid,         pid_t  (296 + 264)
)

// FSCALE is the fixed-point scale of fixpt_t (1 << FSHIFT, <sys/param.h>).
const FSCALE = 2048

// Process states, from <sys/proc.h>.  They are rendered as the letters
// /proc/<pid>/stat uses for the same states, so one report reads the same on
// both systems.
const (
	stateIdle     = 1 // SIDL
	stateRun      = 2 // SRUN
	stateSleep    = 3 // SSLEEP
	stateStopped  = 4 // SSTOP
	stateZombie   = 5 // SZOMB
	StateRunning  = "R"
	StateSleeping = "S"
	StateStopped  = "T"
	StateZombie   = "Z"
	StateIdle     = "I"
	StateUnknown  = "?"
)

// Proc is one process as kern.proc.all describes it.
//
// It is short because kinfo_proc is short.  Resident memory is not in it at
// all: that lives in libproc (proc_pidinfo) and in the Mach task port, neither
// of which Go reaches without cgo, so the host leaves it unmeasured rather
// than filling it with zeros.
type Proc struct {
	PID   int
	PPID  int
	UID   int
	State string
	Comm  string
	// PctCPU is p_pctcpu turned into a percentage of one core.
	//
	// The field is a leftover of BSD process accounting, and a kernel that
	// no longer keeps it simply leaves it zero for every process.  The
	// caller must therefore treat a whole list of zeros as "not published"
	// rather than as "every process is idle" — see AnyCPU.
	PctCPU float64
}

// AnyCPU reports whether any process in the list carries a non-zero p_pctcpu,
// which is how the caller tells a kernel that fills the field from one that
// does not.
func AnyCPU(procs []Proc) bool {
	for _, p := range procs {
		if p.PctCPU > 0 {
			return true
		}
	}
	return false
}

// ParseProcs decodes the kern.proc.all buffer into one Proc per record.
//
// The buffer may be one byte short of a whole number of records — see Padded
// for why — and is restored before the split.  Anything else that does not
// divide into records is refused: a partial record means the record is not the
// size we think it is.
func ParseProcs(b []byte) ([]Proc, bool) {
	if len(b)%KinfoProcSize == KinfoProcSize-1 {
		b = append(append(make([]byte, 0, len(b)+1), b...), 0)
	}
	if len(b) == 0 || len(b)%KinfoProcSize != 0 {
		return nil, false
	}
	out := make([]Proc, 0, len(b)/KinfoProcSize)
	for off := 0; off+KinfoProcSize <= len(b); off += KinfoProcSize {
		rec := b[off : off+KinfoProcSize]
		out = append(out, Proc{
			PID:    int(i32(rec, offPID)),
			PPID:   int(i32(rec, offPPID)),
			UID:    int(u32(rec, offUID)),
			State:  stateLetter(rec[offPStat]),
			Comm:   cstring(rec[offComm : offComm+lenComm]),
			PctCPU: 100 * float64(u32(rec, offPctCPU)) / FSCALE,
		})
	}
	return out, true
}

func stateLetter(c byte) string {
	switch c {
	case stateRun:
		return StateRunning
	case stateSleep:
		return StateSleeping
	case stateStopped:
		return StateStopped
	case stateZombie:
		return StateZombie
	case stateIdle:
		return StateIdle
	}
	return StateUnknown
}

// Verify reports whether a decoded list really describes the machine it was
// read on: the process that asked must be in it, carrying the parent and the
// user it knows it has.
//
// This is the runtime proof that the offsets above still match the kernel's
// idea of kinfo_proc.  Three fields spread over 520 bytes — the pid at 40, the
// uid at 420, the parent at 560 — all landing on values the caller already
// knows cannot happen by accident, while any one of them alone could: a
// shifted layout still finds some plausible small integer where a pid belongs.
// When Verify says no, the caller publishes no process list at all and names
// the reason — an empty section is honest, a list of misread numbers is not.
func Verify(procs []Proc, pid, ppid, uid int) bool {
	for _, p := range procs {
		if p.PID == pid {
			return p.PPID == ppid && p.UID == uid
		}
	}
	return false
}
