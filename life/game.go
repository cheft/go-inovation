package life

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 2560
	ScreenHeight = 1920
	ScreenScale  = 0.3
	Title        = "独立人生"
)

func init() {
}

type Game struct {
	transparent bool
	human       *Sprite
}

func NewGame() (*Game, error) {
	human := NewSprite()
	return &Game{transparent: false, human: human}, nil
}

func (g *Game) SetTransparent() {
	g.transparent = true
}

func (g *Game) Update() error {
	g.human.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	DrawMap(screen)
	g.human.Draw(screen)
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	// Text(screen, "中文支持1234567890abcdwert", 0, 200)
	// TextByBotw(screen, "中文支持1234567890abcdwert", 0, 228)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
