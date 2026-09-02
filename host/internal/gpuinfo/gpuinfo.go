// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package gpuinfo reads what a machine publishes about its video cards.
//
// Files first.  Everything here is an ordinary read of /sys and /proc: the
// cards themselves come from /sys/class/drm and /sys/bus/pci, the numbers from
// the files the driver publishes next to them, and the names from a database
// file the distribution ships.  Running somebody else's program to learn about
// our own machine is the last resort, not the first, and it never happens
// unless the reader was told it may (Reader.Tool) — see nvsmi.go.
//
// What a card publishes is entirely the driver's business, and the drivers do
// not agree:
//
//	amdgpu   gpu_busy_percent, mem_info_vram_used/total, hwmon (°C, Вт, МГц)
//	i915/xe  hwmon on the newer chips; no busy share, no memory of its own
//	nvidia   name, bus and firmware in /proc/driver/nvidia; no counters at all
//	mgag200  nothing but its name — a service processor's display controller
//
// So a field that one card fills is empty on the next, and an empty field is
// left empty: nil, printed as a dash.  Why it is empty is said once, in
// Result.NoNumbers, and the caller puts it where explanations belong — behind
// `--why`, never beside the number.
//
// The roots are fields so tests can point them at a captured tree instead of
// the live kernel, exactly as the Linux collector does.
package gpuinfo

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Card is one video card as the system describes it.
//
// Every measured number is a pointer: a card whose driver publishes no load
// share must not be printed as a card that is doing nothing.
type Card struct {
	// Node is the kernel's name for the card, "card0", when it has one.
	Node string `json:"node,omitempty"`
	// Bus is the PCI address, "0000:02:00.0", when the card is on PCI.
	Bus      string `json:"bus,omitempty"`
	VendorID string `json:"vendor_id,omitempty"`
	DeviceID string `json:"device_id,omitempty"`
	Vendor   string `json:"vendor,omitempty"`
	Name     string `json:"name"`
	Driver   string `json:"driver,omitempty"`

	MemoryTotalBytes *uint64  `json:"memory_total_bytes"`
	MemoryUsedBytes  *uint64  `json:"memory_used_bytes"`
	BusyPercent      *float64 `json:"busy_percent"`
	Celsius          *float64 `json:"celsius"`
	Watts            *float64 `json:"watts"`
	MHz              *float64 `json:"mhz"`

	// Source names where the numbers came from: a directory of the
	// filesystem, or the program that was run.  It travels with the card
	// into the JSON and onto the screen, because "99%" from a file the
	// driver publishes and "99%" from a program we started are not the same
	// claim.
	Source string `json:"source"`
	// Outside is true when a program outside digitdisk was run to fill this
	// card in.
	Outside bool `json:"outside_tool"`
}

// Result is the whole answer: the cards, and the reasons for what is not in
// them.
type Result struct {
	Cards []Card
	// NoCards says why the list is empty, when it is.
	NoCards string
	// NoNumbers says why the cards that were found carry no load or memory.
	NoNumbers string
	// NoPower says why a card that publishes a power file has no watts
	// beside it — see wattsFrom.
	NoPower string
}

// Reader reads the cards from a set of roots.
type Reader struct {
	Sys  string // usually /sys
	Proc string // usually /proc
	// IDs are the candidate paths of the pci.ids database shipped by the
	// distribution.  The first one that opens is used; none is required.
	IDs []string
	// Tool allows running an outside program (nvidia-smi) for the numbers
	// no file publishes.  Off unless the caller was told to turn it on.
	Tool bool
	// Run is how an outside program is run, injected so the fallback is
	// testable without one.  nil means the real one.
	Run func(name string, args ...string) ([]byte, error)
}

// New returns a Reader pointed at the live system.
func New() Reader {
	return Reader{Sys: "/sys", Proc: "/proc", IDs: DefaultPCIIDs}
}

// DefaultPCIIDs are where distributions put the PCI name database.  It is
// data, not a program: a text file mapping the numbers a card answers with to
// the name a person calls it by.
var DefaultPCIIDs = []string{
	"/usr/share/hwdata/pci.ids",
	"/usr/share/misc/pci.ids",
	"/usr/share/pci.ids",
	"/var/lib/pciutils/pci.ids",
}

// classDisplay is the PCI class of a display controller (base class 0x03).
const classDisplay = 0x03

// Read lists the cards of this machine.
func (r Reader) Read() Result {
	var out Result
	oddPower := false
	seen := map[string]bool{} // by bus address, to join the three sources

	for _, node := range r.drmNodes() {
		c := r.fromDevice(filepath.Join(r.Sys, "class/drm", node, "device"), &oddPower)
		c.Node = node
		if c.Bus != "" {
			seen[c.Bus] = true
		}
		out.Cards = append(out.Cards, c)
	}

	// A card can be on the bus with no DRM node of its own: bound to no
	// driver, or handed to a virtual machine.  It is still a video card the
	// owner paid for, so it is named even though nothing about it moves.
	for _, dir := range r.pciDisplay() {
		if seen[filepath.Base(dir)] {
			continue
		}
		c := r.fromDevice(dir, &oddPower)
		seen[c.Bus] = true
		out.Cards = append(out.Cards, c)
	}

	// The NVIDIA driver publishes a directory of its own, and it is the only
	// place that spells the model out.  It also lists cards whose DRM node
	// is switched off.
	for bus, info := range r.nvidiaCards() {
		if i := indexByBus(out.Cards, bus); i >= 0 {
			out.Cards[i].mergeNvidia(info)
			continue
		}
		c := Card{Bus: bus, Driver: "nvidia", Source: filepath.Join(r.Proc, "driver/nvidia/gpus", bus)}
		c.mergeNvidia(info)
		out.Cards = append(out.Cards, c)
	}

	sort.SliceStable(out.Cards, func(i, j int) bool {
		a, b := out.Cards[i], out.Cards[j]
		if (a.Node == "") != (b.Node == "") {
			return a.Node != ""
		}
		if a.Node != b.Node {
			return nodeIndex(a.Node) < nodeIndex(b.Node)
		}
		return a.Bus < b.Bus
	})

	r.name(out.Cards)
	if r.Tool {
		r.askNvidiaSMI(out.Cards)
	}

	if len(out.Cards) == 0 {
		out.NoCards = "в " + filepath.Join(r.Sys, "class/drm") + " и на шине PCI видеокарт не нашлось"
	} else if silent := silentDrivers(out.Cards); len(silent) > 0 {
		out.NoNumbers = mute(silent, r.Tool)
	}
	if oddPower {
		out.NoPower = "драйвер отдаёт мощность карты не в микроваттах, как обещает документация hwmon: " +
			"полученное число меньше полуватта, то есть единица измерения у него другая — такое число мы не печатаем"
	}
	return out
}

// silentDrivers names the drivers whose cards published not one number.  A
// machine can have both kinds at once — this one had an AMD card answering
// four questions and an NVIDIA card answering none — so the answer is a list
// of drivers and not a yes or no about the machine.
func silentDrivers(cards []Card) []string {
	drivers := map[string]bool{}
	for _, c := range cards {
		if c.BusyPercent != nil || c.MemoryTotalBytes != nil || c.Celsius != nil || c.Watts != nil {
			continue
		}
		name := c.Driver
		if name == "" {
			name = "без драйвера"
		}
		drivers[name] = true
	}
	out := make([]string, 0, len(drivers))
	for d := range drivers {
		out = append(out, d)
	}
	sort.Strings(out)
	return out
}

// mute is the one sentence that explains the cards with no numbers beside
// them.  It is written here and printed only by `--why`.
func mute(silent []string, tool bool) string {
	for _, d := range silent {
		if d == "nvidia" && !tool {
			return "драйвер nvidia не публикует ни загрузки, ни памяти, ни температуры карты в файлах — их знает только его собственная программа nvidia-smi, и запускается она лишь по ключу --gpu-tool"
		}
	}
	if len(silent) == 1 && silent[0] == "без драйвера" {
		return "у найденной карты нет драйвера, а без него ядро не публикует о ней ничего, кроме имени"
	}
	return "драйвер " + strings.Join(silent, ", ") + " не публикует ни загрузки, ни памяти карты в файлах ядра"
}

func indexByBus(cards []Card, bus string) int {
	for i, c := range cards {
		if c.Bus != "" && c.Bus == bus {
			return i
		}
	}
	return -1
}

// drmNodes lists the card directories of /sys/class/drm, in kernel order.
// The connectors that live beside them ("card0-DP-3") are not cards and are
// left out: their names carry a dash, a card's name never does.
func (r Reader) drmNodes() []string {
	entries, err := os.ReadDir(filepath.Join(r.Sys, "class/drm"))
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		name := e.Name()
		if !strings.HasPrefix(name, "card") || strings.Contains(name, "-") {
			continue
		}
		if _, err := strconv.Atoi(strings.TrimPrefix(name, "card")); err != nil {
			continue
		}
		out = append(out, name)
	}
	sort.Slice(out, func(i, j int) bool { return nodeIndex(out[i]) < nodeIndex(out[j]) })
	return out
}

func nodeIndex(node string) int {
	n, err := strconv.Atoi(strings.TrimPrefix(node, "card"))
	if err != nil {
		return 1 << 30
	}
	return n
}

// pciDisplay lists the PCI devices whose class says display controller.
func (r Reader) pciDisplay() []string {
	entries, err := os.ReadDir(filepath.Join(r.Sys, "bus/pci/devices"))
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		dir := filepath.Join(r.Sys, "bus/pci/devices", e.Name())
		class, err := strconv.ParseUint(strings.TrimPrefix(trimmed(dir, "class"), "0x"), 16, 32)
		if err != nil || class>>16 != classDisplay {
			continue
		}
		out = append(out, dir)
	}
	sort.Strings(out)
	return out
}

// fromDevice reads one device directory: what the card is, and whatever
// numbers its driver publishes there.
func (r Reader) fromDevice(dir string, oddPower *bool) Card {
	c := Card{Source: dir}
	for _, line := range strings.Split(read(filepath.Join(dir, "uevent")), "\n") {
		key, val, ok := strings.Cut(strings.TrimSpace(line), "=")
		if !ok {
			continue
		}
		switch key {
		case "DRIVER":
			c.Driver = val
		case "PCI_ID":
			if v, d, ok := strings.Cut(strings.ToLower(val), ":"); ok {
				c.VendorID, c.DeviceID = v, d
			}
		case "PCI_SLOT_NAME":
			c.Bus = val
		}
	}
	if c.VendorID == "" {
		c.VendorID = strings.TrimPrefix(trimmed(dir, "vendor"), "0x")
		c.DeviceID = strings.TrimPrefix(trimmed(dir, "device"), "0x")
	}
	if c.Bus == "" {
		if target, err := filepath.EvalSymlinks(dir); err == nil {
			c.Bus = filepath.Base(target)
		}
	}
	if c.Driver == "" {
		if target, err := os.Readlink(filepath.Join(dir, "driver")); err == nil {
			c.Driver = filepath.Base(target)
		}
	}
	r.numbers(dir, &c, oddPower)
	return c
}

// numbers fills in whatever the driver publishes beside the card.  Each value
// is checked against what it claims to be — a share is a share, a temperature
// is a temperature — and refused rather than believed when it is not.
func (r Reader) numbers(dir string, c *Card, oddPower *bool) {
	if v, ok := number(trimmed(dir, "gpu_busy_percent")); ok && v >= 0 && v <= 100 {
		c.BusyPercent = &v
	}
	// The names differ by driver: amdgpu counts its video memory,
	// i915 counts local memory on the cards that have any.
	for _, pair := range [][2]string{
		{"mem_info_vram_total", "mem_info_vram_used"},
		{"lmem_total_bytes", "lmem_avail_bytes"},
	} {
		total, ok := unsigned(trimmed(dir, pair[0]))
		if !ok || total == 0 {
			continue
		}
		c.MemoryTotalBytes = &total
		if used, ok := unsigned(trimmed(dir, pair[1])); ok && used <= total {
			if pair[0] == "lmem_total_bytes" {
				used = total - used // that file counts what is free
			}
			c.MemoryUsedBytes = &used
		}
		break
	}
	r.hwmon(dir, c, oddPower)
}

// wattsFrom turns a hwmon power reading into watts, or refuses it.
//
// The documented unit is the microwatt: Documentation/hwmon/amdgpu.rst says
// power1_input is "instantaneous power used by the GPU in microWatts", and
// that is what is divided by here.  But a card that answers 18000 — which is
// eighteen milliwatts by that unit, and eighteen watts by the one some
// drivers actually use for integrated graphics — is a card whose unit we do
// not know, and a reading whose unit is unknown is not a reading.
//
// So the answer has to land in the range a video card can really draw.  Half
// a watt is the floor: no card that is powered at all draws less, and every
// value below it is a driver counting in something other than microwatts.
const (
	minWatts = 0.5
	maxWatts = 2000
)

func wattsFrom(microwatts float64) (float64, bool) {
	w := microwatts / 1e6
	return w, w >= minWatts && w <= maxWatts
}

// hwmon reads the sensor chip a graphics driver registers next to the card.
// It is the same interface the temperature section already reads, one
// directory deeper.
func (r Reader) hwmon(dir string, c *Card, oddPower *bool) {
	chips, err := filepath.Glob(filepath.Join(dir, "hwmon", "hwmon*"))
	if err != nil {
		return
	}
	sort.Strings(chips)
	for _, chip := range chips {
		if c.Celsius == nil {
			if v, ok := number(trimmed(chip, "temp1_input")); ok && v > -50_000 && v < 200_000 {
				celsius := v / 1000
				c.Celsius = &celsius
			}
		}
		if c.Watts == nil {
			for _, file := range []string{"power1_average", "power1_input"} {
				v, ok := number(trimmed(chip, file))
				if !ok {
					continue
				}
				if watts, ok := wattsFrom(v); ok {
					c.Watts = &watts
				} else {
					*oddPower = true
				}
				break
			}
		}
		if c.MHz == nil {
			if v, ok := number(trimmed(chip, "freq1_input")); ok && v > 0 && v < 1e13 {
				mhz := v / 1e6 // the file counts hertz
				c.MHz = &mhz
			}
		}
	}
}

// nvidiaInfo is what /proc/driver/nvidia/gpus/<bus>/information says.
type nvidiaInfo struct {
	Model    string
	VBIOS    string
	Firmware string
}

// nvidiaCards reads the directory the NVIDIA driver publishes, keyed by bus.
func (r Reader) nvidiaCards() map[string]nvidiaInfo {
	root := filepath.Join(r.Proc, "driver/nvidia/gpus")
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	out := map[string]nvidiaInfo{}
	for _, e := range entries {
		text := read(filepath.Join(root, e.Name(), "information"))
		if strings.TrimSpace(text) == "" {
			continue
		}
		info := ParseNvidiaInformation(text)
		bus := e.Name()
		if info.Bus != "" {
			bus = info.Bus
		}
		out[bus] = nvidiaInfo{Model: info.Model, VBIOS: info.VBIOS, Firmware: info.Firmware}
	}
	return out
}

func (c *Card) mergeNvidia(info nvidiaInfo) {
	if info.Model != "" {
		c.Name = info.Model
	}
}

// name gives every card a name.  The driver's own word comes first, then the
// database the distribution ships, and only then the numbers themselves — a
// card is never left nameless, and a name is never invented.
func (r Reader) name(cards []Card) {
	var want [][2]string
	for i := range cards {
		if cards[i].Name == "" && cards[i].VendorID != "" {
			want = append(want, [2]string{cards[i].VendorID, cards[i].DeviceID})
		}
	}
	names := r.lookupPCIIDs(want)
	for i := range cards {
		c := &cards[i]
		key := c.VendorID + ":" + c.DeviceID
		if n, ok := names[key]; ok {
			if c.Vendor == "" {
				c.Vendor = n.vendor
			}
			if c.Name == "" {
				c.Name = n.device
			}
		}
		if c.Vendor == "" {
			c.Vendor = knownVendors[c.VendorID]
		}
		if c.Name == "" {
			c.Name = strings.TrimSpace(c.Vendor + " " + fallbackName(*c))
		}
	}
}

// fallbackName is what a card is called when nothing names it: its own
// numbers.  Two hexadecimal words are not a name a person likes, but they are
// the truth, and the alternative is an empty line where a card is.
func fallbackName(c Card) string {
	switch {
	case c.VendorID != "" && c.DeviceID != "":
		return c.VendorID + ":" + c.DeviceID
	case c.Node != "":
		return c.Node
	default:
		return "видеокарта"
	}
}

// knownVendors names the four vendors whose cards this tool is likely to meet,
// for the machine that ships no pci.ids at all.  It is a courtesy, not a
// database: anything else falls through to its number.
var knownVendors = map[string]string{
	"10de": "NVIDIA",
	"1002": "AMD",
	"1022": "AMD",
	"8086": "Intel",
	"102b": "Matrox",
	"1a03": "ASPEED",
	"15ad": "VMware",
	"1af4": "Red Hat",
}

// read returns the whole of a small pseudo-file, or "".
func read(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

func trimmed(parts ...string) string {
	return strings.TrimSpace(read(filepath.Join(parts...)))
}

func number(s string) (float64, bool) {
	if s == "" {
		return 0, false
	}
	v, err := strconv.ParseFloat(s, 64)
	return v, err == nil
}

func unsigned(s string) (uint64, bool) {
	if s == "" {
		return 0, false
	}
	v, err := strconv.ParseUint(s, 10, 64)
	return v, err == nil
}
