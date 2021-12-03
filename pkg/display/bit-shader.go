package display

var Bits2 = &Ind2Bit{
	Threshold1_5: 0.15,
	Threshold2:   0.3,
	Threshold2_5: 0.45,
	Threshold3:   0.55,
	Threshold3_5: 0.7,
	Threshold4:   0.8,
}

type Ind2Bit struct {
	Threshold1_5 float64
	Threshold2   float64
	Threshold2_5 float64
	Threshold3   float64
	Threshold3_5 float64
	Threshold4   float64
}

type Indexizer interface {
	Indexize(intens float64, posx, posy int) byte
}

// Color2Bit calculates 2-bit color index for input intensity (clamped to 0-1).
// To introduce lighting, add offset and/or multiply intensity by scale factor.
func (i *Ind2Bit) Indexize(intens float64, posx, posy int) byte {
	var color byte
	var dither bool
	switch {
	case intens >= i.Threshold4:
		color = 4
	case intens >= i.Threshold3_5:
		color = 4
		dither = true
	case intens >= i.Threshold3:
		color = 3
	case intens >= i.Threshold2_5:
		color = 3
		dither = true
	case intens >= i.Threshold2:
		color = 2
	case intens >= i.Threshold1_5:
		color = 2
		dither = true
	default:
		color = 1
	}
	if dither && (posx%2 == posy%2) {
		color--
	}
	return color
}

func TriShaderIndexed(iim IndexedImage, tx0, ty0, tx1, ty1, tx2, ty2 int,
	lightOffset, lightScale float64) func(o *TriangleShaderOpts) {
	tx0m := float64(tx0)
	tx1m := float64(tx1)
	tx2m := float64(tx2)
	ty0m := float64(ty0)
	ty1m := float64(ty1)
	ty2m := float64(ty2)
	return func(o *TriangleShaderOpts) {
		cx := int(tx0m*o.W0 + tx1m*o.W1 + tx2m*o.W2)
		cy := int(ty0m*o.W0 + ty1m*o.W1 + ty2m*o.W2)
		c := iim.Pixels[cx+cy*iim.Width]
		if c == 0 {
			return
		}
		in := o.Lights.Light(c, int(o.X), int(o.Y))
		col := o.Indexizer.Indexize(in, int(o.X), int(o.Y))
		o.Buffer[o.BufferOffset] = col
	}
}

func PerspTriShaderIndexed(iim IndexedImage, tx0, ty0, tx1, ty1, tx2, ty2 int,
	z0, z1, z2 float64) func(o *TriangleShaderOpts) {
	tx0m := float64(tx0) * z0
	tx1m := float64(tx1) * z1
	tx2m := float64(tx2) * z2
	ty0m := float64(ty0) * z0
	ty1m := float64(ty1) * z1
	ty2m := float64(ty2) * z2

	return func(o *TriangleShaderOpts) {
		rz := 1 / (z0*o.W0 + z1*o.W1 + z2*o.W2)
		cx := int((tx0m*o.W0 + tx1m*o.W1 + tx2m*o.W2) * rz)
		cy := int((ty0m*o.W0 + ty1m*o.W1 + ty2m*o.W2) * rz)
		c := iim.Pixels[cx+cy*iim.Width]
		if c == 0 {
			return
		}
		in := o.Lights.Light(c, int(o.X), int(o.Y))
		col := o.Indexizer.Indexize(in, int(o.X), int(o.Y))
		o.Buffer[o.BufferOffset] = col
	}
}

func RectShaderIndexed(iim IndexedImage, tx0, ty0, tx1, ty1 int) func(o *RectangleShaderOpts) {
	tw := float64(tx1-tx0) + 1
	th := float64(ty1-ty0) + 1

	return func(o *RectangleShaderOpts) {
		cx := int(float64(tw*o.Px)) + tx0
		cy := int(float64(th*o.Py)) + ty0
		c := iim.Pixels[cx+cy*iim.Width]
		if c == 0 {
			return
		}
		in := o.Lights.Light(c, int(o.X), int(o.Y))
		col := o.Indexizer.Indexize(in, int(o.X), int(o.Y))
		o.Buffer[o.BufferOffset] = col
	}
}
