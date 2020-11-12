package ino

import (
	"go-inovation/ino/internal/draw"
	"go-inovation/ino/internal/input"
)

type PositionF struct {
	X float64
	Y float64
}

type View struct {
	position PositionF
}

func NewView(position PositionF) *View {
	return &View{
		position: position,
	}
}

func (v *View) ToScreenPosition(p PositionF) PositionF {
	x := p.X - v.GetPosition().X + draw.ScreenWidth/2
	y := p.Y - v.GetPosition().Y + draw.ScreenHeight/2
	return PositionF{x, y}
}

func (v *View) GetPosition() PositionF {
	p := v.position
	if input.Current().IsTouchEnabled() {
		p.Y += 16
	}
	return p
}

func (v *View) Update(position, speed PositionF) {
	const VIEW_DIRECTION_OFFSET = 30.0
	v.position.X = v.position.X*0.95 + (position.X+speed.X*VIEW_DIRECTION_OFFSET)*0.05
	v.position.Y = v.position.Y*0.95 + position.Y*0.05
}
