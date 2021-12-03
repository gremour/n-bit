package display

type Display struct {
	Screen     IndexedImage
	RGBA       []byte
	Palette    Palette
	Rasterizer Rasterizer
	Atlases    map[string]*IndexedImage
	Sprites    map[string]*Sprite
	Lights     Lights
	Indexizer  Indexizer

	reportedSprite map[string]struct{}
}

func (d *Display) InitBuffers(w, h int) {
	d.Screen = IndexedImage{
		Width:  w,
		Height: h,
		Pixels: make([]byte, w*h),
	}
	d.RGBA = make([]byte, w*h*4)
}
