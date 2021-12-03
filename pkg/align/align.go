package align

type Align int

const (
	Center      Align = 0
	Left        Align = 1
	Right       Align = 2
	Top         Align = 4
	Bottom      Align = 8
	TopLeft     Align = Top | Left
	TopRight    Align = Top | Right
	BottomLeft  Align = Bottom | Left
	BottomRight Align = Bottom | Right
)

// RectInt returns top-left coordinates of the rectangle with
// size (w, h) aligned to the point (x, y).
func (a Align) RectInt(x, y, w, h int) (nx, ny int) {
	return (a & (Left | Right)).LineInt(x, w), (a & (Top | Bottom)).LineInt(y, h)
}

// LineInt returns left coordinate of the line aligned to the point
// (both horizontal and vertical constants work).
func (a Align) LineInt(x, w int) int {
	center := (a&Left != a&Right) || (a&Top != a&Bottom)
	left := (a&Left == Left) || (a&Top == Top)
	if left {
		return x
	}
	if center {
		return x - w/2
	}
	return x - w
}
