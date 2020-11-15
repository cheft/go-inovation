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
	speed       = 2
)

var (
	spriteImage *ebiten.Image
)

type Sprite struct {
	count  int
	x      float64
	y      float64
	isMove bool
}

func NewSprite() *Sprite {
	// img, _, _ := image.Decode(bytes.NewReader(images.Runner_png))
	// spriteImage = ebiten.NewImageFromImage(img)
	spriteImage, _, _ = ebitenutil.NewImageFromFile("life/assets/img/human1.png")
	return &Sprite{count: 0, x: 100, y: 100, isMove: false}
}

func (s *Sprite) Update() error {
	s.count++
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		s.x -= speed
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		s.x += speed
	} else if ebiten.IsKeyPressed(ebiten.KeyW) {
		s.y -= speed
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		s.y += speed
	} else {

	}
	return nil
}

func (s *Sprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// 居中
	// op.GeoM.Translate(ScreenWidth/2, ScreenHeight/2)
	// op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Translate(s.x, s.y)
	i := (s.count / 10) % frameNum
	sx, sy := frameOX+i*frameWidth, frameOY
	screen.DrawImage(spriteImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}
