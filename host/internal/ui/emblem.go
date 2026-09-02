// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package ui

// The mark drawn beside the system's own facts, the way a screenful of them
// has been drawn since the first tool of this kind.
//
// The marks are ours.  Every one of them was drawn here, stroke by stroke, and
// none is copied from another project's collection: those collections are
// somebody's work under somebody's licence, and a picture is code enough to
// carry one.  They are also deliberately plain — a mark is a hint at which
// system this is, not a portrait of a logo, and a hint is what survives being
// cut to forty columns.
//
// Only the printable ASCII of the first half of the table is used, and that is
// a decision, not an accident: the rest of the screen speaks Russian and needs
// a terminal that draws it, but the mark must still be a mark on a terminal
// with a font from 1990 and LANG=C.  Nothing here is wider than
// emblemMaxWidth, so the fields beside it keep their column on a narrow
// terminal.
import "strings"

// emblemMaxWidth is the widest a mark may be drawn.  It is checked by a test,
// because a mark one cell too wide moves every field on the page.
const emblemMaxWidth = 14

// emblem is one mark: the picture and the palette role it is painted in.
type emblem struct {
	art  []string
	role func(Palette) slot
}

// emblems are the marks by the system's own identifier — the ID= line of
// /etc/os-release, which is the name a distribution calls itself by.
var emblems = map[string]emblem{
	"ubuntu": {art: lines(`
     .--.
   ,'    ',
  /   ()   \
 |  ()      |
  \      () /
   ',    ,'
     '--'
`), role: func(p Palette) slot { return p.Orange }},

	"debian": {art: lines(`
     .---.
   ,' .-. ',
  /  /   \  \
 |  |    |  |
  \  \   /  /
   ',  -'  '
     '---'
`), role: func(p Palette) slot { return p.Red }},

	"fedora": {art: lines(`
     .---.
   ,'  _  ',
  /   |_    \
 |    |     |
  \   |    /
   ',     ,'
     '---'
`), role: func(p Palette) slot { return p.Blue }},

	"arch": {art: lines(`
      /\
     /  \
    / /\ \
   / /  \ \
  / /____\ \
 / /  __  \ \
/_/  /  \  \_\
`), role: func(p Palette) slot { return p.Blue }},

	"alpine": {art: lines(`
 +----------+
 |    /\    |
 |   /  \   |
 |  / /\ \  |
 | / /  \ \ |
 |/_/    \_\|
 +----------+
`), role: func(p Palette) slot { return p.Blue }},

	"nixos": {art: lines(`
   \\    //
    \\  //
 ====\\//====
     //\\
    //  \\
   //    \\
`), role: func(p Palette) slot { return p.AccentSoft }},

	"opensuse": {art: lines(`
    .-----.
   / o   o \
  |    ^    |
  |   ---   |
   \       /
    '.___.'
`), role: func(p Palette) slot { return p.Green }},

	"linuxmint": {art: lines(`
 +----------+
 |  |\  /|  |
 |  | \/ |  |
 |  |    |  |
 |  |    |  |
 +----------+
`), role: func(p Palette) slot { return p.Green }},

	"gentoo": {art: lines(`
    .----.
   / .--. \
  | |    | |
   \ '--'  |
    '.    /
      '--'
`), role: func(p Palette) slot { return p.Purple }},

	"rhel": {art: lines(`
    _______
   /       \
  ( o     o )
   |  ___  |
   |_______|
  '---------'
`), role: func(p Palette) slot { return p.Red }},

	"macos": {art: lines(`
        .-'
      .'
    .-----.
   /       \
  |         |
   \       /
    '.   .'
      '-'
`), role: func(p Palette) slot { return p.AccentSoft }},

	"bsd": {art: lines(`
   .-------.
  /  ^   ^  \
 |   (o o)   |
 |    ---    |
  \  \___/  /
   '-------'
`), role: func(p Palette) slot { return p.Orange }},

	// The mark for a system with no mark of its own: a screen and a stand.
	// It is a machine, which is the one thing every case here has in common.
	"": {art: lines(`
 +----------+
 |          |
 |   >_     |
 |          |
 +----------+
     |__|
   _______
`), role: func(p Palette) slot { return p.Accent }},
}

// family maps the distributions that are one another's rebuilds onto the mark
// of the family.  A machine gets the mark of what it is, not an empty space,
// and there is no pretending a rebuild is a different system.
var family = map[string]string{
	"rocky":       "rhel",
	"almalinux":   "rhel",
	"centos":      "rhel",
	"rhel":        "rhel",
	"fedora":      "fedora",
	"ol":          "rhel",
	"pop":         "ubuntu",
	"neon":        "ubuntu",
	"elementary":  "ubuntu",
	"zorin":       "ubuntu",
	"linuxmint":   "linuxmint",
	"mint":        "linuxmint",
	"raspbian":    "debian",
	"devuan":      "debian",
	"kali":        "debian",
	"manjaro":     "arch",
	"endeavouros": "arch",
	"arch":        "arch",
	"artix":       "arch",
	"opensuse":    "opensuse",
	"sles":        "opensuse",
	"nixos":       "nixos",
	"gentoo":      "gentoo",
	"alpine":      "alpine",
	"void":        "gentoo",
	"debian":      "debian",
	"ubuntu":      "ubuntu",
	"macos":       "macos",
	"darwin":      "macos",
	"freebsd":     "bsd",
	"openbsd":     "bsd",
	"netbsd":      "bsd",
	"dragonfly":   "bsd",
}

// emblemFor picks the mark for a system by what it calls itself.  Anything
// unrecognised gets the general mark: an empty space where the picture should
// be would look like a fault, and this is not one.
func emblemFor(id, pretty string) emblem {
	key := strings.ToLower(strings.TrimSpace(id))
	if name, ok := family[key]; ok {
		return emblems[name]
	}
	// A system that does not name itself in a field we know is still often
	// named in the line a person reads.
	low := strings.ToLower(pretty)
	for _, guess := range []string{"ubuntu", "debian", "fedora", "arch", "alpine",
		"nixos", "opensuse", "gentoo", "macos", "bsd"} {
		if strings.Contains(low, guess) {
			if name, ok := family[guess]; ok {
				return emblems[name]
			}
			return emblems[guess]
		}
	}
	return emblems[""]
}

// lines cuts a drawing into rows and pads them all to one width, so the fields
// printed beside the mark start in the same column on every row.
func lines(art string) []string {
	rows := strings.Split(strings.Trim(art, "\n"), "\n")
	width := 0
	for _, r := range rows {
		if n := runes(r); n > width {
			width = n
		}
	}
	for i, r := range rows {
		rows[i] = r + strings.Repeat(" ", width-runes(r))
	}
	return rows
}

// emblemWidth is how wide a drawn mark is.
func emblemWidth(art []string) int {
	if len(art) == 0 {
		return 0
	}
	return runes(art[0])
}
