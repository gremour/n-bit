package display

// FullLight is a lighting model that produces a perfect lit environment.
var FullLight = FixedLight{1}

// Lights calculates light intensity in a given point of the screen.
// It gets color value (intesity) at given point in the range of 1 (black) to
// 255 (white) and modifies it to produce intensity in the range (0-1).
// This intensity is supposed to be changed by Indexizer into indexed color.
type Lights interface {
	Light(val byte, posx, posy int) (intens float64)
}

// LightSource provides information about lighting for each point of screen.
// The most common implementation is a point light source that has a diminishing
// light the farther it is located from the light source.
// When light source Alive method returns false, the light source will be dropped
// from the current model.
type LightSource interface {
	OffsetScale(posx, posy int) (offs float64, scale float64)
	Alive() bool
}

// LightSet is a lighting model that includes minimum and maximum
// scale and offset amd a number of light sources that add to scale
// and offset at each pixel position.
// Offset is added to original intensity (scaled to 0-1), then it's
// multiplied by scale.
// The more sources are used, the slower drawing functions are.
type LightSet struct {
	Sources   []LightSource
	MinScale  float64
	MaxScale  float64
	MinOffset float64
	MaxOffset float64
}

// Light calculates light intensity (0-1) at the point using current light model.
func (l *LightSet) Light(val byte, posx, posy int) (intens float64) {
	offs := l.MinOffset
	scale := l.MinScale

	for i := 0; i < len(l.Sources); i++ {
		if !l.Sources[i].Alive() {
			n := len(l.Sources) - 1
			l.Sources[i] = l.Sources[n]
			l.Sources[n] = nil
			l.Sources = l.Sources[:n]
		}
		of, sc := l.Sources[i].OffsetScale(posx, posy)
		offs += of
		scale += sc
	}
	if offs > l.MaxOffset {
		offs = l.MaxOffset
	} else if offs < l.MinOffset {
		offs = l.MinOffset
	}
	if scale > l.MaxScale {
		scale = l.MaxScale
	} else if scale < l.MinScale {
		scale = l.MinScale
	}

	return (float64(val)/255 + offs) * scale
}

// TrackCircle adds a circle light source tied to an object that
// conforms a Tracked interface.
func (l *LightSet) TrackCircle(t Tracked, intens, rfall float64) {
	cs := &CircleSource{
		Tracked:    t,
		Intensity:  intens,
		FallRadius: rfall,
	}
	l.Sources = append(l.Sources, cs)
}

// CircleSource is a point light source with circle area, which tracks an object
// position, size and alive status.
// At fall radius intensity is reduced by two times. The dependency is quadric.
type CircleSource struct {
	Tracked    Tracked
	Intensity  float64
	FallRadius float64
}

// Tracked is an interface that describes object that can be tracked by
// the light system. Required methods are Pos which gives current position
// of a light source, SizeMod that can control momentary radius of the light,
// and Alive property, which, when returns false, will cause the removing
// of the tracked object.
// When creating game, any objects that acts as a point light source, needs
// to implement this interface. When object is deleted (removed from the game),
// set an internal variable to return Alive equals to false, so light source
// could be unregistered.
type Tracked interface {
	Pos() (float64, float64)
	Alive() bool
	SizeMod() float64
}

func (s *CircleSource) OffsetScale(posx, posy int) (offs float64, scale float64) {
	x, y := s.Tracked.Pos()
	dx := float64(posx) - x
	dy := float64(posy) - y
	d2 := dx*dx + dy*dy
	fr2 := s.FallRadius * s.FallRadius * s.Tracked.SizeMod()
	scale += s.Intensity * fr2 / (d2*10 + fr2)
	return 0, scale
}

func (s *CircleSource) Alive() bool {
	al := s.Tracked.Alive()
	if !al {
		s.Tracked = nil
	}
	return al
}

// FixedLight returns a constant intensity of light regardless of screen position.
// Intensity of 1 will produce the same as image intensity, while lower values
// will result in a dimmer color.
type FixedLight struct {
	Intensity float64
}

func (l FixedLight) Light(val byte, posx, posy int) float64 {
	return l.Intensity * float64(val) / 255
}
