package ino

import (
	"fmt"
	"syscall/js"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/text/language"

	"go-inovation/ino/internal/text"

	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func tryLoseContext() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) && js.Global().Truthy() {
		doc := js.Global().Get("document")
		canvas := doc.Call("getElementsByTagName", "canvas").Index(0)
		context := canvas.Call("getContext", "webgl")
		context.Call("getExtension", "WEBGL_lose_context").Call("loseContext")
		fmt.Println("Context Lost!")
		return true
	}

	return false
}

func systemLang() language.Tag {
	nav := js.Global().Get("navigator")
	if nav.IsNull() {
		return language.Japanese
	}
	str := nav.Get("language").String()
	newLang, _ := language.Parse(str)
	base, _ := newLang.Base()
	newLang, _ = language.Compose(base)
	for _, l := range text.Languages() {
		if newLang == l {
			return newLang
		}
	}
	return language.Japanese
}
