package pages

import (
	"ttl-monitor/internal/services"
	"ttl-monitor/internal/services/display"
)

type HelpPage struct {
	display.AbstractPage
}

func NewHelpPage(s *services.Serial) *HelpPage {
	page := &HelpPage{}
	return page
}
