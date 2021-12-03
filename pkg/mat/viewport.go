package mat

type Viewport struct {
	Left   float64
	Top    float64
	Width  float64
	Height float64
	Near   float64
	Far    float64

	w2   float64
	h2   float64
	fmn2 float64
	fpn2 float64
}

func NewViewport(l, t, w, h, n, f float64) Viewport {
	return Viewport{
		Left:   l,
		Top:    t,
		Width:  w,
		Height: h,
		Near:   n,
		Far:    f,
		w2:     w / 2,
		h2:     h / 2,
		fmn2:   (f - n) / 2,
		fpn2:   (f + n) / 2,
	}
}
