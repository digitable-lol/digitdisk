// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package gpuinfo

// Two ways to turn the numbers a card answers with into the name a person
// calls it by, both of them files.
//
// The first is the database of PCI names that distributions ship as
// /usr/share/hwdata/pci.ids — a text file, kept by the pciutils project, in
// the public domain and under BSD-3-Clause at the same time.  It is read, not
// copied: nothing of it enters this tree, and a machine without it simply gets
// a card named by its numbers.
//
// The second is the directory the NVIDIA driver publishes for each of its
// cards, which spells the model out in the words the box was sold under.

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

// NvidiaInformation is what /proc/driver/nvidia/gpus/<bus>/information says
// about one card.  The file is a list of "Ключ: значение" lines and nothing
// else; the driver has published it in this shape for as long as it has
// published anything.
type NvidiaInformation struct {
	Model    string
	Bus      string
	VBIOS    string
	Firmware string
	UUID     string
}

// ParseNvidiaInformation reads that file.  A line it does not know is
// skipped: the driver adds lines between releases, and a new one must not
// cost us the model name.
func ParseNvidiaInformation(text string) NvidiaInformation {
	var out NvidiaInformation
	for _, line := range strings.Split(text, "\n") {
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		val = strings.TrimSpace(val)
		if val == "" {
			continue
		}
		switch strings.TrimSpace(key) {
		case "Model":
			out.Model = val
		case "Bus Location":
			out.Bus = strings.ToLower(val)
		case "Video BIOS":
			out.VBIOS = val
		case "GPU Firmware":
			out.Firmware = val
		case "GPU UUID":
			out.UUID = val
		}
	}
	return out
}

// pciName is what the database calls one card.
type pciName struct{ vendor, device string }

// namesCache remembers what was already looked up.  The live screen takes a
// snapshot every couple of seconds and the database is a megabyte and a half
// of text; a card's name, unlike its temperature, does not change between two
// readings.
var namesCache = struct {
	sync.Mutex
	m map[string]pciName
}{m: map[string]pciName{}}

// lookupPCIIDs answers the vendor and device names for the pairs asked about.
// A pair nobody could name is simply absent from the answer.
func (r Reader) lookupPCIIDs(want [][2]string) map[string]pciName {
	out := make(map[string]pciName, len(want))
	need := map[string][2]string{}
	namesCache.Lock()
	for _, pair := range want {
		key := pair[0] + ":" + pair[1]
		if n, ok := namesCache.m[key]; ok {
			out[key] = n
			continue
		}
		need[key] = pair
	}
	namesCache.Unlock()
	if len(need) == 0 {
		return out
	}

	for _, path := range r.IDs {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		found := scanPCIIDs(f, need)
		f.Close()
		if len(found) == 0 {
			continue
		}
		namesCache.Lock()
		for key, n := range found {
			namesCache.m[key] = n
			out[key] = n
		}
		namesCache.Unlock()
		break
	}
	return out
}

// scanPCIIDs walks the database once, looking only for the pairs asked about.
//
// The format is two levels of indentation and nothing else: a vendor at the
// left margin, its devices one tab in, a device's subsystems two tabs in.  The
// subsystem lines are skipped — they name a board built around the chip, and
// two of them can carry the same words for different cards.
func scanPCIIDs(f *os.File, need map[string][2]string) map[string]pciName {
	byVendor := map[string]bool{}
	for _, pair := range need {
		byVendor[pair[0]] = true
	}
	out := map[string]pciName{}
	// The vendor names of the vendors asked about, kept as they go by: the
	// device line that would carry the name may not be in the file at all,
	// and the vendor's own line has passed by then.
	vendorNames := map[string]string{}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	vendorID, vendorName := "", ""
	for sc.Scan() {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// The device classes follow the vendors in the same file, under
		// lines that begin with "C ".  Nothing after that is a vendor.
		if strings.HasPrefix(line, "C ") {
			break
		}
		switch {
		case !strings.HasPrefix(line, "\t"):
			id, name, ok := strings.Cut(line, "  ")
			if !ok {
				continue
			}
			vendorID, vendorName = strings.ToLower(strings.TrimSpace(id)), strings.TrimSpace(name)
			if byVendor[vendorID] {
				vendorNames[vendorID] = vendorName
			}
		case !strings.HasPrefix(line, "\t\t"):
			if !byVendor[vendorID] {
				continue
			}
			id, name, ok := strings.Cut(strings.TrimPrefix(line, "\t"), "  ")
			if !ok {
				continue
			}
			key := vendorID + ":" + strings.ToLower(strings.TrimSpace(id))
			if _, ok := need[key]; !ok {
				continue
			}
			out[key] = pciName{vendor: vendorName, device: strings.TrimSpace(name)}
			if len(out) == len(need) {
				return out
			}
		}
	}
	// A vendor with no matching device still names the vendor, which is
	// better than a card called by two hexadecimal words alone.
	for key, pair := range need {
		if _, ok := out[key]; ok {
			continue
		}
		if name := vendorNames[pair[0]]; name != "" {
			out[key] = pciName{vendor: name}
		}
	}
	return out
}
