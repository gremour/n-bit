package display

import (
	"math"
	"sync"
)

// TriangleRasterInput part of options that is passed to rasterizer workers.
type TriangleRasterInput struct {
	// Pixel buffer. Shader needs to output color to buffer with
	// boffs offset and value calculated by shader. Buffer can use multiple
	// bytes for pixel. Multiply boffs by pixel size in this case.
	Buffer       []byte
	BufferWidth  int
	BufferHeight int

	// Function to call for each fragment.
	// Shader is responsible for saving pixel to the buffer. Use lambda function that
	// has access to the buffer to do that. See bit-shader for example.
	Shader func(o *TriangleShaderOpts)

	// Raster chunks size (in bits, e.g. 3 bits = 8 pixels).
	ChunkBits int

	// Lights model.
	Lights Lights
	// Object to convert to index color.
	Indexizer Indexizer

	// Any other data to pass to shader.
	Extra interface{}
}

// TriangleRasterStatic part of options including input and precalculated for the whole triangle.
type TriangleRasterStatic struct {
	TriangleRasterInput

	// Barycentric coords steps.
	A01, B01 float64
	A12, B12 float64
	A20, B20 float64

	// 1 / area.
	InvArea float64
}

// TriangleShaderOpts options for shader.
type TriangleShaderOpts struct {
	TriangleRasterStatic

	// Cooridnates in buffer, pixels.
	X, Y float64

	// Barycentric coords.
	W0, W1, W2 float64

	// Offset in buffer, pixels. Equal to int(x) + int(y) * BufferWidth.
	BufferOffset int
}

// TriangleInfo input data to DrawTriangle call.
type TriangleInfo struct {
	TriangleRasterInput

	X0, Y0 float64
	X1, Y1 float64
	X2, Y2 float64
}

// TriangleChunk data for triangle chunk to pass to rasterizer worker.
type TriangleChunk struct {
	TriangleRasterStatic

	// Origin x, y, pixels.
	Xo, Yo float64

	// Origin barycentric coords.
	W0o, W1o, W2o float64

	// Offset in buffer, pixels. Equal to x + y * BufferWidth.
	BufferOffset int

	// Chunk width in pixels (can be less than chunk size).
	Width float64

	// Chunk height in pixels (can be less than chunk size).
	Height float64

	WG *sync.WaitGroup
}

func (r *Rasterizer) DrawTriangle(ti TriangleInfo) {
	r.Run()

	if ti.Buffer == nil {
		panic("Rasterizer.DrawTriangle: buffer is not set")
	}
	if ti.Shader == nil {
		panic("Rasterizer.DrawTriangle: shader is not set")
	}
	if ti.Lights == nil {
		panic("Rasterizer.DrawTriangle: lights is not set")
	}
	if ti.Indexizer == nil {
		panic("Rasterizer.DrawTriangle: indexizer is not set")
	}

	if ti.ChunkBits == 0 {
		ti.ChunkBits = 3
	}

	// Structure that will contain values precomputed for all pixels.
	rs := TriangleRasterStatic{
		TriangleRasterInput: ti.TriangleRasterInput,
	}

	// Rounded to nearest pixel coords.
	sm := 1 - math.SmallestNonzeroFloat64
	x0, y0 := math.Floor(ti.X0+sm), math.Floor(ti.Y0+sm)
	x1, y1 := math.Floor(ti.X1+sm), math.Floor(ti.Y1+sm)
	x2, y2 := math.Floor(ti.X2+sm), math.Floor(ti.Y2+sm)

	// Barycentric steps.
	rs.A01, rs.B01 = ti.Y0-ti.Y1, ti.X1-ti.X0
	rs.A12, rs.B12 = ti.Y1-ti.Y2, ti.X2-ti.X1
	rs.A20, rs.B20 = ti.Y2-ti.Y0, ti.X0-ti.X2

	// Clip.
	minX := math.Min(x0, math.Min(x1, x2))
	minY := math.Min(y0, math.Min(y1, y2))
	maxX := math.Max(x0, math.Max(x1, x2))
	maxY := math.Max(y0, math.Max(y1, y2))

	minX = math.Max(minX, 0)
	minY = math.Max(minY, 0)
	maxX = math.Min(maxX, float64(ti.BufferWidth-1))
	maxY = math.Min(maxY, float64(ti.BufferHeight-1))

	chunkSize := 1 << ti.ChunkBits
	chunkSizef := float64(chunkSize)

	area := edgeFunc(ti.X0, ti.Y0, ti.X1, ti.Y1, ti.X2, ti.Y2)
	if area == 0 {
		return
	}
	// Barycentric coords of origin (top-left corner of box surrounding the triangle).
	w0o := edgeFunc(ti.X1, ti.Y1, ti.X2, ti.Y2, minX, minY)
	w1o := edgeFunc(ti.X2, ti.Y2, ti.X0, ti.Y0, minX, minY)
	w2o := edgeFunc(ti.X0, ti.Y0, ti.X1, ti.Y1, minX, minY)

	// Top-left fill correction.
	if !isTopLeft(x1, y1, x2, y2) {
		w0o -= math.SmallestNonzeroFloat64
	}
	if !isTopLeft(x2, y2, x0, y0) {
		w1o -= math.SmallestNonzeroFloat64
	}
	if !isTopLeft(x0, y0, x1, y1) {
		w2o -= math.SmallestNonzeroFloat64
	}

	rs.InvArea = 1 / area

	// Barycentric steps for chunks.
	w0xStep := rs.A12 * chunkSizef
	w1xStep := rs.A20 * chunkSizef
	w2xStep := rs.A01 * chunkSizef
	w0yStep := rs.B12 * chunkSizef
	w1yStep := rs.B20 * chunkSizef
	w2yStep := rs.B01 * chunkSizef

	// Buffer offset for top-left bounding box corner.
	boffso := int(minX) + int(minY)*ti.BufferWidth

	var wg sync.WaitGroup

	for y := minY; y < maxY; y += chunkSizef {
		w0, w1, w2 := w0o, w1o, w2o
		boffs := boffso
		chunkHeight := maxY - y
		if chunkHeight > chunkSizef {
			chunkHeight = chunkSizef
		}

		for x := minX; x < maxX; x += chunkSizef {
			wg.Add(1)
			chunk := TriangleChunk{
				TriangleRasterStatic: rs,
				Xo:                   x,
				Yo:                   y,
				W0o:                  w0,
				W1o:                  w1,
				W2o:                  w2,
				BufferOffset:         boffs,
				Width:                maxX - x,
				Height:               chunkHeight,
				WG:                   &wg,
			}
			if chunk.Width > chunkSizef {
				chunk.Width = chunkSizef
			}

			r.renderChunk(chunk)
			w0 += w0xStep
			w1 += w1xStep
			w2 += w2xStep
			boffs += chunkSize
		}
		w0o += w0yStep
		w1o += w1yStep
		w2o += w2yStep
		boffso += ti.BufferWidth * chunkSize
	}
	wg.Wait()
}

func (c TriangleChunk) Process() {
	defer c.WG.Done()

	maxX := c.Xo + c.Width
	maxY := c.Yo + c.Height
	so := TriangleShaderOpts{
		TriangleRasterStatic: c.TriangleRasterStatic,
	}

	cnt := 0
	for y := c.Yo; y < maxY; y++ {
		w0, w1, w2 := c.W0o, c.W1o, c.W2o
		boffs := c.BufferOffset
		for x := c.Xo; x < maxX; x++ {
			if w0 >= 0 && w1 >= 0 && w2 >= 0 {
				so.BufferOffset = boffs
				so.W0 = w0 * c.InvArea
				so.W1 = w1 * c.InvArea
				so.W2 = w2 * c.InvArea
				so.X = x
				so.Y = y
				c.Shader(&so)
			}
			w0 += c.A12
			w1 += c.A20
			w2 += c.A01
			boffs++
			cnt++
		}
		c.W0o += c.B12
		c.W1o += c.B20
		c.W2o += c.B01
		c.BufferOffset += c.BufferWidth
	}
}

// edgeFunc calculates triangle edge function.
func edgeFunc(ax, ay, bx, by, cx, cy float64) float64 {
	return (bx-ax)*(cy-ay) - (by-ay)*(cx-ax)
}

func isTopLeft(ax, ay, bx, by float64) bool {
	top := ay == by && ax < bx
	left := ay > by
	return top || left
}
