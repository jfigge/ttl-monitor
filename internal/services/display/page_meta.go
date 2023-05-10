package display

import (
	"context"

	"ttl-monitor/internal/models"
)

type PageMeta struct {
	id     PageId
	page   Page
	in     chan models.KeyInput
	out    chan<- TermOutput
	hidden bool
	canvas *Canvas
	cancel context.CancelFunc
}

func (m *PageMeta) isHidden() bool {
	return m.hidden
}
func (m *PageMeta) setHidden(hidden bool) {
	if m.hidden != hidden {
		m.hidden = hidden
		if !hidden {
			m.page.SetDirty(true)
		}
	}
}

func (m *PageMeta) SetDirty(dirty bool) {
	m.page.SetDirty(dirty)
	for _, area := range m.page.Regions() {
		area.SetDirty(dirty)
	}
}

func (m *PageMeta) needsRendering() bool {
	if m.page.IsDirty() {
		return true
	}
	for _, area := range m.page.Regions() {
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
		m.page.SetDirty(true)
	}
	if !m.needsRendering() { // && !p.hidden {
		return
	}
	m.page.PreDraw(m.canvas)
	m.page.Draw(m.canvas)
	for _, area := range m.page.Regions() {
		area.Draw()
	}
	m.page.SetDirty(false)
	m.page.PostDraw(m.canvas)
}

func (m *PageMeta) terminate() {
	m.page.Terminate()
	m.hidden = true
	m.cancel()
}

func (m *PageMeta) ProcessInput(keys *models.KeyInput) {
	m.in <- *keys
}
