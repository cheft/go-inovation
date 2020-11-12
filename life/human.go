package life

import (
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	frameOX     = 0
	frameOY     = 0
	frameWidth  = 32
	frameHeight = 48
	frameNum    = 3
)

var (
	spriteImage *ebiten.Image
)

type Sprite struct {
	count int
}

func NewSprite() *Sprite {
	// img, _, _ := image.Decode(bytes.NewReader(images.Runner_png))
	// spriteImage = ebiten.NewImageFromImage(img)
	spriteImage, _, _ = ebitenutil.NewImageFromFile("life/assets/img/human1.png")
	return &Sprite{0}
}

func (s *Sprite) Update() error {
	s.count++
	return nil
}

func (s *Sprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Translate(ScreenWidth/2, ScreenHeight/2)
	i := (s.count / 10) % frameNum
	sx, sy := frameOX+i*frameWidth, frameOY
	screen.DrawImage(spriteImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}
