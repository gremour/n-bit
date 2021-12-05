package display

type DrawSpriteOpts struct {
	Name string
	SX   float64
	SY   float64
	SW   float64
	SH   float64
	DX   float64
	DY   float64
	DH   float64
	DW   float64
}

func (d *Display) DrawSprite(name string, x, y float64) {
	d.DrawSpriteAdvanced(DrawSpriteOpts{
		Name: name,
		DX:   x,
		DY:   y,
	})
}

func (d *Display) DrawSpriteAdvanced(o DrawSpriteOpts) {
	s, ok := d.Sprites[o.Name]
	if !ok {
		d.ReportSprite(o.Name)
		return
	}
	if o.SW == 0 {
		o.SW = float64(s.Width)
	}
	if o.SH == 0 {
		o.SH = float64(s.Height)
	}
	if o.DW == 0 {
		o.DW = o.SW
	}
	if o.DH == 0 {
		o.DH = o.SH
	}
	d.drawSpriteAdvanced(s, o)
}

func (d *Display) drawSpriteAdvanced(s *Sprite, o DrawSpriteOpts) {
	ri := RectangleInfo{
		RectangleRasterInput: RectangleRasterInput{
			Buffer:       d.Screen.Pixels,
			BufferWidth:  d.Screen.Width,
			BufferHeight: d.Screen.Height,
			Shader: RectShaderIndexed(
				*s.Atlas,
				s.X+int(o.SX),
				s.Y+int(o.SY),
				s.X+int(o.SX)+int(o.SW),
				s.Y+int(o.SY)+int(o.SH),
			),
			Indexizer: d.Indexizer,
			Lights:    d.Lights,
		},
		X: o.DX - float64(s.XOrigin),
		Y: o.DY - float64(s.YOrigin),
		W: o.DW,
		H: o.DH,
	}
	d.Rasterizer.DrawRectangle(ri)
}

// deprecated: 3 times slower vs drawSpriteAdvanced.
func (d *Display) drawSpriteDirect(s *Sprite, x, y float64) {
	// x, y in Screen coords of the top-left sprite coords
	ox, oy := int(x)-s.XOrigin, int(y)-s.YOrigin

	// index of the current line in screen Pixels
	tyind := ox + oy*d.Screen.Width
	// index of the current line in sprite pixels
	syind := s.X + s.Y*s.Atlas.Width
	// current y screen coord
	yy := int(y)
	for sy := 0; sy < s.Height; sy++ {
		// index of the current line in screen Pixels (copy)
		tind := tyind
		// index of the current line in sprite pixels (copy)
		sind := syind
		// current x screen coord
		xx := int(x)
		if tyind >= 0 && tyind < len(d.Screen.Pixels) &&
			syind >= 0 && syind < len(s.Atlas.Pixels) {
			for sx := 0; sx < s.Width; sx++ {
				cursx := s.X + sx
				curtx := ox + sx
				if curtx >= 0 && curtx < d.Screen.Width &&
					cursx >= 0 && cursx < s.Atlas.Width {

					// Sprite pixel color
					scol := s.Atlas.Pixels[sind]
					if scol == 0 {
						continue
					}
					// Light & indexize
					in := d.Lights.Light(scol, xx, yy)
					col := d.Indexizer.Indexize(in, xx, yy)
					d.Screen.Pixels[tind] = col
				}
				tind++
				sind++
				xx++
			}
		}
		tyind += d.Screen.Width
		syind += s.Atlas.Width
		yy++
	}
}
