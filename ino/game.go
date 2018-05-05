package ino

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"runtime/pprof"
	"strconv"

	"github.com/gopherjs/gopherjs/js"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"

	"github.com/hajimehoshi/go-inovation/ino/internal/assets"
	"github.com/hajimehoshi/go-inovation/ino/internal/font"
	"github.com/hajimehoshi/go-inovation/ino/internal/input"
)

type Game struct {
	resourceLoadedCh chan error
	gameState        GameState
	gameData         *GameData
	img              map[string]*ebiten.Image
	screen           *ebiten.Image
	cpup             *os.File
	sourceRect       image.Rectangle
}

var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")

func (g *Game) Loop(screen *ebiten.Image) error {
	// exp
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) && js.Global != nil {
		doc := js.Global.Get("document")
		canvas := doc.Call("getElementsByTagName", "canvas").Index(0)
		context := canvas.Call("getContext", "webgl")
		context.Call("getExtension", "WEBGL_lose_context").Call("loseContext")
		fmt.Println("Context Lost!")
		return nil
	}
	
	if g.resourceLoadedCh != nil {
		select {
		case err := <-g.resourceLoadedCh:
			if err != nil {
				return err
			}
			g.resourceLoadedCh = nil
		default:
		}
	}
	if g.resourceLoadedCh != nil {
		ebitenutil.DebugPrint(screen, "Now Loading...")
		return nil
	}

	input.Current().Update()

	if input.Current().IsKeyJustPressed(ebiten.KeyF) {
		f := ebiten.IsFullscreen()
		ebiten.SetFullscreen(!f)
		ebiten.SetCursorVisible(f)
	}

	if input.Current().IsKeyJustPressed(ebiten.KeyP) && *cpuProfile != "" && g.cpup == nil {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			panic(err)
		}
		g.cpup = f
		pprof.StartCPUProfile(f)
		fmt.Println("Start CPU Profiling")
	}

	if input.Current().IsKeyJustPressed(ebiten.KeyQ) && g.cpup != nil {
		pprof.StopCPUProfile()
		g.cpup.Close()
		g.cpup = nil
		fmt.Println("Stop CPU Profiling")
	}

	g.screen = screen
	if g.gameState == nil {
		g.gameState = &TitleMain{}
	} else {
		switch g.gameState.Msg() {
		case GAMESTATE_MSG_REQ_TITLE:
			if err := PauseBGM(); err != nil {
				return err
			}
			g.gameState = &TitleMain{}
		case GAMESTATE_MSG_REQ_OPENING:
			if err := PlayBGM(BGM1); err != nil {
				return err
			}
			g.gameState = &OpeningMain{}
		case GAMESTATE_MSG_REQ_GAME:
			g.gameState = NewGameMain(g)
		case GAMESTATE_MSG_REQ_ENDING:
			if err := PlayBGM(BGM1); err != nil {
				return err
			}
			g.gameState = &EndingMain{}
		case GAMESTATE_MSG_REQ_SECRET1:
			if err := PlayBGM(BGM1); err != nil {
				return err
			}
			g.gameState = NewSecretMain(1)
		case GAMESTATE_MSG_REQ_SECRET2:
			if err := PlayBGM(BGM1); err != nil {
				return err
			}
			g.gameState = NewSecretMain(2)
		}
	}
	g.gameState.Update(g)
	if !ebiten.IsRunningSlowly() {
		g.gameState.Draw(g)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nFPS: %.2f", ebiten.CurrentFPS()))
	return nil
}

var (
	imageItemFrame *ebiten.Image
)

func init() {
	imageItemFrame, _ = ebiten.NewImage(32, 32, ebiten.FilterNearest)
	imageItemFrame.Fill(color.Black)
	ebitenutil.DrawRect(imageItemFrame, 2, 2, 28, 28, color.White)
}

func (g *Game) DrawItemFrame(x, y int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	g.screen.DrawImage(imageItemFrame, op)
}

func (g *Game) Draw(key string, px, py, sx, sy, sw, sh int) {
	op := &ebiten.DrawImageOptions{}
	g.sourceRect.Min.X = sx
	g.sourceRect.Min.Y = sy
	g.sourceRect.Max.X = sx + sw
	g.sourceRect.Max.Y = sy + sh
	op.SourceRect = &g.sourceRect
	op.GeoM.Translate(float64(px), float64(py))
	g.screen.DrawImage(g.img[key], op)
}

func (g *Game) DrawNumber(num int, x, y int) {
	font.DrawText(g.screen, strconv.Itoa(num), x, y)
}

func (g *Game) DrawTouchButtons() {
	img := g.img["touch"]
	w, h := img.Size()
	w /= 4
	dx := 0
	dy := ScreenHeight - h
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 0.4)
	for _, i := range []int{0, 1, 3} {
		r := image.Rect(i*w, 0, (i+1)*w, h)
		op.SourceRect = &r
		op.GeoM.Reset()
		op.GeoM.Translate(float64(dx+i*w), float64(dy))
		g.screen.DrawImage(img, op)
	}
	// Render 'down' button
	op = &ebiten.DrawImageOptions{}
	r := image.Rect(2*w, 0, 3*w, h)
	op.SourceRect = &r
	op.GeoM.Translate(float64(dx+2*w), float64(dy))
	alpha := 0.0
	if input.Current().IsActionKeyPressed() {
		alpha = 0.4
	} else {
		alpha = 0.1
	}
	op.ColorM.Scale(1, 1, 1, alpha)
	g.screen.DrawImage(img, op)
}

func (g *Game) loadImages() error {
	for _, f := range []string{"ino", "msg", "bg", "touch"} {
		b, err := assets.Asset("resources/images/color/" + f + ".png")
		if err != nil {
			return err
		}
		origImg, _, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			return err
		}
		img, err := ebiten.NewImageFromImage(origImg, ebiten.FilterNearest)
		if err != nil {
			return err
		}
		g.img[f] = img
	}

	fontImgs := map[rune]*ebiten.Image{}
	for n := 48; n <= 57; n++ {
		b, err := assets.Asset(fmt.Sprintf("resources/font/%d.png", n))
		if err != nil {
			return err
		}
		origImg, _, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			return err
		}
		img, err := ebiten.NewImageFromImage(origImg, ebiten.FilterNearest)
		if err != nil {
			return err
		}
		fontImgs[rune(n)] = img
	}
	return nil
}

func NewGame() (*Game, error) {
	game := &Game{
		img:              map[string]*ebiten.Image{},
		resourceLoadedCh: make(chan error),
	}
	go func() {
		if err := game.loadImages(); err != nil {
			game.resourceLoadedCh <- err
			return
		}
		if err := loadAudio(); err != nil {
			game.resourceLoadedCh <- err
			return
		}
		close(game.resourceLoadedCh)
	}()
	if err := finalizeAudio(); err != nil {
		return nil, err
	}
	return game, nil
}
