package display

type Region interface {
	IsDirty() bool
	SetDirty(dirty bool)
	Draw()
}

type AbstractRegion struct {
	width    int
	height   int
	renderer Renderer
	display  *Canvas
	dirty    bool
}

func NewAbstractRegion(x, y, w, h int, full bool, renderer Renderer) *AbstractRegion {
	return &AbstractRegion{
		width:    w,
		height:   h,
		renderer: renderer,
		dirty:    true,
		display:  NewCanvas(x, y, w, h, full),
	}
}

func (a *AbstractRegion) IsDirty() bool {
	return a.dirty
}

func (a *AbstractRegion) getDisplay() *Canvas {
	return a.display
}

func (a *AbstractRegion) SetDirty(dirty bool) {
	a.dirty = dirty
}

func (a *AbstractRegion) Draw() {
	if !a.dirty {
		return
	}
	a.renderer(a.display)
	a.dirty = false
}
