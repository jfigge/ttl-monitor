package display

import (
	"context"
	"os"
	"sync"
	"time"

	"ttl-monitor/internal/common"
	"ttl-monitor/internal/models"

	"github.com/pkg/term"
	xterm "golang.org/x/term"
)

type PageId int

const (
	PageBusy    PageId = -2
	InvalidPage PageId = -1
)

var nextPage PageId = 1

type Terminal struct {
	name       string
	height     int
	width      int
	dirty      bool
	pages      map[PageId]*PageMeta
	activePage *PageMeta
	wg         *sync.WaitGroup
	pageStack  []PageId
	fd         int
	ioTerm     *term.Term
	pageOut    chan TermOutput
}

type ActionCode int

type TermOutput struct {
	action ActionCode
	id     PageId
}

const (
	Close ActionCode = iota
)

func NewTerminal(minWidth int, minHeight int) (*Terminal, error) {
	fd := int(os.Stdin.Fd())
	if !xterm.IsTerminal(fd) {
		return nil, common.Error("not in a terminal")
	}

	time.Sleep(100 * time.Millisecond)
	width, height, err := xterm.GetSize(fd)
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
		pages:     make(map[PageId]*PageMeta),
		pageStack: make([]PageId, 0, 10),
		wg:        &sync.WaitGroup{},
		fd:        fd,
		pageOut:   make(chan TermOutput),
	}

	go handlePageOutput(t.pageOut, t)

	return t, nil
}

func (t *Terminal) Height() int {
	return t.height
}

func (t *Terminal) Width() int {
	return t.width
}

func (t *Terminal) ActivePage() *PageMeta {
	return t.activePage
}

func (t *Terminal) AddPage(page Page) PageId {
	nextPage++
	pageMeta := &PageMeta{
		id:     nextPage,
		page:   page,
		in:     make(chan models.KeyInput, 10),
		out:    t.pageOut,
		hidden: true,
		canvas: NewCanvas(0, 0, t.width, t.height, true),
	}

	ctx := context.Background()
	ctx, pageMeta.cancel = context.WithCancel(ctx)
	go handleInput(ctx, pageMeta)

	t.pages[pageMeta.id] = pageMeta
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

	var ok bool
	t.activePage, ok = t.pages[id]
	if !ok {
		return InvalidPage
	}

	t.pageStack = append(t.pageStack, id)
	t.wg.Add(1)

	t.activePage.hidden = false
	if lastPageId != InvalidPage {
		t.pages[lastPageId].hidden = true
	}
	t.activePage.render(true)

	return id
}

func (t *Terminal) ClosePage() {
	pageCount := len(t.pageStack)
	if pageCount == 0 {
		return
	}

	pageMeta := t.pages[t.pageStack[pageCount-1]]
	pageMeta.terminate()
	t.activePage = nil
	t.pageStack = t.pageStack[:pageCount-1]
	t.wg.Done()
	pageCount--

	if pageCount == 0 {
		return
	}

	t.activePage = t.pages[t.pageStack[pageCount-1]]
	t.activePage.hidden = false
	t.activePage.render(true)
}

func (t *Terminal) Render(force bool) {
	if t.activePage == nil {
		return
	}
	t.activePage.render(force)
}

func (t *Terminal) Wait() {
	t.wg.Wait()
}

func (t *Terminal) Terminate() {
	t.activePage = nil
	for _, page := range t.pages {
		page.terminate()
	}
	t.pages = make(map[PageId]*PageMeta)
	l := len(t.pageStack)
	t.pageStack = make([]PageId, 0, 10)
	for i := 0; i < l; i++ {
		t.wg.Done()
	}
}

func handleInput(ctx context.Context, pageMeta *PageMeta) {
	for true {
		select {
		case keyEvent := <-pageMeta.in:
			if pageMeta.page.ProcessInput(&keyEvent) {
				pageMeta.out <- TermOutput{
					action: Close,
					id:     pageMeta.id,
				}
			}
		case <-ctx.Done():
			break
		}
	}
}

func handlePageOutput(out <-chan TermOutput, t *Terminal) {
	select {
	case cmd := <-out:
		switch cmd.action {
		case Close:
			t.ClosePage()
		}
	}
}
