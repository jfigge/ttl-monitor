package managers

import (
	"fmt"
)

type PageFunc func(input *KeyInput) bool

type Page interface {
	Initialize()
	AddArea(area Area)
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

func (p *AbstractPage) AddArea(area Area) {
	if p.areas == nil {
		p.areas = make([]Area, 0, 10)
	}
	p.areas = append(p.areas, area)
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
	fmt.Printf("Key Code: %d, Ascii Code: %d\n", keyInput.KeyCode, keyInput.Ascii)
	return false
}
