package display

import (
	"math"
	"sync"
)

// RectangleRasterInput part of options that is passed to rasterizer workers.
type RectangleRasterInput struct {
	// Pixel buffer. Shader needs to output color to buffer with
	// boffs offset and value calculated by shader. Buffer can use multiple
	// bytes for pixel. Multiply boffs by pixel size in this case.
	Buffer       []byte
	BufferWidth  int
	BufferHeight int

	// Function to call for each fragment.
	// Shader is responsible for saving pixel to the buffer. Use lambda function that
	// has access to the buffer to do that. See bit-shader for example.
	Shader func(o *RectangleShaderOpts)

	// Raster chunks size (in bits, e.g. 3 bits = 8 pixels).
	ChunkBits int

	// Lights model.
	Lights Lights
	// Object to convert to index color.
	Indexizer Indexizer

	// Any other data to pass to shader.
	Extra interface{}
}

// RectangleRasterStatic part of options including input and precalculated for the whole triangle.
type RectangleRasterStatic struct {
	RectangleRasterInput

	// Percentage step per pixel.
	Pxs, Pys float64
}

// ShaderOpts options for shader.
type RectangleShaderOpts struct {
	RectangleRasterStatic

	// Cooridnates in buffer, pixels.
	X, Y float64

	// Percentage coords.
	Px, Py float64

	// Offset in buffer, pixels. Equal to int(x) + int(y) * BufferWidth.
	BufferOffset int
}

// RectangleInfo input data to DrawRectangle call.
type RectangleInfo struct {
	RectangleRasterInput

	X, Y float64
	W, H float64
}

// RectangleChunk data for rectangle chunk to pass to rasterizer worker.
type RectangleChunk struct {
	RectangleRasterStatic

	// Origin x, y, pixels.
	Xo, Yo float64

	// Origin percentage coords.
	Pxo, Pyo float64

	// Offset in buffer, pixels. Equal to x + y * BufferWidth.
	BufferOffset int

	// Chunk width in pixels (can be less than chunk size).
	Width float64

	// Chunk height in pixels (can be less than chunk size).
	Height float64

	WG *sync.WaitGroup
}

func (r *Rasterizer) DrawRectangle(ri RectangleInfo) {
	r.Run()

	if ri.Buffer == nil {
		panic("Rasterizer.DrawRectangle: buffer is not set")
	}
	if ri.Shader == nil {
		panic("Rasterizer.DrawRectangle: shader is not set")
	}

	if ri.ChunkBits == 0 {
		ri.ChunkBits = 3
	}

	if ri.W == 0 || ri.H == 0 {
		return
	}

	if ri.W < 0 {
		ri.X += ri.W
		ri.W = -ri.W
	}

	if ri.H < 0 {
		ri.Y += ri.H
		ri.H = -ri.H
	}

	// Structure that will contain values precomputed for all pixels.
	rs := RectangleRasterStatic{
		RectangleRasterInput: ri.RectangleRasterInput,
	}

	// Rounded to nearest pixel coords.
	sm := 1 - math.SmallestNonzeroFloat64
	x, y := math.Floor(ri.X+sm), math.Floor(ri.Y+sm)
	w, h := math.Floor(ri.W), math.Floor(ri.H)

	rs.Pxs = 1 / w
	rs.Pys = 1 / h

	// Clip.
	minX := math.Max(x, 0)
	minY := math.Max(y, 0)
	maxX := math.Min(x+w, float64(ri.BufferWidth-1))
	maxY := math.Min(y+h, float64(ri.BufferHeight-1))

	chunkSize := 1 << ri.ChunkBits
	chunkSizef := float64(chunkSize)

	// Steps for chunks.
	pxCh := chunkSizef * rs.Pxs
	pyCh := chunkSizef * rs.Pys

	// Buffer offset for top-left bounding box corner.
	boffso := int(minX) + int(minY)*ri.BufferWidth

	var wg sync.WaitGroup

	var py float64
	for y := minY; y < maxY; y += chunkSizef {
		var px float64
		boffs := boffso
		chunkHeight := maxY - y
		if chunkHeight > chunkSizef {
			chunkHeight = chunkSizef
		}

		for x := minX; x < maxX; x += chunkSizef {
			wg.Add(1)
			chunk := RectangleChunk{
				RectangleRasterStatic: rs,
				Xo:                    x,
				Yo:                    y,
				Pxo:                   px,
				Pyo:                   py,
				BufferOffset:          boffs,
				Width:                 maxX - x,
				Height:                chunkHeight,
				WG:                    &wg,
			}
			if chunk.Width > chunkSizef {
				chunk.Width = chunkSizef
			}

			r.renderChunk(chunk)
			px += pxCh
			boffs += chunkSize
		}
		py += pyCh
		boffso += ri.BufferWidth * chunkSize
	}
	wg.Wait()
}

func (c RectangleChunk) Process() {
	defer c.WG.Done()

	maxX := c.Xo + c.Width
	maxY := c.Yo + c.Height
	so := RectangleShaderOpts{
		RectangleRasterStatic: c.RectangleRasterStatic,
	}

	cnt := 0
	py := c.Pyo
	for y := c.Yo; y < maxY; y++ {
		px := c.Pxo
		boffs := c.BufferOffset
		for x := c.Xo; x < maxX; x++ {
			so.BufferOffset = boffs
			so.X = x
			so.Y = y
			so.Px = px
			so.Py = py
			c.Shader(&so)
			px += c.Pxs
			boffs++
			cnt++
		}
		py += c.Pys
		c.BufferOffset += c.BufferWidth
	}
}
