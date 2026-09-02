// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build linux

package sysinfo

import (
	"os"
	"path/filepath"
	"testing"
)

// A captured tree with a video card in it: the card reaches the snapshot, and
// so does the reason its driver publishes no numbers.
func TestVideoCardsReachTheSnapshot(t *testing.T) {
	dir := t.TempDir()
	write := func(rel, content string) {
		p := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("proc/stat", "cpu  100 0 100 800 0 0 0 0 0 0\ncpu0 100 0 100 800 0 0 0 0 0 0\nbtime 1700000000\n")
	write("proc/cpuinfo", "processor\t: 0\nmodel name\t: Придуманный процессор 9000\n")
	write("sys/class/dmi/id/sys_vendor", "Своя мастерская\n")
	write("sys/class/dmi/id/product_name", "Машина 1\n")
	write("sys/class/drm/card0/device/uevent",
		"DRIVER=amdgpu\nPCI_ID=1002:13C0\nPCI_SLOT_NAME=0000:77:00.0\n")
	write("sys/class/drm/card0/device/gpu_busy_percent", "42\n")
	write("sys/class/drm/card0/device/mem_info_vram_total", "2147483648\n")
	write("sys/class/drm/card0/device/mem_info_vram_used", "1073741824\n")

	c := Collector{Proc: filepath.Join(dir, "proc"), Sys: filepath.Join(dir, "sys"), Etc: filepath.Join(dir, "etc"), Top: 3}
	st := c.Collect()

	if len(st.GPUs) != 1 {
		t.Fatalf("видеокарт %d, ждали одну: %+v", len(st.GPUs), st.GPUs)
	}
	g := st.GPUs[0]
	if g.BusyPercent == nil || *g.BusyPercent != 42 {
		t.Errorf("загрузка карты %v", g.BusyPercent)
	}
	if g.MemoryUsedBytes == nil || *g.MemoryUsedBytes != 1<<30 {
		t.Errorf("память карты %v", g.MemoryUsedBytes)
	}
	if g.Driver != "amdgpu" || g.Bus != "0000:77:00.0" {
		t.Errorf("карта опознана как %+v", g)
	}
	if _, unmeasured := st.Unmeasured(FactGPUs); unmeasured {
		t.Error("найденная карта объявлена ненайденной")
	}
	if st.Host.Model != "Своя мастерская Машина 1" {
		t.Errorf("модель машины %q", st.Host.Model)
	}
	if st.Host.CPUModel != "Придуманный процессор 9000" {
		t.Errorf("модель процессора %q", st.Host.CPUModel)
	}

	// Two readings of a capture that does not move: nothing is published
	// about the processors, and the reason says so instead of showing
	// zeros.
	if len(st.Load.Cores) != 0 {
		t.Errorf("по неподвижному снимку насчитаны доли ядер: %+v", st.Load.Cores)
	}
	if why, unmeasured := st.Unmeasured(FactCores); !unmeasured || why == "" {
		t.Error("отсутствие долей по ядрам не объяснено")
	}
}

// A machine with no video card at all says so, and says it once.
func TestNoVideoCardsIsNamed(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "proc"), 0o755); err != nil {
		t.Fatal(err)
	}
	c := Collector{Proc: filepath.Join(dir, "proc"), Sys: filepath.Join(dir, "sys"), Etc: filepath.Join(dir, "etc"), Top: 3}
	st := c.Collect()
	if len(st.GPUs) != 0 {
		t.Fatalf("на пустом дереве нашлись карты: %+v", st.GPUs)
	}
	if why, unmeasured := st.Unmeasured(FactGPUs); !unmeasured || why == "" {
		t.Error("пустой список карт не объяснён")
	}
}
