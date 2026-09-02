// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package ui

import (
	"fmt"
	"os"
	"strings"
)

// Digitable Focus — фирменная палитра стека, канон живёт в
// courses/products/workbench/themes/focus-palettes.json (версия 0.1.0,
// схема https://digitable.life/schemas/focus-palettes-v1.json).  Те же
// значения повторены как CSS-переменные портала (--digitable-*) и лендинга
// (--portal-*).  Здесь они переписаны числами, а не кодом: цвет — это факт о
// бренде, и своей лицензии он не приносит.
//
// Раскладка на ANSI взята оттуда же (workbench/scripts/targets/_ansi.mjs):
// cyan = accent, brightCyan = accentSoft, magenta = purple, white = muted,
// brightWhite = foreground.  Терминал без 24 бит цвета получает ту же
// айдентику огрублённой, а не чужую.

// RGB is one colour of the palette.
type RGB struct{ R, G, B uint8 }

// ansi16 is the slot this colour falls back to on a sixteen-colour terminal,
// following the stack-wide mapping named above.  30..37 dim, 90..97 bright.
type slot struct {
	c    RGB
	ansi int
}

// Palette names the roles the screen draws with.  The roles are the ones the
// canon file names; nothing is invented here.
type Palette struct {
	Name       string
	Dark       bool
	Background slot
	Surface    slot
	Foreground slot
	Muted      slot
	Subtle     slot
	Border     slot
	Accent     slot
	AccentSoft slot
	Blue       slot
	Green      slot
	Yellow     slot
	Orange     slot
	Purple     slot
	Red        slot
}

// Carbon is the default palette of the brand: dark canvas, cyan signal.
var Carbon = Palette{
	Name: "Digitable Focus Carbon", Dark: true,
	Background: slot{RGB{0x05, 0x08, 0x0D}, 30},
	Surface:    slot{RGB{0x07, 0x10, 0x18}, 30},
	Foreground: slot{RGB{0xF5, 0xF7, 0xFA}, 97},
	Muted:      slot{RGB{0x9B, 0xAA, 0xB8}, 37},
	Subtle:     slot{RGB{0x71, 0x86, 0x95}, 90},
	Border:     slot{RGB{0x15, 0x56, 0x6A}, 90},
	Accent:     slot{RGB{0x00, 0xE5, 0xE5}, 36},
	AccentSoft: slot{RGB{0x00, 0xD8, 0xFF}, 96},
	Blue:       slot{RGB{0x3C, 0xA9, 0xFF}, 34},
	Green:      slot{RGB{0x7C, 0xFF, 0x6B}, 32},
	Yellow:     slot{RGB{0xFF, 0xC2, 0x47}, 33},
	Orange:     slot{RGB{0xFF, 0x8A, 0x2A}, 33},
	Purple:     slot{RGB{0xB6, 0x5C, 0xFF}, 35},
	Red:        slot{RGB{0xFF, 0x5B, 0x5B}, 31},
}

// Paper is the light palette of the same canon.
var Paper = Palette{
	Name: "Digitable Focus Paper", Dark: false,
	Background: slot{RGB{0xF4, 0xF7, 0xF8}, 97},
	Surface:    slot{RGB{0xFF, 0xFF, 0xFF}, 97},
	Foreground: slot{RGB{0x10, 0x20, 0x2A}, 30},
	Muted:      slot{RGB{0x52, 0x6A, 0x78}, 90},
	Subtle:     slot{RGB{0x6B, 0x82, 0x8F}, 90},
	Border:     slot{RGB{0xA9, 0xC2, 0xC9}, 37},
	Accent:     slot{RGB{0x00, 0x7C, 0x83}, 36},
	AccentSoft: slot{RGB{0x00, 0x6E, 0x78}, 36},
	Blue:       slot{RGB{0x1C, 0x5E, 0xA8}, 34},
	Green:      slot{RGB{0x28, 0x7A, 0x36}, 32},
	Yellow:     slot{RGB{0x8A, 0x5B, 0x00}, 33},
	Orange:     slot{RGB{0xA5, 0x4A, 0x00}, 33},
	Purple:     slot{RGB{0x71, 0x39, 0xA3}, 35},
	Red:        slot{RGB{0xB4, 0x23, 0x2E}, 31},
}

// Signal is the high-contrast palette of the same canon.
var Signal = Palette{
	Name: "Digitable Focus Signal", Dark: true,
	Background: slot{RGB{0x00, 0x00, 0x00}, 30},
	Surface:    slot{RGB{0x00, 0x00, 0x00}, 30},
	Foreground: slot{RGB{0xFF, 0xFF, 0xFF}, 97},
	Muted:      slot{RGB{0xE6, 0xF7, 0xFF}, 37},
	Subtle:     slot{RGB{0x9F, 0xC4, 0xD4}, 90},
	Border:     slot{RGB{0x66, 0xFF, 0xFF}, 96},
	Accent:     slot{RGB{0x00, 0xFF, 0xFF}, 96},
	AccentSoft: slot{RGB{0x66, 0xFF, 0xFF}, 96},
	Blue:       slot{RGB{0x5C, 0xC8, 0xFF}, 94},
	Green:      slot{RGB{0x7C, 0xFF, 0x00}, 92},
	Yellow:     slot{RGB{0xFF, 0xD8, 0x00}, 93},
	Orange:     slot{RGB{0xFF, 0x9D, 0x3D}, 93},
	Purple:     slot{RGB{0xE4, 0x7A, 0xFF}, 95},
	Red:        slot{RGB{0xFF, 0x4D, 0x5A}, 91},
}

// PaletteByName returns a palette of the canon by its short name.  An unknown
// name yields Carbon, the default of the stack.
func PaletteByName(name string) Palette {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "paper":
		return Paper
	case "signal":
		return Signal
	default:
		return Carbon
	}
}

// depth is how much colour the terminal will accept.
type depth int

const (
	depthNone depth = iota // NO_COLOR, or a terminal that cannot colour
	depth16                // the sixteen ANSI slots
	depth256               // the xterm cube
	depthTrue              // twenty-four bits
)

// Theme is a palette bound to what the terminal can show.
type Theme struct {
	P Palette
	d depth
}

// NewTheme reads the environment and decides how the palette will be spoken.
// NO_COLOR is honoured (no-color.org): the screen still runs, it is simply
// drawn without colour.
func NewTheme(p Palette) Theme {
	return Theme{P: p, d: detectDepth()}
}

func detectDepth() depth {
	if _, off := os.LookupEnv("NO_COLOR"); off {
		return depthNone
	}
	term := os.Getenv("TERM")
	if term == "" || term == "dumb" {
		return depthNone
	}
	switch strings.ToLower(os.Getenv("COLORTERM")) {
	case "truecolor", "24bit":
		return depthTrue
	}
	if strings.Contains(term, "direct") {
		return depthTrue
	}
	if strings.Contains(term, "256") {
		return depth256
	}
	return depth16
}

// Canvas is the sequence that sets the brand background.  Every drawn line
// opens with it, and ESC[K then erases the rest of the line in that colour, so
// the screen is a Digitable canvas rather than whatever the terminal had.
func (t Theme) Canvas() string {
	if t.d == depthNone {
		return ""
	}
	return t.seq(t.P.Background, true)
}

// reset ends a painted run.  Plain ESC[0m would drop the canvas with the
// colour, so the canvas is laid back down straight after it.
func (t Theme) reset() string { return "\x1b[0m" + t.Canvas() }

// Fg paints s in the foreground colour of the given role.
func (t Theme) Fg(s slot, text string) string {
	if t.d == depthNone || text == "" {
		return text
	}
	return t.seq(s, false) + text + t.reset()
}

// Bold paints s in the role and turns the weight up.
func (t Theme) Bold(s slot, text string) string {
	if text == "" {
		return text
	}
	if t.d == depthNone {
		return "\x1b[1m" + text + "\x1b[0m"
	}
	return "\x1b[1m" + t.seq(s, false) + text + t.reset()
}

// Chip paints text as a filled label: the role as the background, the canvas
// colour as the letters.  This is how the header mark and the current section
// are set apart.
func (t Theme) Chip(s slot, text string) string {
	if t.d == depthNone || text == "" {
		return "\x1b[7m" + text + "\x1b[0m"
	}
	return "\x1b[1m" + t.seq(s, true) + t.seq(t.P.Background, false) + text + t.reset()
}

func (t Theme) seq(s slot, background bool) string {
	switch t.d {
	case depthTrue:
		lead := 38
		if background {
			lead = 48
		}
		return fmt.Sprintf("\x1b[%d;2;%d;%d;%dm", lead, s.c.R, s.c.G, s.c.B)
	case depth256:
		lead := 38
		if background {
			lead = 48
		}
		return fmt.Sprintf("\x1b[%d;5;%dm", lead, cube256(s.c))
	case depth16:
		n := s.ansi
		if background {
			n += 10 // 30..37 -> 40..47, 90..97 -> 100..107
		}
		return fmt.Sprintf("\x1b[%dm", n)
	}
	return ""
}

// cube256 places an RGB colour in the xterm-256 palette: the 6×6×6 cube, or
// the grey ramp when the three channels sit close together.
func cube256(c RGB) int {
	r, g, b := int(c.R), int(c.G), int(c.B)
	if abs(r-g) < 12 && abs(g-b) < 12 && abs(r-b) < 12 {
		grey := (r + g + b) / 3
		switch {
		case grey < 8:
			return 16
		case grey > 248:
			return 231
		}
		return 232 + (grey-8)*23/240
	}
	q := func(v int) int { return (v*5 + 127) / 255 }
	return 16 + 36*q(r) + 6*q(g) + q(b)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
