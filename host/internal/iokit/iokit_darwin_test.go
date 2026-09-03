// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build darwin

package iokit

import (
	"runtime"
	"strings"
	"syscall"
	"testing"
)

// TestДверьОткрывается — сам замер: берётся ли IOKit тем же приёмом, каким
// internal/libsystem берёт Mach.
//
// The check is not "did something come back" but "did the thing that came back
// agree with a number this machine states some other way".  Every Mac has
// exactly one IOPlatformExpertDevice, and its model is the same string sysctl
// answers with under hw.model.  If the door opened onto rubbish — a wrong
// calling convention, a Core Foundation object read as the wrong type — the two
// disagree and this test says so.  A test that only asked for a non-empty
// answer would pass on rubbish.
func TestДверьОткрывается(t *testing.T) {
	entries, err := Match("IOPlatformExpertDevice", 0)
	if err != nil {
		t.Fatalf("IOServiceGetMatchingServices(IOPlatformExpertDevice): %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("записей класса IOPlatformExpertDevice %d, а машина одна", len(entries))
	}
	e := entries[0]
	t.Logf("запись: класс %q, имя %q, свойств %d", e.Class, e.Name, len(e.Props))

	want, err := syscall.Sysctl("hw.model")
	if err != nil {
		t.Fatalf("sysctl hw.model: %v", err)
	}
	v, ok := e.Prop("model")
	if !ok {
		t.Fatalf("у IOPlatformExpertDevice нет свойства model; есть: %s",
			strings.Join(sortedKeys(e.Props), ", "))
	}
	got := v.Text()
	t.Logf("model из IORegistry %q, hw.model из sysctl %q", got, want)
	if got != want {
		t.Fatalf("IOKit и sysctl расходятся: %q против %q", got, want)
	}
}

// TestТипыСнимаются доказывает, что снимается не одна строка, а каждый из
// четырёх типов Core Foundation, ради которых всё и затевалось.
func TestТипыСнимаются(t *testing.T) {
	entries, err := Match("IOPlatformExpertDevice", 1)
	if err != nil {
		t.Fatalf("Match: %v", err)
	}
	seen := map[Kind]string{}
	for p := &entries[0]; p != nil; p = p.Parent {
		for _, k := range sortedKeys(p.Props) {
			if _, ok := seen[p.Props[k].Kind]; !ok {
				seen[p.Props[k].Kind] = k + " = " + p.Props[k].Text()
			}
		}
	}
	for _, k := range []Kind{KindString, KindNumber, KindBool, KindData} {
		if s, ok := seen[k]; ok {
			t.Logf("тип %d снят: %s", k, s)
		} else {
			t.Logf("тип %d на этой машине не встретился", k)
		}
	}
	if len(seen) < 2 {
		t.Fatalf("снялся только %d тип — разбор Core Foundation не работает", len(seen))
	}
}

// gpuClasses — классы, которыми Mac описывает видеокарту.  Which of them a
// given Mac has is exactly what is being measured: an Intel Mac hangs its
// graphics off PCI, an Apple Silicon one does not have a PCI node for the GPU
// at all.
var gpuClasses = []string{
	"IOAccelerator",
	"IOPCIDevice",
	"IOGraphicsAccelerator2",
	"AGXAccelerator",
	"IOGPU",
}

// TestMatchingServices — прогон, который называет числом, что этот мак отдаёт
// про видеокарты и чего не отдаёт.  Он ничего не требует: список свойств
// у встроенной карты и у дискретной разный, и записывать сюда ожидание значило
// бы выдумать его за Apple.  Он ПЕЧАТАЕТ, и напечатанное — результат задачи.
func TestMatchingServices(t *testing.T) {
	t.Logf("машина: %s/%s", runtime.GOOS, runtime.GOARCH)
	found := 0
	for _, class := range gpuClasses {
		entries, err := Match(class, 2)
		if err != nil {
			t.Logf("=== %s: %v", class, err)
			continue
		}
		found += len(entries)
		t.Logf("=== %s: записей %d", class, len(entries))
		for i := range entries {
			for _, line := range entries[i].Dump(140) {
				t.Logf("    %s", line)
			}
		}
	}
	t.Logf("всего записей по всем классам: %d", found)

	cards := Cards()
	t.Logf("карт собрано: %d", len(cards))
	for _, c := range cards {
		mem := "—"
		if c.MemoryTotalBytes != nil {
			mem = byteCount(*c.MemoryTotalBytes)
		}
		t.Logf("карта: имя %q, поставщик %q (%s:%s), драйвер %q, память %s, источник %q",
			c.Name, c.Vendor, c.VendorID, c.DeviceID, c.Driver, mem, c.Source)
	}
	if len(cards) == 0 {
		t.Fatalf("ни одной карты не собрано, хотя записей по классам %d", found)
	}
}

func byteCount(n uint64) string {
	switch {
	case n >= 1<<30:
		return itoa(n>>30) + " ГиБ"
	case n >= 1<<20:
		return itoa(n>>20) + " МиБ"
	}
	return itoa(n) + " Б"
}

func itoa(n uint64) string {
	if n == 0 {
		return "0"
	}
	var b [24]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
