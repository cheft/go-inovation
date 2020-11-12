package life

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 960
	ScreenHeight = 540
	ScreenScale  = 1
	Title        = "精彩人生"
)

func init() {
}

func NewGame() (*Game, error) {
	return &Game{}, nil
}

type Game struct {
	transparent bool
}

func (g *Game) SetTransparent() {
	g.transparent = true
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "1234567890abcdwert")
	Text(screen, "中文支持1234567890abcdwert", 0, 40)
	TextByBotw(screen, "中文支持1234567890abcdwert", 0, 70)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
