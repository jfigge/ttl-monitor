package historypage

import (
	"fmt"
	"ttl-monitor/internal/managers"
)

type HistoryPage struct {
	managers.AbstractPage
}

func (p *HistoryPage) KeyEvent(keyInput *managers.KeyInput) bool {
	// TODO implement me
	panic("implement me")
}

func (p *HistoryPage) PostDraw() {
	// TODO implement me
	panic("implement me")
}

func NewHistoryPage() *HistoryPage {
	page := &HistoryPage{}
	return page
}

func (p *HistoryPage) Draw() {
	fmt.Println("History")
}
