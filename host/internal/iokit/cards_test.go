// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package iokit

import "testing"

// Разбор записей реестра проверяется на любой машине, а не только на маке:
// вызов и разбор того, что вернулось, разделены здесь так же, как в
// internal/libsystem и internal/darwinsys.
//
// The two trees below are not invented.  They are what the checks read off the
// two Macs of .github/workflows/check.yml on 3 сентября 2026, cut down to the
// keys this package looks at — see the run of the step «IOKit без cgo».  A
// third tree, the discrete card, is built from the key names Apple documents,
// and it is marked as such: neither runner has a card with memory of its own,
// so that branch of the decoder would otherwise go unread.

func str(s string) Value   { return Value{Kind: KindString, Str: s} }
func num(n int64) Value    { return Value{Kind: KindNumber, Num: n} }
func data(b ...byte) Value { return Value{Kind: KindData, Data: b} }
func ascii(s string) Value { return Value{Kind: KindData, Data: append([]byte(s), 0)} }

// intelTree — macos-15-intel, Macmini6,2: акселератор висит на PCI-узле
// display, и всё, что человек читает, лежит на этом узле, а не на самом
// акселераторе.
func intelTree() Entry {
	return Entry{
		Name:  "AppleParavirtAccelerator",
		Class: "AppleParavirtAccelerator",
		Props: map[string]Value{
			"AccelCaps":            num(7),
			"IOAccelRevision":      num(2),
			"IOMatchCategory":      str("IOAccelerator"),
			"MetalPluginClassName": str("AppleParavirtDevice"),
		},
		Parent: &Entry{
			Name:  "display",
			Class: "IOPCIDevice",
			Props: map[string]Value{
				"IOName":     str("display"),
				"class-code": data(0x00, 0x00, 0x03, 0x00),
				"vendor-id":  data(0x6b, 0x10, 0x00, 0x00),
				"device-id":  data(0xee, 0xee, 0x00, 0x00),
				"model":      ascii("Apple Paravirtualized Graphics Device"),
				"pcidebug":   str("0:29:0"),
			},
			Parent: &Entry{
				Name:  "AppleACPIPCI",
				Class: "AppleACPIPCI",
				Props: map[string]Value{
					"IOClass":            str("AppleACPIPCI"),
					"CFBundleIdentifier": str("com.apple.driver.AppleACPIPlatform"),
				},
			},
		},
	}
}

// armTree — macos-latest, VirtualMac2,1: шины у видеокарты нет вовсе, имени
// модели тоже, и это не поломка разбора, а устройство машины.
func armTree() Entry {
	return Entry{
		Name:  "AppleParavirtGPU",
		Class: "AppleParavirtGPU",
		Props: map[string]Value{
			"IOClass":            str("AppleParavirtGPU"),
			"CFBundleIdentifier": str("com.apple.driver.AppleParavirtGPUIOGPUFamily"),
			"DisplayPortCount":   num(1),
			"IOMatchCategory":    str("IOAccelerator"),
		},
		Parent: &Entry{
			Name:  "gfx",
			Class: "AppleARMIODevice",
			Props: map[string]Value{
				"compatible": ascii("paravirtualizedgraphics,gpu"),
				"name":       ascii("gfx"),
			},
		},
	}
}

func TestКартаСIntelМака(t *testing.T) {
	e := intelTree()
	c := card(&e)
	if c.Name != "Apple Paravirtualized Graphics Device" {
		t.Errorf("имя карты %q", c.Name)
	}
	if c.Driver != "AppleParavirtAccelerator" {
		t.Errorf("драйвер %q — взят у родителя вместо самой записи", c.Driver)
	}
	if c.VendorID != "106b" || c.DeviceID != "eeee" || c.Vendor != "Apple" {
		t.Errorf("поставщик %q %q/%q", c.Vendor, c.VendorID, c.DeviceID)
	}
	if c.Bus != "0000:00:1d.0" {
		t.Errorf("адрес на шине %q, а pcidebug говорит 0:29:0", c.Bus)
	}
	if c.MemoryTotalBytes != nil {
		t.Errorf("память %d, а карта её не публикует", *c.MemoryTotalBytes)
	}
	if c.Source != "IORegistry" {
		t.Errorf("источник %q", c.Source)
	}
}

func TestКартаСAppleSilicon(t *testing.T) {
	e := armTree()
	c := card(&e)
	if c.Name != "AppleParavirtGPU" {
		t.Errorf("имя карты %q, а модели в реестре нет — ждали имя драйвера", c.Name)
	}
	if c.Driver != "AppleParavirtGPU" {
		t.Errorf("драйвер %q", c.Driver)
	}
	if c.Bus != "" || c.VendorID != "" || c.Vendor != "" {
		t.Errorf("шина %q, поставщик %q %q — у карты без PCI их быть не может", c.Bus, c.Vendor, c.VendorID)
	}
	if c.MemoryTotalBytes != nil {
		t.Errorf("память %d, а карта делит память с машиной", *c.MemoryTotalBytes)
	}
}

// TestПамятьДискретнойКарты — единственная ветка разбора, которой на бегунках
// нет живого примера: обе маковские машины проверок виртуальные, и своей
// памяти у их видеокарт нет. Ключи взяты из документации Apple, и это сказано
// прямо, а не спрятано.
func TestПамятьДискретнойКарты(t *testing.T) {
	e := Entry{
		Name: "AMDRadeonX6000", Class: "AMDRadeonX6000",
		Props: map[string]Value{"IOClass": str("AMDRadeonX6000")},
		Parent: &Entry{
			Name: "GFX0", Class: "IOPCIDevice",
			Props: map[string]Value{
				"class-code":     data(0x00, 0x00, 0x03, 0x00),
				"vendor-id":      data(0x02, 0x10, 0x00, 0x00),
				"device-id":      data(0x1e, 0x73, 0x00, 0x00),
				"model":          ascii("AMD Radeon Pro 5500M"),
				"pcidebug":       str("1:0:0"),
				"VRAM,totalsize": data(0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00),
			},
		},
	}
	c := card(&e)
	if c.Vendor != "AMD" || c.VendorID != "1002" || c.DeviceID != "731e" {
		t.Errorf("поставщик %q %q/%q", c.Vendor, c.VendorID, c.DeviceID)
	}
	if c.Bus != "0000:01:00.0" {
		t.Errorf("адрес на шине %q", c.Bus)
	}
	if c.MemoryTotalBytes == nil {
		t.Fatalf("память не снялась")
	}
	if got := *c.MemoryTotalBytes; got != 8<<30 {
		t.Errorf("память %d, ожидалось %d", got, uint64(8)<<30)
	}
	// В мегабайтах — вторая запись, которой пользуются старые драйверы.
	e.Parent.Props = map[string]Value{"VRAM,totalMB": num(4096)}
	if n, ok := vram(&e); !ok || n != 4<<30 {
		t.Errorf("VRAM,totalMB дало %d (%v)", n, ok)
	}
}

// TestНеВидеокартаОтсеивается — отрицательный контроль: на маковском бегунке
// записей класса IOPCIDevice двенадцать, и видеокарта среди них одна. Разбор,
// который пропустит сетевую карту, покажет человеку железо, которого он не
// спрашивал.
func TestНеВидеокартаОтсеивается(t *testing.T) {
	ether := Entry{
		Name: "ethernet", Class: "IOPCIDevice",
		Props: map[string]Value{
			"class-code": data(0x00, 0x00, 0x02, 0x00),
			"vendor-id":  data(0xf4, 0x1a, 0x00, 0x00),
		},
	}
	bridge := Entry{
		Name: "pci106b,1a05", Class: "IOPCIDevice",
		Props: map[string]Value{"class-code": data(0x00, 0x00, 0x06, 0x00)},
	}
	display := *intelTree().Parent
	for _, c := range []struct {
		name string
		e    Entry
		want bool
	}{
		{"ethernet", ether, false},
		{"мост", bridge, false},
		{"display", display, true},
	} {
		if got := isDisplay(&c.e); got != c.want {
			t.Errorf("%s: isDisplay = %v, ожидалось %v", c.name, got, c.want)
		}
	}

	seen := map[string]bool{}
	cards := cardsFrom("IOPCIDevice", []Entry{ether, bridge, display}, seen, nil)
	if len(cards) != 1 {
		t.Fatalf("карт из трёх записей PCI собралось %d, ожидалась одна", len(cards))
	}
	// Тот же мак отдаёт эту карту дважды — как акселератор и как устройство
	// PCI. Второй раз она не должна попасть в список.
	cards = cardsFrom("IOAccelerator", []Entry{intelTree()}, seen, cards)
	if len(cards) != 1 {
		t.Fatalf("после второго класса карт %d — двойник не отсеялся", len(cards))
	}
}

func TestЗначенияЧитаютсяКакТекст(t *testing.T) {
	for _, c := range []struct {
		v    Value
		want string
	}{
		{str("AppleParavirtGPU"), "AppleParavirtGPU"},
		{num(-1073741822), "-1073741822"},
		{Value{Kind: KindBool, Bool: true}, "true"},
		{ascii("Apple Paravirtualized Graphics Device"), "Apple Paravirtualized Graphics Device"},
		{data(0x6b, 0x10, 0x00, 0x00), "6b100000"},
		{Value{}, "?"},
	} {
		if got := c.v.Text(); got != c.want {
			t.Errorf("Text() = %q, ожидалось %q", got, c.want)
		}
	}
}
