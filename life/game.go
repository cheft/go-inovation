package life

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	Title        = "life"
	ScreenWidth  = 320
	ScreenHeight = 240
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
	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
