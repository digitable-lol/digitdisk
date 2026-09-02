// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package gpuinfo

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// write puts one pseudo-file into the captured tree.
func write(t *testing.T, path, text string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

// amdTree writes what an amdgpu card publishes, copied field for field from a
// live machine (Radeon Graphics, драйвер amdgpu, 2026-09-02).
func amdTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	dev := filepath.Join(root, "sys/class/drm/card1/device")
	write(t, filepath.Join(dev, "uevent"), "DRIVER=amdgpu\nPCI_CLASS=30000\nPCI_ID=1002:13C0\nPCI_SLOT_NAME=0000:77:00.0\n")
	write(t, filepath.Join(dev, "gpu_busy_percent"), "37\n")
	write(t, filepath.Join(dev, "mem_info_vram_total"), "2147483648\n")
	write(t, filepath.Join(dev, "mem_info_vram_used"), "16449536\n")
	hw := filepath.Join(dev, "hwmon/hwmon4")
	write(t, filepath.Join(hw, "name"), "amdgpu\n")
	write(t, filepath.Join(hw, "temp1_input"), "43000\n")
	write(t, filepath.Join(hw, "power1_input"), "18000000\n")
	write(t, filepath.Join(hw, "freq1_input"), "600000000\n")
	// A connector is not a card, and it sits in the same directory.
	write(t, filepath.Join(root, "sys/class/drm/card1-DP-1/uevent"), "DEVTYPE=drm_connector\n")
	write(t, filepath.Join(root, "sys/class/drm/version"), "drm 1.1.0 20060810\n")
	return root
}

func TestCapturedAmdCardIsReadWhole(t *testing.T) {
	root := amdTree(t)
	res := Reader{Sys: filepath.Join(root, "sys"), Proc: filepath.Join(root, "proc")}.Read()
	if len(res.Cards) != 1 {
		t.Fatalf("карт %d, ждали одну: %+v", len(res.Cards), res.Cards)
	}
	c := res.Cards[0]
	if c.Node != "card1" || c.Bus != "0000:77:00.0" || c.Driver != "amdgpu" {
		t.Errorf("карта опознана как %+v", c)
	}
	if c.VendorID != "1002" || c.DeviceID != "13c0" {
		t.Errorf("номера карты %s:%s", c.VendorID, c.DeviceID)
	}
	if c.BusyPercent == nil || *c.BusyPercent != 37 {
		t.Errorf("загрузка %v, ждали 37", c.BusyPercent)
	}
	if c.MemoryTotalBytes == nil || *c.MemoryTotalBytes != 2147483648 {
		t.Errorf("память всего %v", c.MemoryTotalBytes)
	}
	if c.MemoryUsedBytes == nil || *c.MemoryUsedBytes != 16449536 {
		t.Errorf("память занято %v", c.MemoryUsedBytes)
	}
	if c.Celsius == nil || *c.Celsius != 43 {
		t.Errorf("температура %v, ждали 43 °C", c.Celsius)
	}
	if c.Watts == nil || *c.Watts != 18 {
		t.Errorf("мощность %v, ждали 18 Вт", c.Watts)
	}
	if c.MHz == nil || *c.MHz != 600 {
		t.Errorf("частота %v, ждали 600 МГц", c.MHz)
	}
	if res.NoNumbers != "" {
		t.Errorf("карта с числами объявлена немой: %q", res.NoNumbers)
	}
}

// A card whose driver publishes nothing must come back with empty fields and
// with the reason set — never with zeros.
func TestSilentCardHasEmptyFieldsAndAReason(t *testing.T) {
	root := t.TempDir()
	dev := filepath.Join(root, "sys/class/drm/card0/device")
	write(t, filepath.Join(dev, "uevent"), "DRIVER=mgag200\nPCI_ID=102B:0536\nPCI_SLOT_NAME=0000:62:00.0\n")
	res := Reader{Sys: filepath.Join(root, "sys"), Proc: filepath.Join(root, "proc")}.Read()
	if len(res.Cards) != 1 {
		t.Fatalf("карт %d", len(res.Cards))
	}
	c := res.Cards[0]
	if c.BusyPercent != nil || c.MemoryTotalBytes != nil || c.Celsius != nil {
		t.Errorf("немая карта принесла числа: %+v", c)
	}
	if c.Name == "" {
		t.Error("карта осталась без имени")
	}
	if res.NoNumbers == "" {
		t.Error("молчание карты не объяснено")
	}
	if res.NoCards != "" {
		t.Error("найденная карта объявлена ненайденной")
	}
}

func TestNoCardsIsSaidOutLoud(t *testing.T) {
	root := t.TempDir()
	res := Reader{Sys: filepath.Join(root, "sys"), Proc: filepath.Join(root, "proc")}.Read()
	if len(res.Cards) != 0 {
		t.Fatalf("на пустом дереве нашлись карты: %+v", res.Cards)
	}
	if res.NoCards == "" {
		t.Error("пустой список не объяснён")
	}
}

// The NVIDIA driver is the case the whole fallback exists for: it names the
// card and publishes not one number about it.
func TestNvidiaCardIsNamedFromItsOwnDirectory(t *testing.T) {
	root := t.TempDir()
	dev := filepath.Join(root, "sys/class/drm/card0/device")
	write(t, filepath.Join(dev, "uevent"), "DRIVER=nvidia\nPCI_ID=10DE:26B1\nPCI_SLOT_NAME=0000:02:00.0\n")
	write(t, filepath.Join(root, "proc/driver/nvidia/gpus/0000:02:00.0/information"),
		"Model: \t\t NVIDIA RTX 6000 Ada Generation\nIRQ:   \t\t 146\nVideo BIOS: \t 95.02.59.00.09\nBus Location: \t 0000:02:00.0\n")

	r := Reader{Sys: filepath.Join(root, "sys"), Proc: filepath.Join(root, "proc")}
	res := r.Read()
	if len(res.Cards) != 1 {
		t.Fatalf("карт %d: %+v", len(res.Cards), res.Cards)
	}
	if got := res.Cards[0].Name; got != "NVIDIA RTX 6000 Ada Generation" {
		t.Errorf("имя карты %q", got)
	}
	if res.NoNumbers == "" {
		t.Error("немота карты nvidia не объяснена")
	}

	// With the outside program allowed, the numbers arrive — and they arrive
	// marked as having come from outside.
	r.Tool = true
	r.Run = func(name string, args ...string) ([]byte, error) {
		if name != "nvidia-smi" {
			t.Fatalf("запущена посторонняя программа %q", name)
		}
		return []byte("00000000:02:00.0, NVIDIA RTX 6000 Ada Generation, 29035, 49140, 99, 85, 297.12, 2505\n"), nil
	}
	res = r.Read()
	c := res.Cards[0]
	if c.BusyPercent == nil || *c.BusyPercent != 99 {
		t.Errorf("загрузка %v", c.BusyPercent)
	}
	if c.MemoryTotalBytes == nil || *c.MemoryTotalBytes != 49140*1024*1024 {
		t.Errorf("память всего %v", c.MemoryTotalBytes)
	}
	if c.Celsius == nil || *c.Celsius != 85 {
		t.Errorf("температура %v", c.Celsius)
	}
	if !c.Outside || c.Source != "nvidia-smi" {
		t.Errorf("числа чужой программы не помечены: outside=%v source=%q", c.Outside, c.Source)
	}
}

// A row about hardware the files never saw is not a card: it is a claim about
// another machine.
func TestOutsideToolCannotAddHardware(t *testing.T) {
	root := t.TempDir()
	dev := filepath.Join(root, "sys/class/drm/card0/device")
	write(t, filepath.Join(dev, "uevent"), "DRIVER=nvidia\nPCI_ID=10DE:26B1\nPCI_SLOT_NAME=0000:02:00.0\n")
	r := Reader{
		Sys: filepath.Join(root, "sys"), Proc: filepath.Join(root, "proc"), Tool: true,
		Run: func(string, ...string) ([]byte, error) {
			return []byte("00000000:02:00.0, A, 1, 2, 3, 4, 5, 6\n00000000:81:00.0, Выдуманная, 1, 2, 3, 4, 5, 6\n"), nil
		},
	}
	res := r.Read()
	if len(res.Cards) != 1 {
		t.Fatalf("чужая программа добавила железа: %+v", res.Cards)
	}
}

func TestNvidiaSMIRefusesWhatIsNotANumber(t *testing.T) {
	rows := ParseNvidiaSMI("00000000:02:00.0, Карта, [N/A], [N/A], [N/A], 41, [N/A], [N/A]\nмусор\n")
	if len(rows) != 1 {
		t.Fatalf("строк %d", len(rows))
	}
	r := rows[0]
	if r.UsedMiB != nil || r.TotalMiB != nil || r.Busy != nil || r.Watts != nil {
		t.Errorf("«[N/A]» стало числом: %+v", r)
	}
	if r.Celsius == nil || *r.Celsius != 41 {
		t.Errorf("измеренная температура потеряна: %v", r.Celsius)
	}
	if r.Bus != "0000:02:00.0" {
		t.Errorf("адрес на шине %q", r.Bus)
	}
}

func TestNormalizeBusMeetsTheKernelSpelling(t *testing.T) {
	for in, want := range map[string]string{
		"00000000:02:00.0": "0000:02:00.0",
		"0000:02:00.0":     "0000:02:00.0",
		"0:02:00.0":        "0000:02:00.0",
		"вздор":            "",
		"":                 "",
	} {
		if got := NormalizeBus(in); got != want {
			t.Errorf("NormalizeBus(%q) = %q, ждали %q", in, got, want)
		}
	}
}

func TestNvidiaInformationIsReadByItsKeys(t *testing.T) {
	info := ParseNvidiaInformation("Model: \t\t NVIDIA RTX 6000 Ada Generation\n" +
		"GPU UUID: \t GPU-a253\nBus Location: \t 0000:02:00.0\nЧего-то новое: \t 1\n")
	if info.Model != "NVIDIA RTX 6000 Ada Generation" {
		t.Errorf("модель %q", info.Model)
	}
	if info.Bus != "0000:02:00.0" {
		t.Errorf("шина %q", info.Bus)
	}
}

func TestPCIIDsNamesACardAndCachesNothingWrong(t *testing.T) {
	root := t.TempDir()
	ids := filepath.Join(root, "pci.ids")
	write(t, ids, "# комментарий\n"+
		"1002  Advanced Micro Devices, Inc. [AMD/ATI]\n"+
		"\t13c0  Raphael\n"+
		"\t\t1043 8877  Плата\n"+
		"10de  NVIDIA Corporation\n"+
		"\t26b1  AD102GL\n"+
		"C 03  Display controller\n"+
		"\t00  VGA compatible controller\n")

	dev := filepath.Join(root, "sys/class/drm/card1/device")
	write(t, filepath.Join(dev, "uevent"), "DRIVER=amdgpu\nPCI_ID=1002:13C0\nPCI_SLOT_NAME=0000:77:00.0\n")
	r := Reader{Sys: filepath.Join(root, "sys"), Proc: filepath.Join(root, "proc"), IDs: []string{ids}}
	c := r.Read().Cards[0]
	if c.Name != "Raphael" {
		t.Errorf("имя карты %q, ждали Raphael", c.Name)
	}
	if c.Vendor != "Advanced Micro Devices, Inc. [AMD/ATI]" {
		t.Errorf("изготовитель %q", c.Vendor)
	}
}

// A vendor whose device is not in the database still names the vendor: the
// card is then "NVIDIA Corporation 10de:ffff" and not two words of hex alone.
func TestUnknownDeviceKeepsTheVendorName(t *testing.T) {
	root := t.TempDir()
	ids := filepath.Join(root, "pci.ids")
	write(t, ids, "10de  NVIDIA Corporation\n\t26b1  AD102GL\n")
	dev := filepath.Join(root, "sys/class/drm/card0/device")
	write(t, filepath.Join(dev, "uevent"), "DRIVER=nouveau\nPCI_ID=10DE:FFFF\nPCI_SLOT_NAME=0000:03:00.0\n")
	r := Reader{Sys: filepath.Join(root, "sys"), Proc: filepath.Join(root, "proc"), IDs: []string{ids}}
	c := r.Read().Cards[0]
	if c.Vendor != "NVIDIA Corporation" {
		t.Errorf("изготовитель %q", c.Vendor)
	}
	if c.Name != "NVIDIA Corporation 10de:ffff" {
		t.Errorf("имя %q", c.Name)
	}
}

// A card on the bus with no DRM node of its own is still a card.
func TestPCIDeviceWithoutDrmNodeIsFound(t *testing.T) {
	root := t.TempDir()
	dev := filepath.Join(root, "sys/bus/pci/devices/0000:81:00.0")
	write(t, filepath.Join(dev, "class"), "0x030000\n")
	write(t, filepath.Join(dev, "uevent"), "PCI_ID=10DE:1234\nPCI_SLOT_NAME=0000:81:00.0\n")
	// And a device of another class in the same directory is not a card.
	other := filepath.Join(root, "sys/bus/pci/devices/0000:81:00.1")
	write(t, filepath.Join(other, "class"), "0x040300\n")
	write(t, filepath.Join(other, "uevent"), "PCI_ID=10DE:1235\nPCI_SLOT_NAME=0000:81:00.1\n")

	res := Reader{Sys: filepath.Join(root, "sys"), Proc: filepath.Join(root, "proc")}.Read()
	if len(res.Cards) != 1 {
		t.Fatalf("карт %d: %+v", len(res.Cards), res.Cards)
	}
	if res.Cards[0].Bus != "0000:81:00.0" {
		t.Errorf("нашлась не та карта: %+v", res.Cards[0])
	}
}

// A power file whose value is not in microwatts, whatever else it may be in,
// leaves the card without watts and says so once.
func TestPowerInTheWrongUnitIsRefused(t *testing.T) {
	root := t.TempDir()
	dev := filepath.Join(root, "sys/class/drm/card0/device")
	write(t, filepath.Join(dev, "uevent"), "DRIVER=amdgpu\nPCI_ID=1002:13C0\nPCI_SLOT_NAME=0000:77:00.0\n")
	write(t, filepath.Join(dev, "gpu_busy_percent"), "7\n")
	// 18000 microwatts is eighteen milliwatts, which no video card draws;
	// the driver is counting in something else, and we do not guess in what.
	write(t, filepath.Join(dev, "hwmon/hwmon4/power1_input"), "18000\n")
	res := Reader{Sys: filepath.Join(root, "sys"), Proc: filepath.Join(root, "proc")}.Read()
	if got := res.Cards[0].Watts; got != nil {
		t.Errorf("мощность напечатана как %v Вт", *got)
	}
	if res.NoPower == "" {
		t.Error("отказ от мощности не объяснён")
	}

	// A believable reading passes and is in watts.
	write(t, filepath.Join(dev, "hwmon/hwmon4/power1_input"), "18000000\n")
	res = Reader{Sys: filepath.Join(root, "sys"), Proc: filepath.Join(root, "proc")}.Read()
	if got := res.Cards[0].Watts; got == nil || *got != 18 {
		t.Errorf("мощность %v, ждали 18 Вт", got)
	}
	if res.NoPower != "" {
		t.Errorf("правдоподобная мощность объявлена сомнительной: %q", res.NoPower)
	}
}

// One machine can hold a card that answers and a card that does not, and the
// reason has to name the silent one rather than the machine.
func TestSilentCardIsNamedEvenWhenAnotherCardAnswers(t *testing.T) {
	root := t.TempDir()
	amd := filepath.Join(root, "sys/class/drm/card1/device")
	write(t, filepath.Join(amd, "uevent"), "DRIVER=amdgpu\nPCI_ID=1002:13C0\nPCI_SLOT_NAME=0000:77:00.0\n")
	write(t, filepath.Join(amd, "gpu_busy_percent"), "5\n")
	nv := filepath.Join(root, "sys/class/drm/card0/device")
	write(t, filepath.Join(nv, "uevent"), "DRIVER=nvidia\nPCI_ID=10DE:26B1\nPCI_SLOT_NAME=0000:02:00.0\n")

	res := Reader{Sys: filepath.Join(root, "sys"), Proc: filepath.Join(root, "proc")}.Read()
	if len(res.Cards) != 2 {
		t.Fatalf("карт %d", len(res.Cards))
	}
	if res.NoNumbers == "" {
		t.Fatal("молчащая карта не объяснена, потому что вторая карта ответила")
	}
	if !strings.Contains(res.NoNumbers, "nvidia") || !strings.Contains(res.NoNumbers, "--gpu-tool") {
		t.Errorf("объяснение не называет ни драйвер, ни ключ: %q", res.NoNumbers)
	}
}
