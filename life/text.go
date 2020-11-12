package life

import (
	"image/color"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	ipixFace font.Face
	botwFace font.Face
)

func init() {
	ipixFace = loadTTF("life/assets/font/IPix.ttf")
	botwFace = loadTTF("life/assets/font/Botw_cn_zyt.ttf")
}

func Text(screen *ebiten.Image, msg string, x int, y int) {
	text.Draw(screen, msg, ipixFace, x, y, color.White)
}

func TextByBotw(screen *ebiten.Image, msg string, x int, y int) {
	text.Draw(screen, msg, botwFace, x, y, color.White)
}

func loadTTF(path string) font.Face {
	ttf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err, "readFile")
	}

	tt, err := opentype.Parse(ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	return face
}
