package ino

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"go-inovation/ino/internal/audio"
	"go-inovation/ino/internal/draw"
	"go-inovation/ino/internal/field"
	"go-inovation/ino/internal/fieldtype"
	"go-inovation/ino/internal/input"
)

type PlayerState int

const (
	PLAYERSTATE_START PlayerState = iota
	PLAYERSTATE_NORMAL
	PLAYERSTATE_ITEMGET
	PLAYERSTATE_MUTEKI
	PLAYERSTATE_DEAD
)

const (
	PLAYER_SPEED         = 2.0
	PLAYER_GRD_ACCRATIO  = 0.04
	PLAYER_AIR_ACCRATIO  = 0.01
	PLAYER_JUMP          = -4.0
	PLAYER_GRAVITY       = 0.2
	PLAYER_FALL_SPEEDMAX = 4.0
	WAIT_TIMER_INTERVAL  = 10
	LIFE_RATIO           = 400
	MUTEKI_INTERVAL      = 50
	START_WAIT_INTERVAL  = 50
	SCROLLPANEL_SPEED    = 2.0

	LUNKER_JUMP_DAMAGE1 = 40.0
	LUNKER_JUMP_DAMAGE2 = 96.0
)

type Player struct {
	life        int
	jumpCnt     int
	timer       int
	position    PositionF
	speed       PositionF
	direction   int
	jumpedPoint PositionF
	state       PlayerState
	itemGet     fieldtype.FieldType
	waitTimer   int
	gameData    *GameData // TODO(hajimehoshi): Remove this?
	view        *View
	field       *field.Field
}

func NewPlayer(gameData *GameData) *Player {
	f := field.New(field_data)
	x, y := f.GetStartPoint()
	startPointF := PositionF{float64(x), float64(y)}
	audio.PlayBGM(audio.BGM0)
	return &Player{
		gameData:    gameData,
		field:       f,
		life:        gameData.lifeMax * LIFE_RATIO,
		position:    startPointF,
		jumpedPoint: startPointF,
		view:        NewView(startPointF),
	}
}

func (p *Player) onWall() bool {
	if p.toFieldOfsY() > field.CHAR_SIZE/4 {
		return false
	}
	if p.field.IsRidable(p.toFieldX(), p.toFieldY()+1) && p.toFieldOfsX() < field.CHAR_SIZE*7/8 {
		return true
	}
	if p.field.IsRidable(p.toFieldX()+1, p.toFieldY()+1) && p.toFieldOfsX() > field.CHAR_SIZE/8 {
		return true
	}
	return false
}

func (p *Player) isFallable() bool {
	if !p.onWall() {
		return false
	}
	if p.field.IsWall(p.toFieldX(), p.toFieldY()+1) && p.toFieldOfsX() < field.CHAR_SIZE*7/8 {
		return false
	}
	if p.field.IsWall(p.toFieldX()+1, p.toFieldY()+1) && p.toFieldOfsX() > field.CHAR_SIZE/8 {
		return false
	}
	return true
}

func (p *Player) isUpperWallBoth() bool {
	if p.toFieldOfsY() < field.CHAR_SIZE/2 {
		return false
	}
	if p.field.IsWall(p.toFieldX(), p.toFieldY()) && p.field.IsWall(p.toFieldX()+1, p.toFieldY()) {
		return true
	}
	return false
}

func (p *Player) isUpperWall() bool {
	if p.toFieldOfsY() < field.CHAR_SIZE/2 {
		return false
	}
	if p.field.IsWall(p.toFieldX(), p.toFieldY()) && p.toFieldOfsX() < field.CHAR_SIZE*7/8 {
		return true
	}
	if p.field.IsWall(p.toFieldX()+1, p.toFieldY()) && p.toFieldOfsX() > field.CHAR_SIZE/8 {
		return true
	}
	return false
}

func (p *Player) isLeftWall() bool {
	if p.field.IsWall(p.toFieldX(), p.toFieldY()) {
		return true
	}
	if p.field.IsWall(p.toFieldX(), p.toFieldY()+1) && p.toFieldOfsY() > field.CHAR_SIZE/8 {
		return true
	}
	return false
}

func (p *Player) isRightWall() bool {
	if p.field.IsWall(p.toFieldX()+1, p.toFieldY()) {
		return true
	}
	if p.field.IsWall(p.toFieldX()+1, p.toFieldY()+1) && p.toFieldOfsY() > field.CHAR_SIZE/8 {
		return true
	}
	return false
}

func (p *Player) normalizeToRight() {
	p.position.X = float64(p.toFieldX() * field.CHAR_SIZE)
	p.speed.X = 0
}

func (p *Player) normalizeToLeft() {
	p.position.X = float64((p.toFieldX() + 1) * field.CHAR_SIZE)
	p.speed.X = 0
}

func (p *Player) normalizeToUpper() {
	if p.speed.Y < 0 {
		p.speed.Y = 0
	}
	p.position.Y = float64(field.CHAR_SIZE * (p.toFieldY() + 1))
}

func (p *Player) toFieldX() int {
	return int(p.position.X) / field.CHAR_SIZE
}

func (p *Player) toFieldY() int {
	return int(p.position.Y) / field.CHAR_SIZE
}

func (p *Player) toFieldOfsX() int {
	return int(p.position.X) % field.CHAR_SIZE
}

func (p *Player) toFieldOfsY() int {
	return int(p.position.Y) % field.CHAR_SIZE
}

func (p *Player) Update() GameStateMsg {
	msg := GAMESTATE_MSG_NONE
	p.field.Update()
	switch p.state {
	case PLAYERSTATE_START:
		p.waitTimer++
		if p.waitTimer > START_WAIT_INTERVAL {
			p.state = PLAYERSTATE_NORMAL
		}

	case PLAYERSTATE_NORMAL:
		p.moveByInput()
		p.moveNormal()
		if p.life < p.gameData.lifeMax*LIFE_RATIO {
			o_life := p.life
			p.life++
			if (p.life / LIFE_RATIO) != (o_life / LIFE_RATIO) {
				audio.PlaySE(audio.SE_HEAL)
			}
		}

	case PLAYERSTATE_ITEMGET:
		p.moveItemGet()
		if p.state != PLAYERSTATE_ITEMGET {
			if p.gameData.IsGameClear() {
				msg = GAMESTATE_MSG_REQ_ENDING
			}
		}

	case PLAYERSTATE_MUTEKI:
		p.moveByInput()
		p.moveNormal()
		p.waitTimer++
		if p.waitTimer > MUTEKI_INTERVAL {
			p.state = PLAYERSTATE_NORMAL
		}

	case PLAYERSTATE_DEAD:
		p.moveNormal()
		audio.PauseBGM()
		if input.Current().IsActionKeyPressed() && p.waitTimer > 15 {
			msg = GAMESTATE_MSG_REQ_TITLE
		}
	}
	if p.life < LIFE_RATIO {
		if p.state != PLAYERSTATE_DEAD {
			p.waitTimer = 0
		}
		p.state = PLAYERSTATE_DEAD
		p.direction = 0
		p.waitTimer++
	}
	return msg
}

func (p *Player) moveNormal() {
	p.timer++
	p.gameData.Update()

	// 移動＆落下
	p.speed.Y += PLAYER_GRAVITY
	p.position.X += p.speed.X
	p.position.Y += p.speed.Y

	if p.speed.Y > PLAYER_FALL_SPEEDMAX {
		p.speed.Y = PLAYER_FALL_SPEEDMAX
	}

	if p.state == PLAYERSTATE_NORMAL {
		p.checkCollision()
	}

	// ATARI判定
	hitLeft := false
	hitRight := false
	hitUpper := false
	if p.onWall() && p.speed.Y >= 0 {
		if p.gameData.lunkerMode {
			if p.position.Y-p.jumpedPoint.Y > LUNKER_JUMP_DAMAGE1 {
				p.state = PLAYERSTATE_MUTEKI
				p.waitTimer = 0
				p.life -= LIFE_RATIO
				audio.PlaySE(audio.SE_DAMAGE)
			}
			if p.position.Y-p.jumpedPoint.Y > LUNKER_JUMP_DAMAGE2 {
				p.state = PLAYERSTATE_MUTEKI
				p.waitTimer = 0
				p.life -= LIFE_RATIO * 99
				audio.PlaySE(audio.SE_DAMAGE)
			}
		}

		if !input.Current().IsActionKeyPressed() || !input.Current().IsDirectionKeyPressed(input.DirectionDown) || !p.isFallable() {
			if p.speed.Y > 0 {
				p.speed.Y = 0
			}
			p.position.Y = float64(field.CHAR_SIZE * p.toFieldY())
			p.jumpCnt = 0
		}

		p.jumpedPoint = p.position
	}
	if p.isLeftWall() && p.speed.X < 0 {
		hitLeft = true
	}
	if p.isRightWall() && p.speed.X > 0 {
		hitRight = true
	}
	if p.isUpperWall() && p.speed.Y <= 0 {
		hitUpper = true
	}

	if hitUpper && !hitLeft && !hitRight {
		p.normalizeToUpper()
	}
	if !hitUpper && hitLeft {
		p.normalizeToLeft()
	}
	if !hitUpper && hitRight {
		p.normalizeToRight()
	}
	if hitUpper && hitRight {
		if p.isUpperWallBoth() {
			p.normalizeToUpper()
		} else {
			if p.toFieldOfsX() > p.toFieldOfsY() {
				p.normalizeToRight()
			} else {
				p.normalizeToUpper()
			}
		}
	}
	if hitUpper && hitLeft {
		if p.isUpperWallBoth() {
			p.normalizeToUpper()
		} else {
			if field.CHAR_SIZE-p.toFieldOfsX() > p.toFieldOfsY() {
				p.normalizeToLeft()
			} else {
				p.normalizeToUpper()
			}
		}
	}

	// 床特殊効果
	switch p.getOnField() {
	case fieldtype.FIELD_SCROLL_L:
		p.speed.X = p.speed.X*(1.0-PLAYER_GRD_ACCRATIO) + float64(p.direction*PLAYER_SPEED-SCROLLPANEL_SPEED)*PLAYER_GRD_ACCRATIO
	case fieldtype.FIELD_SCROLL_R:
		p.speed.X = p.speed.X*(1.0-PLAYER_GRD_ACCRATIO) + float64(p.direction*PLAYER_SPEED+SCROLLPANEL_SPEED)*PLAYER_GRD_ACCRATIO
	case fieldtype.FIELD_SLIP:
		// Do nothing
	case fieldtype.FIELD_NONE:
		p.speed.X = p.speed.X*(1.0-PLAYER_AIR_ACCRATIO) + float64(p.direction*PLAYER_SPEED)*PLAYER_AIR_ACCRATIO
	default:
		p.speed.X = p.speed.X*(1.0-PLAYER_GRD_ACCRATIO) + float64(p.direction*PLAYER_SPEED)*PLAYER_GRD_ACCRATIO
	}

	p.view.Update(p.position, p.speed)
}

func (p *Player) moveItemGet() {
	if p.waitTimer < WAIT_TIMER_INTERVAL {
		p.waitTimer++
		return
	}
	if input.Current().IsActionKeyJustPressed() {
		p.state = PLAYERSTATE_NORMAL
		audio.ResumeBGM(audio.BGM0)
	}
}

func (p *Player) moveByInput() {
	if input.Current().IsDirectionKeyPressed(input.DirectionLeft) {
		p.direction = -1
	}
	if input.Current().IsDirectionKeyPressed(input.DirectionRight) {
		p.direction = 1
	}

	if input.Current().IsActionKeyJustPressed() {
		if ((p.gameData.jumpMax > p.jumpCnt) || p.onWall()) && !input.Current().IsDirectionKeyPressed(input.DirectionDown) {
			p.speed.Y = PLAYER_JUMP // ジャンプ
			if !p.onWall() {
				p.jumpCnt++
			}

			if math.Abs(p.speed.X) < 0.1 {
				if p.speed.X < 0 {
					p.speed.X -= 0.02
				}
				if p.speed.X > 0 {
					p.speed.X += 0.02
				}
			}
			audio.PlaySE(audio.SE_JUMP)
			p.jumpedPoint = p.position
		}
	}
}

func (p *Player) checkCollision() {
	for xx := 0; xx < 2; xx++ {
		for yy := 0; yy < 2; yy++ {
			// アイテム獲得(STATE_ITEMGETへ遷移)
			if p.field.IsItem(p.toFieldX()+xx, p.toFieldY()+yy) {
				// 隠しアイテムは条件が必要
				if !p.field.IsItemGettable(p.toFieldX()+xx, p.toFieldY()+yy, p.gameData) {
					continue
				}

				p.state = PLAYERSTATE_ITEMGET

				// アイテム効果
				p.itemGet = p.field.GetField(p.toFieldX()+xx, p.toFieldY()+yy)
				switch p.field.GetField(p.toFieldX()+xx, p.toFieldY()+yy) {
				case fieldtype.FIELD_ITEM_POWERUP:
					p.gameData.jumpMax++
				case fieldtype.FIELD_ITEM_LIFE:
					p.gameData.lifeMax++
					p.life = p.gameData.lifeMax * LIFE_RATIO
				default:
					p.gameData.itemGetFlags[p.itemGet] = true
				}
				p.field.EraseField(p.toFieldX()+xx, p.toFieldY()+yy)
				p.waitTimer = 0

				audio.PauseBGM()
				if IsItemForClear(p.itemGet) || p.itemGet == fieldtype.FIELD_ITEM_POWERUP {
					audio.PlaySE(audio.SE_ITEMGET)
				} else {
					audio.PlaySE(audio.SE_ITEMGET2)
				}
				return
			}
			// トゲ(ダメージ)
			if p.field.IsSpike(p.toFieldX()+xx, p.toFieldY()+yy) {
				p.state = PLAYERSTATE_MUTEKI
				p.waitTimer = 0
				p.life -= LIFE_RATIO
				p.speed.Y = PLAYER_JUMP
				p.jumpCnt = -1 // ダメージ・エキストラジャンプ
				audio.PlaySE(audio.SE_DAMAGE)
				return
			}
		}
	}
}

func (p *Player) getOnField() fieldtype.FieldType {
	if !p.onWall() {
		return fieldtype.FIELD_NONE
	}
	x, y := p.toFieldX(), p.toFieldY()
	if p.toFieldOfsX() < field.CHAR_SIZE/2 {
		if p.field.IsRidable(x, y+1) {
			return p.field.GetField(x, y+1)
		}
		return p.field.GetField(x+1, y+1)
	}
	if p.field.IsRidable(x+1, y+1) {
		return p.field.GetField(x+1, y+1)
	}
	return p.field.GetField(x, y+1)
}

func (p *Player) drawPlayer(screen *ebiten.Image, game *Game) {
	v := p.view.ToScreenPosition(p.position)
	vx, vy := int(v.X), int(v.Y)
	if p.state == PLAYERSTATE_DEAD { // 死亡
		anime := (p.timer / 6) % 4
		if game.gameData.lunkerMode {
			draw.Draw(screen, "ino", vx, vy, field.CHAR_SIZE*(2+anime), 128+field.CHAR_SIZE*2, field.CHAR_SIZE, field.CHAR_SIZE)
			return
		}
		draw.Draw(screen, "ino", vx, vy, field.CHAR_SIZE*(2+anime), 128, field.CHAR_SIZE, field.CHAR_SIZE)
		return
	}
	if p.state != PLAYERSTATE_MUTEKI || p.timer%10 < 5 {
		anime := (p.timer / 6) % 2
		if !p.onWall() {
			anime = 0
		}
		if p.direction < 0 {
			if game.gameData.lunkerMode {
				draw.Draw(screen, "ino", vx, vy, field.CHAR_SIZE*anime, 128+field.CHAR_SIZE*2, field.CHAR_SIZE, field.CHAR_SIZE)
				return
			}
			draw.Draw(screen, "ino", vx, vy, field.CHAR_SIZE*anime, 128, field.CHAR_SIZE, field.CHAR_SIZE)
			return
		}
		if game.gameData.lunkerMode {
			draw.Draw(screen, "ino", vx, vy, field.CHAR_SIZE*anime, 128+field.CHAR_SIZE*3, field.CHAR_SIZE, field.CHAR_SIZE)
			return
		}
		draw.Draw(screen, "ino", vx, vy, field.CHAR_SIZE*anime, 128+field.CHAR_SIZE, field.CHAR_SIZE, field.CHAR_SIZE)
		return
	}
}

func (p *Player) drawLife(screen *ebiten.Image, game *Game) {
	for t := 0; t < game.gameData.lifeMax; t++ {
		if p.life < LIFE_RATIO*2 && p.timer%10 < 5 && game.gameData.lifeMax > 1 {
			continue
		}
		if p.life >= (t+1)*LIFE_RATIO {
			draw.Draw(screen, "ino",
				field.CHAR_SIZE*t, 0, field.CHAR_SIZE*3, 128+field.CHAR_SIZE*1, field.CHAR_SIZE, field.CHAR_SIZE)
			continue
		}
		draw.Draw(screen, "ino",
			field.CHAR_SIZE*t, 0, field.CHAR_SIZE*4, 128+field.CHAR_SIZE*1, field.CHAR_SIZE, field.CHAR_SIZE)
	}
}

func (p *Player) drawItems(screen *ebiten.Image, game *Game) {
	for t := fieldtype.FIELD_ITEM_FUJI; t < fieldtype.FIELD_ITEM_MAX; t++ {
		if !game.gameData.itemGetFlags[t] {
			draw.Draw(screen, "ino",
				draw.ScreenWidth-field.CHAR_SIZE/4*(int(fieldtype.FIELD_ITEM_MAX)-2-int(t)), 0, // 無
				field.CHAR_SIZE*5, 128+field.CHAR_SIZE, field.CHAR_SIZE/4, field.CHAR_SIZE/2)
			continue
		}
		// クリア条件アイテムは専用グラフィック
		if IsItemForClear(t) {
			for i, c := range clearFlagItems {
				if c != t {
					continue
				}
				draw.Draw(screen, "ino",
					draw.ScreenWidth-field.CHAR_SIZE/4*(int(fieldtype.FIELD_ITEM_MAX)-2-int(t)), 0,
					field.CHAR_SIZE*5+field.CHAR_SIZE/4*(i+2), 128+field.CHAR_SIZE, field.CHAR_SIZE/4, field.CHAR_SIZE/2)
			}
			continue
		}
		draw.Draw(screen, "ino",
			draw.ScreenWidth-field.CHAR_SIZE/4*(int(fieldtype.FIELD_ITEM_MAX)-2-int(t)), 0, // 有
			field.CHAR_SIZE*5+field.CHAR_SIZE/4, 128+field.CHAR_SIZE, field.CHAR_SIZE/4, field.CHAR_SIZE/2)
	}
}

func (p *Player) drawMessage(screen *ebiten.Image, game *Game) {
	switch p.state {
	case PLAYERSTATE_ITEMGET:
		t := WAIT_TIMER_INTERVAL - p.waitTimer
		draw.DrawItemMessage(screen, p.itemGet, (draw.ScreenHeight-96)/2+24-t*t, game.lang)
		draw.DrawItemFrame(screen, (draw.ScreenWidth-32)/2, (draw.ScreenHeight-96)/2-t*t-24)
		it := int(p.itemGet) - (int(fieldtype.FIELD_ITEM_BORDER) + 1)
		draw.Draw(screen, "ino", (draw.ScreenWidth-16)/2, (draw.ScreenHeight-96)/2-int(t)*int(t)-16,
			(it%16)*field.CHAR_SIZE, (it/16+4)*field.CHAR_SIZE, field.CHAR_SIZE, field.CHAR_SIZE)
	case PLAYERSTATE_START:
		key := "msg_" + game.lang.String()
		draw.Draw(screen, key, (draw.ScreenWidth-256)/2, 64+(draw.ScreenHeight-240)/2, 0, 96, 256, 32)
	case PLAYERSTATE_DEAD:
		key := "msg_" + game.lang.String()
		draw.Draw(screen, key, (draw.ScreenWidth-256)/2, 64+(draw.ScreenHeight-240)/2, 0, 128, 256, 32)
	}
}

func (p *Player) Draw(screen *ebiten.Image, game *Game) {
	po := p.view.GetPosition()
	p.field.Draw(screen, game.gameData, int(po.X), int(po.Y))
	p.drawPlayer(screen, game)
	p.drawLife(screen, game)
	p.drawItems(screen, game)
	p.drawMessage(screen, game)
}
