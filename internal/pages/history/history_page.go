package history

import (
	"ttl-monitor/internal/managers"
)

type Page struct {
	managers.AbstractPage
}

func NewHistoryPage() *Page {
	page := &Page{}
	return page
}
