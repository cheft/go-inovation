package life

import (
	"go-inovation/life/utils"
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	tileSize = 32
	tileXNum = 19
)

var (
	tilesImage *ebiten.Image
	m          *utils.Map
)

func chk(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	tilesImage, _, _ = ebitenutil.NewImageFromFile("life/assets/tmx/Serene_Village_32x32.png")

	r, err := os.Open("life/assets/tmx/farm.tmx")
	chk(err)

	m, err = utils.Decode(r)
	chk(err)
}

func DrawMap(screen *ebiten.Image) {
	// const xNum = ScreenWidth / tileSize
	const xNum = 80
	for _, l := range m.Layers {
		for i, t := range l.GIDs {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64((i%xNum)*tileSize), float64((i/xNum)*tileSize))
			sx := (int(t) % tileXNum) * tileSize
			sy := (int(t) / tileXNum) * tileSize
			screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
		}
	}
}
