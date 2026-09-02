// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package gpuinfo

// The one place in this package that runs somebody else's program.
//
// The NVIDIA driver keeps its counters to itself: /proc/driver/nvidia names
// the card and stops there, and no file under /sys carries its load, its
// memory or its temperature.  The only published way to those numbers is
// nvidia-smi, the program the driver ships — so on a machine with an NVIDIA
// card the choice is between a program we did not write and four empty
// fields.
//
// The choice is the reader's, not ours.  Nothing here runs unless
// `--gpu-tool` was given, every number that arrives this way is marked as
// having come from outside (Card.Outside, Card.Source), and a row about a
// card the files did not already find is thrown away — a program that tells
// us about hardware we cannot see is telling us about a different machine.
import (
	"context"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// nvidiaSMIQuery is what is asked for, in that order.  The format is the
// documented machine-readable one: comma-separated values, no header, no
// units — the units are fixed by the query and written down in the parser.
var nvidiaSMIQuery = []string{
	"--query-gpu=pci.bus_id,name,memory.used,memory.total,utilization.gpu,temperature.gpu,power.draw,clocks.sm",
	"--format=csv,noheader,nounits",
}

// SMIRow is one line of that answer, decoded.
type SMIRow struct {
	Bus      string
	Name     string
	UsedMiB  *float64
	TotalMiB *float64
	Busy     *float64
	Celsius  *float64
	Watts    *float64
	ClockMHz *float64
}

// askNvidiaSMI fills the cards the files could not fill.
func (r Reader) askNvidiaSMI(cards []Card) {
	if !r.wantsSMI(cards) {
		return
	}
	out, err := r.run("nvidia-smi", nvidiaSMIQuery...)
	if err != nil {
		return
	}
	for _, row := range ParseNvidiaSMI(string(out)) {
		i := indexByBus(cards, row.Bus)
		if i < 0 {
			// A card the files never saw.  It is not published: the
			// list of cards is what /sys and /proc say it is, and a
			// program cannot add hardware to a machine.
			continue
		}
		c := &cards[i]
		took := false
		if c.BusyPercent == nil && row.Busy != nil {
			c.BusyPercent, took = row.Busy, true
		}
		if c.Celsius == nil && row.Celsius != nil {
			c.Celsius, took = row.Celsius, true
		}
		if c.Watts == nil && row.Watts != nil {
			c.Watts, took = row.Watts, true
		}
		if c.MHz == nil && row.ClockMHz != nil {
			c.MHz, took = row.ClockMHz, true
		}
		if c.MemoryTotalBytes == nil && row.TotalMiB != nil {
			total := uint64(*row.TotalMiB) * 1024 * 1024
			c.MemoryTotalBytes, took = &total, true
			if row.UsedMiB != nil && *row.UsedMiB <= *row.TotalMiB {
				used := uint64(*row.UsedMiB) * 1024 * 1024
				c.MemoryUsedBytes = &used
			}
		}
		if c.Name == "" && row.Name != "" {
			c.Name = row.Name
		}
		if took {
			c.Outside = true
			c.Source = "nvidia-smi"
		}
	}
}

// wantsSMI reports whether any card is both NVIDIA's and short of numbers.
// A machine whose cards all answered from files never starts the program.
func (r Reader) wantsSMI(cards []Card) bool {
	for _, c := range cards {
		if c.Driver != "nvidia" && c.VendorID != "10de" {
			continue
		}
		if c.BusyPercent == nil || c.MemoryTotalBytes == nil || c.Celsius == nil {
			return true
		}
	}
	return false
}

// run starts the program, with a limit on how long it may take.  A driver in
// trouble makes nvidia-smi hang, and a live screen that hangs with it would
// be worse than a screen with a dash on it.
func (r Reader) run(name string, args ...string) ([]byte, error) {
	if r.Run != nil {
		return r.Run(name, args...)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, name, args...).Output()
}

// ParseNvidiaSMI decodes the comma-separated answer.
//
// Every field is optional: the program writes "[N/A]" for a card that does
// not measure something, and "[N/A]" must never become a zero.  A value
// outside what it claims to be — a share above a hundred, a temperature no
// silicon survives — is dropped the same way.
func ParseNvidiaSMI(text string) []SMIRow {
	var out []SMIRow
	for _, line := range strings.Split(text, "\n") {
		f := strings.Split(line, ",")
		if len(f) < 6 {
			continue
		}
		for i := range f {
			f[i] = strings.TrimSpace(f[i])
		}
		row := SMIRow{Bus: NormalizeBus(f[0]), Name: f[1]}
		if row.Bus == "" {
			continue
		}
		row.UsedMiB = smiNumber(f[2], 0, 1<<40)
		row.TotalMiB = smiNumber(f[3], 0, 1<<40)
		row.Busy = smiNumber(f[4], 0, 100)
		row.Celsius = smiNumber(f[5], -50, 200)
		if len(f) > 6 {
			row.Watts = smiNumber(f[6], 0, 10000)
		}
		if len(f) > 7 {
			row.ClockMHz = smiNumber(f[7], 0, 100000)
		}
		out = append(out, row)
	}
	return out
}

func smiNumber(s string, low, high float64) *float64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil || v < low || v > high {
		return nil
	}
	return &v
}

// NormalizeBus turns a PCI address into the spelling /sys uses.  nvidia-smi
// writes the domain in eight digits ("00000000:02:00.0"), the kernel in four
// ("0000:02:00.0"), and the two must meet somewhere for the answer to be
// joined to the card it belongs to.
func NormalizeBus(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return ""
	}
	domain := strings.TrimLeft(parts[0], "0")
	for len(domain) < 4 {
		domain = "0" + domain
	}
	return domain + ":" + parts[1] + ":" + parts[2]
}
