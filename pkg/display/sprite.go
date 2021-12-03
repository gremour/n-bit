package display

import (
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"log"
	"strings"

	"os"

	"gopkg.in/yaml.v3"
)

type Sprite struct {
	Atlas            *IndexedImage
	X, Y             int
	Width, Height    int
	XOrigin, YOrigin int
}

type altasYaml struct {
	Name    string        `yaml:"name"`
	File    string        `yaml:"file"`
	Sprites []spritesYaml `yaml:"sprites"`
}

type spritesYaml struct {
	Width  int      `yaml:"width"`
	Height int      `yaml:"height"`
	XOffs  int      `yaml:"xoffs"`
	YOffs  int      `yaml:"yoffs"`
	XOrig  int      `yaml:"xorig"`
	YOrig  int      `yaml:"yorig"`
	Names  []string `yaml:"names"`
}

func (d *Display) LoadAtlas(fileName string) error {
	if d.Atlases == nil {
		d.Atlases = make(map[string]*IndexedImage)
	}
	if d.Sprites == nil {
		d.Sprites = make(map[string]*Sprite)
	}
	pl, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	var atl altasYaml
	err = yaml.Unmarshal(pl, &atl)
	if err != nil {
		return err
	}
	atl.File = strings.ReplaceAll(atl.File, `\`, `/`)

	im, err := LoadImage(atl.File)
	if err != nil {
		return fmt.Errorf("failed to open image %v: %v", fileName, err)
	}
	tex := IndexedImageFromImage(im, FromImageOpts{})

	if atl.Name == "" {
		p1 := strings.Split(atl.File, "/")
		p2 := strings.Split(p1[len(p1)-1], ".")
		atl.Name = p2[0]
	}
	d.Atlases[atl.Name] = &tex
	for _, s := range atl.Sprites {
		d.addSprites(&tex, &s, atl.Name)
	}

	return nil
}

func (d *Display) addSprites(atl *IndexedImage, spr *spritesYaml, pref string) {
	x, y := spr.XOffs, spr.YOffs
	for _, n := range spr.Names {
		s := &Sprite{
			Atlas:   atl,
			X:       x,
			Y:       y,
			Width:   spr.Width,
			Height:  spr.Height,
			XOrigin: spr.XOrig,
			YOrigin: spr.YOrig,
		}
		var last bool
		x += spr.Width
		if x >= atl.Width {
			s.Width = atl.Width - s.X
			x = 0
			y += s.Height
			if y >= atl.Height {
				s.Height = atl.Height - s.Y
				last = true
			}
		}
		d.Sprites[pref+"."+n] = s
		if last {
			break
		}
	}
}

func (d *Display) ReportSprite(name string) {
	if d.reportedSprite == nil {
		d.reportedSprite = make(map[string]struct{})
	}
	if _, ok := d.reportedSprite[name]; ok {
		return
	}
	d.reportedSprite[name] = struct{}{}
	log.Printf("Error: sprite %v is not loaded. Sprite name format is 'atlas.sprite'.", name)
}

func LoadImage(fileName string) (image.Image, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}
