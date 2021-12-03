package display

import "math/rand"

var Palettes2Bit = map[string]Palette{
	"cga1": {
		0x0, 0x0, 0x0,
		0xff, 0x55, 0xff,
		0x55, 0xff, 0xff,
		0xff, 0xff, 0xff,
	},
	"cga2": {
		0x0, 0x0, 0x0,
		0xff, 0x55, 0x55,
		0x55, 0xff, 0x55,
		0xff, 0xff, 0x55,
	},
	"red-cyan": {
		0x0, 0x0, 0x0,
		0xff, 0x00, 0x00,
		0x00, 0xff, 0xff,
		0xff, 0xff, 0xff,
	},
	"magenta-green": {
		0x0, 0x0, 0x0,
		0xff, 0x00, 0xff,
		0x00, 0xff, 0x00,
		0xff, 0xff, 0xff,
	},
	"navy-yellow": {
		0x0, 0x0, 0x0,
		0x00, 0x00, 0xff,
		0xff, 0xff, 0x00,
		0xff, 0xff, 0xff,
	},
	"grey-red": {
		0x0, 0x0, 0x0,
		0x77, 0x77, 0x77,
		0xff, 0x00, 0x00,
		0xff, 0xff, 0xff,
	},
	"magenta-pink": {
		0x0, 0x0, 0x0,
		0x33, 0x11, 0x55,
		0xff, 0x55, 0x55,
		0xff, 0xff, 0xff,
	},
	"blue-orange": {
		0x0, 0x0, 0x0,
		0x44, 0x44, 0xff,
		0xff, 0x77, 0x00,
		0xff, 0xff, 0xff,
	},
	"desert": {
		0x0, 0x0, 0x0,
		0xaa, 0x44, 0x22,
		0xff, 0x77, 0x44,
		0xff, 0xff, 0xff,
	},
	"frozen": {
		0x0, 0x0, 0x0,
		0x22, 0x00, 0x77,
		0x00, 0x77, 0xee,
		0xff, 0xff, 0xff,
	},

	"volcanic": {
		0x0, 0x0, 0x0,
		0x33, 0x33, 0x33,
		0xff, 0x33, 0x00,
		0xff, 0xff, 0xff,
	},
	"acid1": {
		0x0, 0x0, 0x0,
		0x55, 0x00, 0x55,
		0x00, 0x77, 0x00,
		0x00, 0xff, 0xff,
	},
	"malachite": {
		0x0, 0x0, 0x0,
		0x00, 0x44, 0x22,
		0x00, 0xff, 0xaa,
		0xff, 0xff, 0xff,
	},
	"green-blood": {
		0x0, 0x0, 0x0,
		0x66, 0x00, 0x00,
		0x00, 0xff, 0x00,
		0xff, 0xff, 0xff,
	},
	"pine": {
		0x0, 0x0, 0x0,
		0x44, 0x22, 0x00,
		0x22, 0x77, 0x22,
		0xff, 0xff, 0xff,
	},
	"darkmoon": {
		0x0, 0x0, 0x0,
		0x00, 0x33, 0x33,
		0x00, 0x77, 0x00,
		0xff, 0xff, 0xff,
	},
	"painted-leather": {
		0x0, 0x0, 0x0,
		0x22, 0x55, 0x77,
		0xbb, 0x44, 0x33,
		0xff, 0xff, 0xff,
	},
}

func RandomPalette2BitName() string {
	n := rand.Intn(len(Palettes2Bit))
	i := 0
	for name := range Palettes2Bit {
		if i == n {
			return name
		}
		i++
	}
	return ""
}

var Palette2BitNames []string

func init() {
	for n := range Palettes2Bit {
		Palette2BitNames = append(Palette2BitNames, n)
	}
}
