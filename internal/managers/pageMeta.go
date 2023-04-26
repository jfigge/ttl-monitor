package managers

import (
	"context"
	"ttl-monitor/internal/display"
)

type PageMeta struct {
	id      PageId
	page    Page
	in      chan KeyInput
	out     chan<- TermOutput
	hidden  bool
	dirty   bool
	display *display.Display
	cancel  context.CancelFunc
}

func (m *PageMeta) isHidden() bool {
	return m.hidden
}
func (m *PageMeta) setHidden(hidden bool) {
	m.hidden = hidden
}

func (m *PageMeta) setDirty(dirty bool) {
	m.dirty = dirty
	for _, area := range m.page.Areas() {
		area.SetDirty(dirty)
	}
}

func (m *PageMeta) needsRendering() bool {
	if m.dirty {
		return m.dirty
	}
	for _, area := range m.page.Areas() {
		if area.IsDirty() {
			return true
		}
	}
	return false
}

func (m *PageMeta) render(force bool) {
	if m.hidden {
		return
	}
	if force {
		m.setDirty(true)
	}
	if !m.needsRendering() { // && !p.hidden {
		return
	}
	m.page.PreDraw()
	if m.dirty {
		m.display.Cls()
	}
	m.page.Draw()
	for _, area := range m.page.Areas() {
		area.Draw()
	}
	m.dirty = false
	m.page.PostDraw()
}

func (m *PageMeta) terminate() {
	m.page.Terminate()
	m.hidden = true
	m.cancel()
}

func (m *PageMeta) ProcessInput(keys *KeyInput) {
	m.in <- *keys
}
