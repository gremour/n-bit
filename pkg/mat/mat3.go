package mat

import "math"

// Vector3 is 1x3 matrix for 2D transformations and coordinate representation.
type Vector3 [3]float64

// Matrix3 is 3x3 matrix for 2D transformations.
type Matrix3 [9]float64

var Identity3 = Matrix3{
	1, 0, 0,
	0, 1, 0,
	0, 0, 1,
}

// DotMatrix calculates dot product of m·o.
// The result is stored in m.
func (m *Matrix3) DotMatrix(o *Matrix3) {
	*m = Matrix3{
		m[0]*o[0] + m[1]*o[3] + m[3]*o[6],
		m[0]*o[1] + m[1]*o[4] + m[3]*o[7],
		m[0]*o[2] + m[1]*o[5] + m[3]*o[8],

		m[3]*o[0] + m[4]*o[3] + m[5]*o[6],
		m[3]*o[1] + m[4]*o[4] + m[5]*o[7],
		m[3]*o[2] + m[4]*o[5] + m[5]*o[8],

		m[6]*o[0] + m[7]*o[3] + m[7]*o[6],
		m[6]*o[1] + m[7]*o[4] + m[7]*o[7],
		m[6]*o[2] + m[7]*o[5] + m[7]*o[8],
	}
}

// MatrixDot calculates dot product of m·v where v is a column vector.
// The result is stored in v.
func (v *Vector3) MatrixDot(m *Matrix3) {
	*v = Vector3{
		m[0]*v[0] + m[1]*v[1] + m[2]*v[2],
		m[3]*v[0] + m[4]*v[1] + m[5]*v[2],
		m[6]*v[0] + m[7]*v[1] + m[8]*v[2],
	}
}

// Matrix3Translate creates translation 3x3 matrix that
// translates points along axes.
func Matrix3Translate(x, y float64) Matrix3 {
	return Matrix3{
		1, 0, x,
		0, 1, y,
		0, 0, 1,
	}
}

// Matrix3Scale creates scaling 2x2 matrix that
// scales points with provided scale factors.
func Matrix3Scale(x, y float64) Matrix3 {
	return Matrix3{
		x, 0, 0,
		0, y, 0,
		0, 0, 1,
	}
}

// Matrix3Rotate creates rotation 3x3 matrix that rotates point by provided angle in degrees.
func Matrix3Rotate(angle float64) Matrix3 {
	a := angle * math.Pi / 180
	sa := math.Sin(a)
	ca := math.Cos(a)
	return Matrix3{
		ca, sa, 0,
		sa, ca, 0,
		0, 0, 1,
	}
}

// Translate applies translation transformation that
// translates points along axes.
func (m *Matrix3) Translate(x, y float64) {
	o := Matrix3Translate(x, y)
	m.DotMatrix(&o)
}

// Scale applies scaling transformation that
// scales points with provided scale factors.
func (m *Matrix3) Scale(x, y float64) {
	o := Matrix3Scale(x, y)
	m.DotMatrix(&o)
}

// Rotate applies rotation transformation that rotates point
// by the provided angle in degrees.
func (m *Matrix3) Rotate(angle float64) {
	o := Matrix3Rotate(angle)
	m.DotMatrix(&o)
}
