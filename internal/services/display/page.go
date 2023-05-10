package display

import (
	"ttl-monitor/internal/models"
)

type PageFunc func(input *models.KeyInput) bool

type Page interface {
	Initialize()
	CreateRegion(x, y, w, h int, renderer Renderer)
	Regions() []Region
	PreDraw(canvas *Canvas)
	Draw(canvas *Canvas)
	PostDraw(canvas *Canvas)
	ProcessInput(input *models.KeyInput) bool
	Terminate()
	SetDirty(bool)
	IsDirty() bool
}

type AbstractPage struct {
	dirty   bool
	regions []Region
}

func (p *AbstractPage) CreateRegion(x, y, w, h int, renderer Renderer) {
	if p.regions == nil {
		p.regions = make([]Region, 0, 10)
	}
	p.regions = append(p.regions, NewAbstractRegion(x, y, w, h, false, renderer))
}

func (p *AbstractPage) Regions() []Region {
	return p.regions
}

func (p *AbstractPage) Initialize() {
}

func (p *AbstractPage) SetDirty(dirty bool) {
	p.dirty = dirty
}

func (p *AbstractPage) IsDirty() bool {
	return p.dirty
}

func (p *AbstractPage) PreDraw(canvas *Canvas) {
	if p.dirty {
		canvas.Cls()
	}
}

func (p *AbstractPage) Draw(canvas *Canvas) {
}

func (p *AbstractPage) PostDraw(canvas *Canvas) {
}

func (p *AbstractPage) Terminate() {
}

func (p *AbstractPage) ProcessInput(keyInput *models.KeyInput) bool {
	return false
}
