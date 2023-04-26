package front

import (
	"fmt"
	"ttl-monitor/internal/display"
	"ttl-monitor/internal/managers"
)

type Page struct {
	managers.AbstractPage
}

func NewFrontPage() *Page {
	page := &Page{}
	page.CreateArea(0, 0, 20, 10, false, display.BoxDrawer)
	page.CreateArea(20, 0, 20, 5, false, display.RandomDrawer)
	page.CreateArea(20, 5, 20, 5, false, display.RandomDrawer)
	return page
}

func (p *Page) ProcessInput(keyInput *managers.KeyInput) bool {
	if keyInput.Ascii == 'q' {
		return true
	}
	fmt.Printf("Key Code: %d, Ascii Code: %d\n", keyInput.KeyCode, keyInput.Ascii)
	return false
}
