package display

// Triplets for each color. Must be multiplier of 3 in size.
type Palette []byte

// ColorsNumber returns number of colors in palette.
func (p Palette) ColorsNumber() int {
	return len(p) / 3
}

// Colorize sets pixels to RGBA values from palette.
// Index starts from 1, 0 sets no-color (zero alpha).
func (p Palette) Colorize(canvas []byte, index byte) {
	arr := [4]byte{}
	if index > 0 && int(index) <= len(p)/3 {
		i := (index - 1) * 3
		copy(arr[:], p[i:i+3])
		arr[3] = 255
	}
	copy(canvas, arr[:])
}
