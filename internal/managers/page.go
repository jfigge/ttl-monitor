package managers

import "ttl-monitor/internal/display"

type AreaId int

const (
	AreaBusy    AreaId = -2
	InvalidArea AreaId = -1
)

var nextArea AreaId = 1

type BasePage struct {
	id      PageId
	areas   map[AreaId]*Area
	dirty   bool
	hidden  bool
	display *display.Display
}

type Page interface {
	// Id() PageId
	SetDisplay(d *display.Display)
	NeedsRendering() bool
	SetDirty(dirty bool)
	IsHidden() bool
	SetHidden(dirty bool)
	Render(force bool)

	PreDraw()
	Draw()
	PostDraw()
	AddArea(area *Area) AreaId
}

func (p *BasePage) SetDisplay(d *display.Display) {
	p.display = d
}

func (p *BasePage) GetDisplay() *display.Display {
	return p.display
}

func (p *BasePage) NeedsRendering() bool {
	if p.dirty {
		return p.dirty
	}
	for _, area := range p.areas {
		if area.IsDirty() {
			return true
		}
	}
	return false
}

func (p *BasePage) SetDirty(dirty bool) {
	p.dirty = dirty
	for _, area := range p.areas {
		area.SetDirty(dirty)
	}
}

func (p *BasePage) IsHidden() bool {
	return p.hidden
}
func (p *BasePage) SetHidden(hidden bool) {
	p.hidden = hidden
}
func (p *BasePage) Id() PageId {
	return p.id
}

func (p *BasePage) AddArea(area *Area) AreaId {
	if p.areas == nil {
		p.areas = make(map[AreaId]*Area)
	}
	nextArea++
	p.areas[nextArea] = area
	return nextArea
}

func (p *BasePage) Render(force bool) {
	if force {
		p.SetDirty(true)
	}
	if !p.NeedsRendering() || p.hidden {
		return
	}
	p.PreDraw()
	p.Draw()
	p.PostDraw()
}

func (p *BasePage) PreDraw() {
	if p.dirty {
		p.display.Cls()
	}
}

func (p *BasePage) Draw() {
	for _, area := range p.areas {
		area.Draw()
	}
}

func (p *BasePage) PostDraw() {
	p.dirty = false
}
