// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package report

import (
	"bytes"
	"strings"
	"testing"

	"digitdisk/internal/core"
	"digitdisk/internal/scan"
	"digitdisk/internal/sysinfo"
)

func TestBytes(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{0, "0 Б"},
		{512, "512 Б"},
		{1024, "1.0 КиБ"},
		{1536, "1.5 КиБ"},
		{1 << 20, "1.0 МиБ"},
		{1 << 30, "1.0 ГиБ"},
		{1 << 40, "1.0 ТиБ"},
		{-2048, "-2.0 КиБ"},
	}
	for _, c := range cases {
		if got := Bytes(c.in); got != c.want {
			t.Errorf("Bytes(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestStatusPrintsDashesForMissingData(t *testing.T) {
	var buf bytes.Buffer
	Status(&buf, sysinfo.Status{})
	out := buf.String()
	if !strings.Contains(out, "узел          —") {
		t.Errorf("an absent hostname must print as —, got:\n%s", out)
	}
	if !strings.Contains(out, "занято ЦП     — (замер не делался)") {
		t.Errorf("an unsampled CPU must say so rather than print 0%%:\n%s", out)
	}
	if strings.Contains(out, "0.0%") && !strings.Contains(out, "средняя       0.00") {
		t.Errorf("empty status invented a percentage:\n%s", out)
	}
}

func TestAnalyzeWarnsWhenDeciderIsAStub(t *testing.T) {
	var buf bytes.Buffer
	Analyze(&buf, scan.Result{
		Root:      "/x",
		Decider:   core.Default().Name(),
		ByClass:   map[core.Class]scan.Bucket{},
		ByVerdict: map[core.Verdict]scan.Bucket{},
	})
	out := buf.String()
	if !strings.Contains(out, "ВНИМАНИЕ") {
		t.Errorf("a report backed by the stub must say so:\n%s", out)
	}
	for _, c := range core.Classes {
		if !strings.Contains(out, string(c)) {
			t.Errorf("разряд %q missing from the report", c)
		}
	}
}
