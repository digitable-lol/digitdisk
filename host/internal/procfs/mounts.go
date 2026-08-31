// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package procfs

import "strings"

// Mount is one line of /proc/mounts (or /proc/self/mounts).
type Mount struct {
	Source  string `json:"source"`
	Point   string `json:"mount_point"`
	FSType  string `json:"fs_type"`
	Options string `json:"options"`
}

// unescapeMount decodes the octal escapes the kernel writes for the
// characters that would otherwise break the space-separated table:
// space, tab, newline and backslash.
func unescapeMount(s string) string {
	if !strings.ContainsRune(s, '\\') {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+3 < len(s) {
			d1, d2, d3 := s[i+1], s[i+2], s[i+3]
			if d1 >= '0' && d1 <= '3' && d2 >= '0' && d2 <= '7' && d3 >= '0' && d3 <= '7' {
				b.WriteByte((d1-'0')<<6 | (d2-'0')<<3 | (d3 - '0'))
				i += 3
				continue
			}
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

// ParseMounts parses the /proc/mounts table.
func ParseMounts(text string) []Mount {
	var out []Mount
	for _, line := range strings.Split(text, "\n") {
		f := strings.Fields(line)
		if len(f) < 3 {
			continue
		}
		m := Mount{
			Source: unescapeMount(f[0]),
			Point:  unescapeMount(f[1]),
			FSType: f[2],
		}
		if len(f) > 3 {
			m.Options = unescapeMount(f[3])
		}
		out = append(out, m)
	}
	return out
}
