package pages

import (
	"ttl-monitor/internal/services"
	"ttl-monitor/internal/services/display"
)

type HistoryPage struct {
	display.AbstractPage
}

func NewHistoryPage(s *services.Serial) *HistoryPage {
	page := &HistoryPage{}
	return page
}
