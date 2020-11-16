package life

import (
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	frameWidth  = 32
	frameHeight = 48
)

var (
	spriteImage *ebiten.Image
)

type Sprite struct {
	count    int
	x        float64
	y        float64
	isMove   bool
	speed    float64
	frameNum int
	frameOX  int
	frameOY  int
}

func NewSprite() *Sprite {
	spriteImage, _, _ = ebitenutil.NewImageFromFile("life/assets/img/human1.png")
	return &Sprite{
		count:    0,
		x:        100,
		y:        100,
		isMove:   false,
		speed:    2,
		frameNum: 3,
		frameOX:  96, // 人物切换
		frameOY:  0,
	}
}

func (s *Sprite) Update() error {
	s.count++
	s.isMove = true
	s.frameNum = 3
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		s.x -= s.speed
		s.frameOY = 48
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		s.x += s.speed
		s.frameOY = 96
	} else if ebiten.IsKeyPressed(ebiten.KeyW) {
		s.y -= s.speed
		s.frameOY = 144
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		s.y += s.speed
		s.frameOY = 0
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if s.frameOX <= 0 {
			s.frameOX = 288
		} else {
			s.frameOX -= 96
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if s.frameOX >= 288 {
			s.frameOX = 0
		} else {
			s.frameOX += 96
		}
	} else {
		s.isMove = false
		s.frameNum = 1
		s.count = 0
	}
	return nil
}

func (s *Sprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(s.x, s.y)
	i := (s.count / 10) % s.frameNum
	sx, sy := s.frameOX+i*frameWidth, s.frameOY
	screen.DrawImage(spriteImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}
