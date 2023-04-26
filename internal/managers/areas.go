package managers

import (
	"ttl-monitor/internal/display"
)

type Area interface {
	IsDirty() bool
	SetDirty(dirty bool)
	GetDisplay() *display.Display
	Draw()
}

type AbstractArea struct {
	width    int
	height   int
	renderer display.Renderer
	display  *display.Display
	dirty    bool
}

func NewAbstractArea(x, y, w, h int, full bool, renderer display.Renderer) *AbstractArea {
	return &AbstractArea{
		width:    w,
		height:   h,
		renderer: renderer,
		display:  display.NewDisplay(x, y, w, h, full),
	}
}

func (a *AbstractArea) IsDirty() bool {
	return a.dirty
}

func (a *AbstractArea) SetDirty(dirty bool) {
	a.dirty = dirty
}

func (a *AbstractArea) GetDisplay() *display.Display {
	return a.display
}

func (a *AbstractArea) Draw() {
	if !a.dirty {
		return
	}
	a.renderer(a.display)
	a.dirty = false
}
