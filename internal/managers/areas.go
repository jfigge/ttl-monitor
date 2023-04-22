package managers

import (
	"ttl-monitor/internal/display"
)

type Area struct {
	id       AreaId
	width    int
	height   int
	renderer display.Renderer
	display  *display.Display
	dirty    bool
}

func NewArea(x, y, w, h int, full bool, renderer display.Renderer) *Area {
	return &Area{
		width:    w,
		height:   h,
		renderer: renderer,
		display:  display.NewDisplay(x, y, w, h, full),
	}
}

func (a *Area) IsDirty() bool {
	return a.dirty
}

func (a *Area) SetDirty(dirty bool) {
	a.dirty = dirty
}

func (p *Area) GetDisplay() *display.Display {
	return p.display
}

func (a *Area) Draw() {
	if !a.dirty {
		return
	}
	a.renderer(a.display)
	a.dirty = false
}
