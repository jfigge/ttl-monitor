package historypage

import (
	"fmt"
	"ttl-monitor/internal/managers"
)

type HistoryPage struct {
	managers.BasePage
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
