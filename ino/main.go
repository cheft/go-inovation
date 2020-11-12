package ino

import (
	"fmt"
	"image/color"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/text/language"

	"go-inovation/ino/internal/audio"
	"go-inovation/ino/internal/draw"
	"go-inovation/ino/internal/font"
	"go-inovation/ino/internal/input"
	"go-inovation/ino/internal/text"
)

const (
	ScreenWidth  = draw.ScreenWidth
	ScreenHeight = draw.ScreenHeight
	ScreenScale  = draw.ScreenScale
	Title        = "INNO VATION 2007 (Go version)"
)

const (
	ENDINGMAIN_STATE_STAFFROLL = iota
	ENDINGMAIN_STATE_RESULT
)

type GameStateMsg int

const (
	GAMESTATE_MSG_NONE GameStateMsg = iota
	GAMESTATE_MSG_REQ_TITLE
	GAMESTATE_MSG_REQ_GAME
	GAMESTATE_MSG_REQ_OPENING
	GAMESTATE_MSG_REQ_ENDING
	GAMESTATE_MSG_REQ_SECRET_COMMAND
	GAMESTATE_MSG_REQ_SECRET_CLEAR
)

type TitleScene struct {
	gameStateMsg  GameStateMsg
	timer         int
	offsetX       int
	offsetY       int
	lunkerMode    bool
	lunkerCommand int
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (t *TitleScene) Update(game *Game) {
	t.timer++
	if t.timer%5 == 0 {
		t.offsetX = rand.Intn(5) - 3
		t.offsetY = rand.Intn(5) - 3
	}

	if (input.Current().IsActionKeyJustPressed() || input.Current().IsSpaceJustTouched()) && t.timer > 5 {
		t.gameStateMsg = GAMESTATE_MSG_REQ_OPENING

		if t.lunkerMode {
			game.gameData = NewGameData(GAMEMODE_LUNKER)
		} else {
			game.gameData = NewGameData(GAMEMODE_NORMAL)
		}
	}

	// ランカー・モード・コマンド
	switch t.lunkerCommand {
	case 0, 1, 2, 6:
		if input.Current().IsKeyJustPressed(ebiten.KeyLeft) {
			t.lunkerCommand++
		} else if input.Current().IsKeyJustPressed(ebiten.KeyRight) || input.Current().IsKeyJustPressed(ebiten.KeyUp) || input.Current().IsKeyJustPressed(ebiten.KeyDown) {
			t.lunkerCommand = 0
		}
	case 3, 4, 5, 7:
		if input.Current().IsKeyJustPressed(ebiten.KeyRight) {
			t.lunkerCommand++
		} else if input.Current().IsKeyJustPressed(ebiten.KeyLeft) || input.Current().IsKeyJustPressed(ebiten.KeyUp) || input.Current().IsKeyJustPressed(ebiten.KeyDown) {
			t.lunkerCommand = 0
		}
	default:
		break
	}
	if t.lunkerCommand > 7 {
		t.lunkerCommand = 0
		t.lunkerMode = !t.lunkerMode
	}

	if input.Current().IsLanguageSwitcherPressed() {
		switch game.lang {
		case language.Japanese:
			game.lang = language.English
		case language.English:
			game.lang = language.Japanese
		}
	}
}

func (t *TitleScene) Draw(screen *ebiten.Image, game *Game) {
	textID := text.TextIDStart
	clr := color.Black
	if !game.transparent {
		if t.lunkerMode {
			draw.Draw(screen, "bg", 0, 0, 0, 240, draw.ScreenWidth, draw.ScreenHeight)
			textID = text.TextIDStartLunker
			clr = color.White
		} else {
			draw.Draw(screen, "bg", 0, 0, 0, 0, draw.ScreenWidth, draw.ScreenHeight)
			if input.Current().IsTouchEnabled() {
				textID = text.TextIDStartTouch
			}
		}
	}

	str := text.Get(game.lang, textID)
	x := (draw.ScreenWidth-font.Width(str))/2 + t.offsetX
	y := (draw.ScreenHeight-240)/2 + 160 + t.offsetY
	font.DrawText(screen, str, x, y, clr)

	// Draw the title.
	key := "msg_" + game.lang.String()
	draw.Draw(screen, key, (draw.ScreenWidth-256)/2, 32+(draw.ScreenHeight-240)/2, 0, 0, 256, 48)

	// Draw the language switcher.
	font.DrawText(screen, "Language", 320-48, 0, color.RGBA{0x80, 0x80, 0x80, 0xff})
}

func (t *TitleScene) Msg() GameStateMsg {
	return t.gameStateMsg
}

type OpeningScene struct {
	gameStateMsg GameStateMsg
	timer        int
}

const (
	OPENING_SCROLL_SPEED = 3
)

func (o *OpeningScene) Update(game *Game) {
	o.timer++

	if input.Current().IsActionKeyPressed() || input.Current().IsSpaceTouched() {
		o.timer += 20
	}
	scrollLen := font.Height(text.Get(game.lang, text.TextIDOpening)) + draw.ScreenHeight
	if o.timer/OPENING_SCROLL_SPEED > scrollLen {
		o.gameStateMsg = GAMESTATE_MSG_REQ_GAME
		audio.PauseBGM()
	}
}

func (o *OpeningScene) Draw(screen *ebiten.Image, game *Game) {
	if !game.transparent {
		draw.Draw(screen, "bg", 0, 0, 0, 480, 320, 240)
	}

	for i, line := range strings.Split(text.Get(game.lang, text.TextIDOpening), "\n") {
		x := (draw.ScreenWidth - font.Width(line)) / 2
		y := draw.ScreenHeight - (o.timer / OPENING_SCROLL_SPEED) + i*font.LineHeight
		font.DrawText(screen, line, x, y, color.Black)
	}
}

func (o *OpeningScene) Msg() GameStateMsg {
	return o.gameStateMsg
}

type EndingScene struct {
	gameStateMsg   GameStateMsg
	timer          int
	bgmFadingTimer int
	state          int
}

const (
	ENDING_SCROLL_SPEED = 3
)

func (e *EndingScene) Update(game *Game) {
	e.timer++
	switch e.state {
	case ENDINGMAIN_STATE_STAFFROLL:
		if input.Current().IsActionKeyPressed() || input.Current().IsSpaceTouched() {
			e.timer += 20
		}
		scrollLen := font.Height(text.Get(game.lang, text.TextIDEnding)) + draw.ScreenHeight
		if e.timer/ENDING_SCROLL_SPEED > scrollLen {
			e.timer = 0
			e.state = ENDINGMAIN_STATE_RESULT
		}
	case ENDINGMAIN_STATE_RESULT:
		const max = 5 * 60
		e.bgmFadingTimer++
		switch {
		case e.bgmFadingTimer == max:
			audio.PauseBGM()
		case e.bgmFadingTimer < max:
			vol := 1 - (float64(e.bgmFadingTimer) / max)
			audio.SetBGMVolume(vol)
		}
		if (input.Current().IsActionKeyJustPressed() || input.Current().IsSpaceJustTouched()) && e.timer > 5 {
			// 条件を満たしていると隠し画面へ
			if game.gameData.IsGetOmega() {
				if game.gameData.lunkerMode {
					e.gameStateMsg = GAMESTATE_MSG_REQ_SECRET_CLEAR
				} else {
					e.gameStateMsg = GAMESTATE_MSG_REQ_SECRET_COMMAND
				}
				return
			}
			e.gameStateMsg = GAMESTATE_MSG_REQ_TITLE
			audio.PauseBGM()
		}
	}
}

func (e *EndingScene) Draw(screen *ebiten.Image, game *Game) {
	if !game.transparent {
		draw.Draw(screen, "bg", 0, 0, 0, 480, 320, 240)
	}

	switch e.state {
	case ENDINGMAIN_STATE_STAFFROLL:
		for i, line := range strings.Split(text.Get(game.lang, text.TextIDEnding), "\n") {
			x := (draw.ScreenWidth - font.Width(line)) / 2
			y := draw.ScreenHeight - (e.timer / ENDING_SCROLL_SPEED) + i*font.LineHeight
			font.DrawText(screen, line, x, y, color.Black)
		}
	case ENDINGMAIN_STATE_RESULT:
		lines := []string{
			text.Get(game.lang, text.TextIDEndingScore1),
			"",
			text.Get(game.lang, text.TextIDEndingScore2),
			strconv.Itoa(game.gameData.GetItemCount()),
			"",
			text.Get(game.lang, text.TextIDEndingScore3),
			fmt.Sprintf("%.2f", float64(game.gameData.TimeInFrame())/60),
		}
		for i, line := range lines {
			x := (draw.ScreenWidth - font.Width(line)) / 2
			font.DrawText(screen, line, x, (draw.ScreenHeight-160)/2+16*i, color.Black)
		}
	}
}

func (e *EndingScene) Msg() GameStateMsg {
	return e.gameStateMsg
}

type SecretType int

const (
	SecretTypeCommand SecretType = iota
	SecretTypeClear
)

type SecretScene struct {
	gameStateMsg GameStateMsg
	timer        int
	secretType   SecretType
}

func NewSecretScene(secretType SecretType) *SecretScene {
	return &SecretScene{
		secretType: secretType,
	}
}

func (s *SecretScene) Update(game *Game) {
	s.timer++
	if (input.Current().IsActionKeyJustPressed() || input.Current().IsSpaceJustTouched()) && s.timer > 5 {
		s.gameStateMsg = GAMESTATE_MSG_REQ_TITLE
	}
}

func (s *SecretScene) Draw(screen *ebiten.Image, game *Game) {
	if !game.transparent {
		draw.Draw(screen, "bg", 0, 0, 0, 240, 320, 240)
	}
	var textID text.TextID
	switch s.secretType {
	case SecretTypeCommand:
		textID = text.TextIDSecretCommand
	case SecretTypeClear:
		textID = text.TextIDSecretClear
	default:
		panic("not reached")
	}
	str := text.Get(game.lang, textID)
	y := (draw.ScreenHeight - font.Height(str)) / 2
	for i, line := range strings.Split(str, "\n") {
		x := (draw.ScreenWidth - font.Width(line)) / 2
		font.DrawText(screen, line, x, y+i*font.LineHeight, color.White)
	}
}

func (s *SecretScene) Msg() GameStateMsg {
	return s.gameStateMsg
}

type GameScene struct {
	gameStateMsg GameStateMsg
	player       *Player
}

func NewGameScene(game *Game) *GameScene {
	g := &GameScene{
		player: NewPlayer(game.gameData),
	}
	return g
}

func (g *GameScene) Update(game *Game) {
	g.gameStateMsg = g.player.Update()
}

func (g *GameScene) Draw(screen *ebiten.Image, game *Game) {
	if !game.transparent {
		if game.gameData.lunkerMode {
			draw.Draw(screen, "bg", 0, 0, 0, 240, draw.ScreenWidth, draw.ScreenHeight)
		} else {
			draw.Draw(screen, "bg", 0, 0, 0, 0, draw.ScreenWidth, draw.ScreenHeight)
		}
	}
	g.player.Draw(screen, game)
	if input.Current().IsTouchEnabled() {
		draw.DrawTouchButtons(screen)
	}
}

func (g *GameScene) Msg() GameStateMsg {
	return g.gameStateMsg
}

type Scene interface {
	Update(g *Game) // TODO: Should return errors
	Draw(screen *ebiten.Image, g *Game)
	Msg() GameStateMsg
}
