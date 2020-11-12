package field

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"go-inovation/ino/internal/draw"
	"go-inovation/ino/internal/fieldtype"
)

const (
	CHAR_SIZE = 16
	maxFieldX = 128
	maxFieldY = 128
)

type Field struct {
	field [maxFieldX * maxFieldY]fieldtype.FieldType
	timer int
}

func New(data string) *Field {
	f := &Field{}
	xm := strings.Split(data, "\n")
	const decoder = " HUB~<>*I PabcdefghijklmnopqrzL@"

	for yy, line := range xm {
		for xx, c := range line {
			n := strings.IndexByte(decoder, byte(c))
			f.field[yy*maxFieldX+xx] = fieldtype.FieldType(n)
		}
	}
	return f
}

func (f *Field) Update() {
	f.timer++
}

func (f *Field) GetStartPoint() (int, int) {
	for yy := 0; yy < maxFieldY; yy++ {
		for xx := 0; xx < maxFieldX; xx++ {
			if f.GetField(xx, yy) == fieldtype.FIELD_ITEM_STARTPOINT {
				x := xx * CHAR_SIZE
				y := yy * CHAR_SIZE
				f.EraseField(xx, yy)
				return x, y
			}
		}
	}
	panic("no start point")
}

func (f *Field) IsWall(x, y int) bool {
	return f.field[y*maxFieldX+x] != fieldtype.FIELD_NONE &&
		f.field[y*maxFieldX+x] != fieldtype.FIELD_HIDEPATH &&
		f.field[y*maxFieldX+x] != fieldtype.FIELD_BAR &&
		!f.IsItem(x, y)
}
func (f *Field) IsRidable(x, y int) bool {
	return f.field[y*maxFieldX+x] != fieldtype.FIELD_NONE &&
		f.field[y*maxFieldX+x] != fieldtype.FIELD_HIDEPATH &&
		!f.IsItem(x, y)
}

func (f *Field) IsSpike(x, y int) bool {
	return f.field[y*maxFieldX+x] == fieldtype.FIELD_SPIKE
}

func (f *Field) GetField(x, y int) fieldtype.FieldType {
	return f.field[y*maxFieldX+x]
}

func (f *Field) IsItem(x, y int) bool {
	return f.field[y*maxFieldX+x] >= fieldtype.FIELD_ITEM_BORDER &&
		f.field[y*maxFieldX+x] != fieldtype.FIELD_ITEM_STARTPOINT
}

func (f *Field) IsItemGettable(x, y int, gameData GameData) bool {
	if !f.IsItem(x, y) {
		return false
	}
	if f.field[y*maxFieldX+x] == fieldtype.FIELD_ITEM_OMEGA && gameData.IsHiddenSecret() {
		return false
	}
	return true
}

func (f *Field) EraseField(x, y int) {
	f.field[y*maxFieldX+x] = fieldtype.FIELD_NONE
}

type GameData interface {
	IsHiddenSecret() bool
}

func (f *Field) Draw(screen *ebiten.Image, gameData GameData, viewPositionX, viewPositionY int) {
	const (
		graphicOffsetX = -16 - 16*2
		graphicOffsetY = 8 - 16*2
	)
	vx, vy := viewPositionX, viewPositionY
	ofs_x := CHAR_SIZE - vx%CHAR_SIZE
	ofs_y := CHAR_SIZE - vy%CHAR_SIZE
	for xx := -(draw.ScreenWidth/CHAR_SIZE/2 + 2); xx < (draw.ScreenWidth/CHAR_SIZE/2 + 2); xx++ {
		fx := xx + vx/CHAR_SIZE
		if fx < 0 || fx >= maxFieldX {
			continue
		}
		for yy := -(draw.ScreenHeight/CHAR_SIZE/2 + 2); yy < (draw.ScreenHeight/CHAR_SIZE/2 + 2); yy++ {
			fy := yy + vy/CHAR_SIZE
			if fy < 0 || fy >= maxFieldY {
				continue
			}

			gy := (f.timer / 10) % 4
			gx := int(f.field[fy*maxFieldX+fx])

			if f.IsItem(fx, fy) {
				gx -= (int(fieldtype.FIELD_ITEM_BORDER) + 1)
				gy = 4 + gx/16
				gx = gx % 16
			}

			if gameData.IsHiddenSecret() && f.field[fy*maxFieldX+fx] == fieldtype.FIELD_ITEM_OMEGA {
				continue
			}

			draw.Draw(screen, "ino",
				(xx+12)*CHAR_SIZE+ofs_x+graphicOffsetX+(draw.ScreenWidth-320)/2,
				(yy+8)*CHAR_SIZE+ofs_y+graphicOffsetY+(draw.ScreenHeight-240)/2,
				gx*16, gy*16, 16, 16)
		}
	}
}
