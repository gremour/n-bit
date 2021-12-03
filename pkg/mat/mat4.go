package mat

import (
	"math"
)

// Vector4 is 1x4 matrix for 3D transformations and coordinate representation.
type Vector4 [4]float64

// Matrix4 is 4x4 matrix for 3D transformations.
type Matrix4 [16]float64

var Identity4 = Matrix4{
	1, 0, 0, 0,
	0, 1, 0, 0,
	0, 0, 1, 0,
	0, 0, 0, 1,
}

// DotMatrix calculates dot product of m·o.
// The result is stored in m.
func (m *Matrix4) DotMatrix(o *Matrix4) {
	*m = Matrix4{
		m[0]*o[0] + m[1]*o[4] + m[2]*o[8] + m[3]*o[12],
		m[0]*o[1] + m[1]*o[5] + m[2]*o[9] + m[3]*o[13],
		m[0]*o[2] + m[1]*o[6] + m[2]*o[10] + m[3]*o[14],
		m[0]*o[3] + m[1]*o[7] + m[2]*o[11] + m[3]*o[15],

		m[4]*o[0] + m[5]*o[4] + m[6]*o[8] + m[7]*o[12],
		m[4]*o[1] + m[5]*o[5] + m[6]*o[9] + m[7]*o[13],
		m[4]*o[2] + m[5]*o[6] + m[6]*o[10] + m[7]*o[14],
		m[4]*o[3] + m[5]*o[7] + m[6]*o[11] + m[7]*o[15],

		m[8]*o[0] + m[9]*o[4] + m[10]*o[8] + m[11]*o[12],
		m[8]*o[1] + m[9]*o[5] + m[10]*o[9] + m[11]*o[13],
		m[8]*o[2] + m[9]*o[6] + m[10]*o[10] + m[11]*o[14],
		m[8]*o[3] + m[9]*o[7] + m[10]*o[11] + m[11]*o[15],

		m[12]*o[0] + m[13]*o[4] + m[14]*o[8] + m[15]*o[12],
		m[12]*o[1] + m[13]*o[5] + m[14]*o[9] + m[15]*o[13],
		m[12]*o[2] + m[13]*o[6] + m[14]*o[10] + m[15]*o[14],
		m[12]*o[3] + m[13]*o[7] + m[14]*o[11] + m[15]*o[15],
	}
}

// MatrixDot calculates dot product of m·v where v is a column vector.
// The result is stored in v.
func (v *Vector4) MatrixDot(m *Matrix4) {
	*v = Vector4{
		m[0]*v[0] + m[1]*v[1] + m[2]*v[2] + m[3]*v[3],
		m[4]*v[0] + m[5]*v[1] + m[6]*v[2] + m[7]*v[3],
		m[8]*v[0] + m[9]*v[1] + m[10]*v[2] + m[11]*v[3],
		m[12]*v[0] + m[13]*v[1] + m[14]*v[2] + m[15]*v[3],
	}
}

// Translate applies translation transformation that
// translates points along axes.
func (m *Matrix4) Translate(x, y, z float64) {
	o := Matrix4Translate(x, y, z)
	m.DotMatrix(&o)
}

// Scale applies scaling transformation that
// scales points with provided scale factors.
func (m *Matrix4) Scale(x, y, z float64) {
	o := Matrix4Scale(x, y, z)
	m.DotMatrix(&o)
}

// Rotate applies rotation transformation that rotates point around unit
// vector (x, y, z) by the provided angle in degrees.
// Vector must be normalized.
func (m *Matrix4) Rotate(angle, x, y, z float64) {
	o := Matrix4Rotate(angle, x, y, z)
	m.DotMatrix(&o)
}

// Matrix4Translate creates translation 4x4 matrix that
// translates points along axes.
func Matrix4Translate(x, y, z float64) Matrix4 {
	return Matrix4{
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, z,
		0, 0, 0, 1,
	}
}

// Matrix4Scale creates scaling 4x4 matrix that
// scales points with provided scale factors.
func Matrix4Scale(x, y, z float64) Matrix4 {
	return Matrix4{
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1,
	}
}

// Matrix4Rotate creates rotation 4x4 matrix that rotates point around unit
// vector (x, y, z) by the provided angle in degrees.
// Vector must be normalized.
func Matrix4Rotate(angle, x, y, z float64) Matrix4 {
	a := angle * math.Pi / 180
	sa := math.Sin(a)
	ca := math.Cos(a)
	ca1 := 1 - ca
	xy1 := x * y * ca1
	xz1 := x * z * ca1
	yz1 := y * z * ca1
	x1 := x * sa
	y1 := y * sa
	z1 := z * sa
	return Matrix4{
		ca + x*x*ca1, xy1 - z1, xz1 + y1, 0,
		xy1 + z1, ca + y*y*ca1, yz1 - x1, 0,
		xz1 - y1, yz1 + x1, ca + z*z*ca1, 0,
		0, 0, 0, 1,
	}
}

// Matrix4Frustum creates frustrum projection 4x4 matrix.
func Matrix4Frustum(l, r, t, b, n, f float64) Matrix4 {
	dx1 := 1 / (r - l)
	dy1 := 1 / (b - t)
	dz1 := 1 / (f - n)
	n2 := n * 2
	return Matrix4{
		n2 * dx1, 0, (r + l) * dx1, 0,
		0, n2 * dy1, (t + b) * dy1, 0,
		0, 0, -(f + n) * dz1, -f * n2 * dz1,
		0, 0, -1, 0,
	}
}

// Matrix4Ortho creates orthographic projection 4x4 matrix.
func Matrix4Ortho(l, r, t, b, n, f float64) Matrix4 {
	dx1 := 1 / (r - l)
	dy1 := 1 / (b - t)
	dz1 := 1 / (f - n)
	return Matrix4{
		2 * dx1, 0, 0, (r + l) * dx1,
		0, 2 * n * dy1, 0, (t + b) * dy1,
		0, 0, -2 * dz1, (f + n) * dz1,
		0, 0, 0, 1,
	}
}

// FovToY returns y coordinate to pass to Matrix4Frustum as
// (-y, y, -y, y) for given field of view (in degrees) and near coordinate.
// Example usage:
// near := 1
// far := 10
// fov := 120
// y := FovToY(fov, near)
// m := Matrix4Frustum(-y, y, -y, y, near, far)
func FovToY(fov, n float64) float64 {
	return n * math.Tan(fov*math.Pi/360)
}

// Normalize returns normalized vector of length 1.
func Normalize(x, y, z float64) (xn, yn, zn float64) {
	l := math.Sqrt(x*x + y*y + z*z)
	return x / l, y / l, z / l
}

func (v *Vector4) WDiv() {
	v[0] /= v[3]
	v[1] /= v[3]
	v[2] /= v[3]
	v[3] = 1 / v[3]
}

func (v *Vector4) ViewportProject(vp *Viewport) {
	v[0] = vp.Left + vp.w2 + vp.w2*v[0]
	v[1] = vp.Top + vp.h2 + vp.h2*v[1]
	v[2] = vp.fpn2 + vp.fmn2*v[2]
	v[2] = 1 / (v[2] * v[2])
}
