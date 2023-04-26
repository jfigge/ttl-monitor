package help

import (
	"ttl-monitor/internal/managers"
)

type Page struct {
	managers.AbstractPage
}

func NewHelpPage() *Page {
	page := &Page{}
	return page
}
