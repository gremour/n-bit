package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/gremour/n-bit/pkg/display"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game implements ebiten.Game interface.
type Game struct {
	Display   display.Display
	Lights    display.LightSet
	DeltaTime float64

	lights []*light

	LastTime time.Time
	CurTime  time.Time
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	g.CurTime = time.Now()
	g.DeltaTime = float64(g.CurTime.Sub(g.LastTime)) / float64(time.Second)
	if g.LastTime.IsZero() {
		g.DeltaTime = 0
	}
	defer func() { g.LastTime = g.CurTime }()

	for _, l := range g.lights {
		l.process(g.DeltaTime)
	}

	return nil
}

// Draw draws the game screen.
// Draw is called every frame.
func (g *Game) Draw(screen *ebiten.Image) {
	g.Display.Screen.Fill(0)

	// start := time.Now()
	// for i := 0; i < 100; i++ {
	g.Display.DrawSprite("spr.one", 130, 50)
	// g.Display.DrawSpriteAdvanced(display.DrawSpriteOpts{
	// 	Name: "av.one",
	// 	DX:   130,
	// 	DY:   0,
	// })
	// }
	// elapsed := time.Since(start)
	// fmt.Println(elapsed)

	for _, l := range g.lights {
		x := int(l.x)
		y := int(l.y)
		if x < 1 {
			x = 1
		}
		if x > 531 {
			x = 531
		}
		if y < 1 {
			y = 1
		}
		if y > 298 {
			y = 298
		}

		g.Display.Screen.Pixels[x+y*g.Display.Screen.Width] = 4
		g.Display.Screen.Pixels[x-1+y*g.Display.Screen.Width] = 4
		g.Display.Screen.Pixels[x+1+y*g.Display.Screen.Width] = 4
		g.Display.Screen.Pixels[x+(y-1)*g.Display.Screen.Width] = 4
		g.Display.Screen.Pixels[x+(y+1)*g.Display.Screen.Width] = 4
	}

	screen.ReplacePixels(g.Display.Screen.ToRGBA(display.ToRGBAOpts{
		Pixels:  g.Display.RGBA,
		Palette: g.Display.Palette,
	}))

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %v, CPU: %v", int(ebiten.CurrentFPS()), runtime.NumCPU()))
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	h := 300
	w := h * outsideWidth / outsideHeight
	if g.Display.Screen.Width == 0 {
		g.Display.InitBuffers(w, h)
		fmt.Printf("Display size: %vx%v\n", w, h)
	}
	return w, h
}

func main() {
	rand.Seed(time.Now().UnixNano())
	pname := "painted-leather" //display.RandomPalette2BitName()
	pal := display.Palettes2Bit[pname]
	if pal == nil {
		log.Fatalf("Palette %v is not found", pname)
	}
	g := &Game{
		Display: display.Display{
			Palette:   pal,
			Indexizer: display.Bits2,
		},
		Lights: display.LightSet{
			MinScale:  0.2,
			MaxScale:  1.0,
			MinOffset: 0,
			MaxOffset: 1,
		},
	}
	g.Display.Lights = &g.Lights

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 6; i++ {
		l := &light{
			x: 200, y: 150,
			sx:   rand.Float64()*50 + 8,
			sy:   rand.Float64()*50 + 8,
			size: rand.Float64()*2 + 1,
		}
		if rand.Float64() < 0.5 {
			l.sx = -l.sx
		}
		if rand.Float64() < 0.5 {
			l.sy = -l.sy
		}
		g.lights = append(g.lights, l)
		g.Lights.TrackCircle(l, 0.7+rand.Float64()*0.5, 50+30*rand.Float64())
	}

	ebiten.SetWindowTitle("n-bit engine")
	ebiten.SetFullscreen(true)

	err := g.Display.LoadAtlas("spr.yaml")
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type light struct {
	x, y   float64
	sx, sy float64
	dead   bool
	size   float64
}

func (l *light) process(dt float64) {
	l.x = l.x + l.sx*dt
	l.y = l.y + l.sy*dt
	if l.x < 0 {
		l.sx *= -1
	}
	if l.x > 533 {
		l.sx *= -1
	}
	if l.y < 0 {
		l.sy *= -1
	}
	if l.y > 300 {
		l.sy *= -1
	}
}

func (l *light) Pos() (float64, float64) {
	return l.x, l.y
}

func (l *light) Alive() bool {
	return !l.dead
}

func (l *light) SizeMod() float64 {
	return l.size
}
