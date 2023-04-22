package helppage

import (
	"fmt"
	"ttl-monitor/internal/managers"
)

type HelpPage struct {
	managers.BasePage
}

func (p *HelpPage) PostDraw() {
	// TODO implement me
	panic("implement me")
}

func NewHeLpPage() *HelpPage {
	page := &HelpPage{}
	return page
}

func (p *HelpPage) Draw() {
	fmt.Println("Help")
}
