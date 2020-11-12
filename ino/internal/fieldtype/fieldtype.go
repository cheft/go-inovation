package fieldtype

import (
	"golang.org/x/text/language"

	"go-inovation/ino/internal/text"
)

type FieldType int

const (
	FIELD_NONE         FieldType = iota // なし
	FIELD_HIDEPATH                      // 隠しルート(見えるけど判定のないブロック)
	FIELD_UNVISIBLE                     // 不可視ブロック(見えないけど判定があるブロック)
	FIELD_BLOCK                         // 通常ブロック
	FIELD_BAR                           // 床。降りたり上ったりできる
	FIELD_SCROLL_L                      // ベルト床左
	FIELD_SCROLL_R                      // ベルト床右
	FIELD_SPIKE                         // トゲ
	FIELD_SLIP                          // すべる
	FIELD_ITEM_BORDER                   // アイテムチェック用
	FIELD_ITEM_POWERUP                  // パワーアップ
	// ふじ系
	FIELD_ITEM_FUJI
	FIELD_ITEM_BUSHI
	FIELD_ITEM_APPLE
	FIELD_ITEM_V
	// たか系
	FIELD_ITEM_TAKA
	FIELD_ITEM_SHUOLDER
	FIELD_ITEM_DAGGER
	FIELD_ITEM_KATAKATA
	// なす系
	FIELD_ITEM_NASU
	FIELD_ITEM_BONUS
	FIELD_ITEM_NURSE
	FIELD_ITEM_NAZUNA
	// くそげー系
	FIELD_ITEM_GAMEHELL
	FIELD_ITEM_GUNDAM
	FIELD_ITEM_POED
	FIELD_ITEM_MILESTONE
	FIELD_ITEM_1YEN
	FIELD_ITEM_TRIANGLE
	FIELD_ITEM_OMEGA      // 隠し
	FIELD_ITEM_LIFE       // ハート
	FIELD_ITEM_STARTPOINT // 開始地点
	FIELD_ITEM_MAX
)

func (f FieldType) ItemMessage(lang language.Tag) string {
	return text.Get(lang, text.TextID(f-FIELD_ITEM_POWERUP)+text.TextIDItemPowerUp)
}
