package managers

import (
	"fmt"
	"ttl-monitor/internal/display"
)

type PageFunc func(input *KeyInput) bool

type Page interface {
	Initialize()
	CreateArea(x, y, w, h int, fullPage bool, renderer display.Renderer)
	Areas() []Area
	PreDraw()
	Draw()
	PostDraw()
	ProcessInput(input *KeyInput) bool
	Terminate()
}

type AbstractPage struct {
	areas []Area
}

func (p *AbstractPage) CreateArea(x, y, w, h int, full bool, renderer display.Renderer) {
	if p.areas == nil {
		p.areas = make([]Area, 0, 10)
	}
	p.areas = append(p.areas, NewAbstractArea(x, y, w, h, full, renderer))
}

func (p *AbstractPage) Areas() []Area {
	return p.areas
}

func (p *AbstractPage) Initialize() {
}

func (p *AbstractPage) PreDraw() {
}

func (p *AbstractPage) Draw() {
}

func (p *AbstractPage) PostDraw() {
}

func (p *AbstractPage) Terminate() {
	fmt.Printf("Terminating page")
}

func (p *AbstractPage) ProcessInput(keyInput *KeyInput) bool {
	return false
}
