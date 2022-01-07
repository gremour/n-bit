package display

import (
	"fmt"
	"image"
)

// IndexedImage is a structure that represents image in indexed mode.
type IndexedImage struct {
	Width  int
	Height int

	// Pixels of the image, 1 byte per pixel.
	// Color 0 is reserved for no color.
	Pixels []byte

	// Number of colors image supports.
	Colors int
}

// FromImageOpts is a structure with options for IndexedImageFromImage function.
type FromImageOpts struct {
	// Number of colors to produce. Default is 255.
	Colors int

	// Alpha threshold. Alpha that is lower this is considered invisible.
	// Default is 0.5. Less than 0 produces no invisible pixels.
	AlphaThreshold float64
}

type ToRGBAOpts struct {
	// Prepared Pixels buffer to use.
	Pixels []byte

	// Palette to convert indexed color to RGBA.
	Palette Palette
}

// IndexedImageFromImage converts stdlib image.Image to IndexedImage with a given number
// of colors. Image is first converted to grascale, then intensity of each pixel
// is mapped to the color range. 0 is reserved for no color. No color is set for
// pixels which alpha is lower than alpha threshold.
func IndexedImageFromImage(img image.Image, o FromImageOpts) IndexedImage {
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	if o.Colors < 1 || o.Colors > 255 {
		o.Colors = 255
	}
	if o.AlphaThreshold == 0 {
		o.AlphaThreshold = 0.5
	}
	iim := IndexedImage{
		Width:  w,
		Height: h,
		Pixels: make([]byte, w*h),
		Colors: o.Colors,
	}
	if w == 0 || h == 0 {
		return iim
	}

	for i := 0; i < len(iim.Pixels); i++ {
		r, g, b, a := img.At(i%iim.Width, i/iim.Width).RGBA()
		fa := float64(a)
		if fa/0xffff < o.AlphaThreshold {
			continue
		}
		// Compensate for alpha-premultiplication and scale to 0-1.
		intensity := float64(r)/fa*0.21 + float64(g)/fa*0.72 + float64(b)/fa*0.07
		const cap = float64(0xffff) / 0x10000
		if intensity > cap {
			intensity = cap
		} else if intensity < 0 {
			intensity = 0
		}
		iim.Pixels[i] = byte(intensity * (float64(o.Colors + 1)))
	}
	return iim
}

// Fill buffer with provided byte.
func (iim IndexedImage) Fill(v byte) {
	// Preload the first value into the slice
	iim.Pixels[0] = v

	// Incrementally duplicate the value into the rest of the container
	for i := 1; i < len(iim.Pixels); i *= 2 {
		copy(iim.Pixels[i:], iim.Pixels[:i])
	}
}

// ToRGBA converts IndexedImage to RGBA mode.
// Resulting slice is 4 times larger than Pixels size.
// You can provide existing buffer in ToRGBAOpts.Pixels
// to avoid memory allocation, but it need to be exactly
// 4 times larger than Pixels of this image or function
// will panic.
func (iim IndexedImage) ToRGBA(o ToRGBAOpts) []byte {
	if o.Pixels == nil {
		o.Pixels = make([]byte, 4*iim.Width*4*iim.Height)
	} else {
		if len(o.Pixels) != 4*iim.Width*iim.Height {
			panic(fmt.Sprintf("IndexedImage.ToRGBA: Passed bytes size do not match, passed: %v expected: %v",
				len(o.Pixels), 4*iim.Width*iim.Height))
		}
	}
	for i := range iim.Pixels {
		o.Palette.Colorize(o.Pixels[i*4:i*4+3], iim.Pixels[i])
	}
	return o.Pixels
}
