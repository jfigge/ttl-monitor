package managers

import (
	"golang.org/x/term"
	"sync"
	"time"
	"ttl-monitor/internal/common"
	"ttl-monitor/internal/display"
)

type PageId int

const (
	PageBusy    PageId = -2
	InvalidPage PageId = -1
)

var nextPage PageId = 1

type Terminal struct {
	name      string
	height    int
	width     int
	dirty     bool
	pages     map[PageId]Page
	wg        *sync.WaitGroup
	pageStack []PageId
}

func NewTerminal(minWidth int, minHeight int) (*Terminal, error) {
	if !term.IsTerminal(0) {
		return nil, common.Error("not in a terminal")
	}

	time.Sleep(100 * time.Millisecond)
	width, height, err := term.GetSize(0)
	if err != nil {
		return nil, common.Errorf("Failed to retrieve display size: %v\n", err)
	}

	if height < minHeight {
		return nil, common.Errorf("Minimal terminal height is %d.  Found %d\n", minHeight, height)
	} else if width < minWidth {
		return nil, common.Errorf("Minimal terminal width is %d.  Found %d\n", minWidth, width)
	}

	t := &Terminal{
		height:    height,
		width:     width,
		pages:     make(map[PageId]Page),
		pageStack: make([]PageId, 0, 10),
		wg:        &sync.WaitGroup{},
	}

	return t, nil
}

func (t *Terminal) Height() int {
	return t.height
}

func (t *Terminal) Width() int {
	return t.width
}

func (t *Terminal) AddPage(page Page) PageId {
	nextPage++
	t.pages[nextPage] = page
	page.SetHidden(true)
	page.SetDisplay(display.NewDisplay(0, 0, t.width, t.height, true))
	return nextPage
}

func (t *Terminal) ShowPage(id PageId) PageId {
	var lastPageId = InvalidPage
	if len(t.pageStack) > 0 {
		lastPageId = t.pageStack[len(t.pageStack)-1]
		if lastPageId == id {
			return id
		}
	}
	for _, pageId := range t.pageStack {
		if pageId == id {
			return PageBusy
		}
	}

	page, ok := t.pages[id]
	if !ok {
		return InvalidPage
	}

	t.pageStack = append(t.pageStack, id)
	t.wg.Add(1)

	page.SetHidden(false)
	if lastPageId != InvalidPage {
		t.pages[lastPageId].SetHidden(true)
	}
	page.Render(true)

	return id
}

func (t *Terminal) ClosePage() {
	pageCount := len(t.pageStack)
	if pageCount == 0 {
		return
	}

	page := t.pages[t.pageStack[pageCount-1]]
	page.SetHidden(true)
	t.pageStack = t.pageStack[:pageCount-1]
	t.wg.Done()
	pageCount--

	if pageCount == 0 {
		return
	}

	page = t.pages[t.pageStack[pageCount-1]]
	page.SetHidden(false)
	page.Render(true)
}

func (t *Terminal) Render(force bool) {
	if len(t.pageStack) == 0 {
		return
	}
	t.pages[t.pageStack[len(t.pageStack)-1]].Render(force)
}

func (t *Terminal) Wait() {
	t.wg.Wait()
}
