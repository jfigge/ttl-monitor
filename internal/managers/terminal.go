package managers

import (
	"context"
	"fmt"
	"github.com/pkg/term"
	xterm "golang.org/x/term"
	"os"
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

type KeyInput struct {
	Ascii   int
	KeyCode int
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
		id:      nextPage,
		page:    page,
		in:      make(chan KeyInput, 10),
		out:     t.pageOut,
		hidden:  true,
		display: display.NewDisplay(0, 0, t.width, t.height, true),
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

func (t *Terminal) ReadChar() (*KeyInput, error) {
	t.ioTerm, _ = term.Open("/dev/tty")
	if t.ioTerm == nil {
		return nil, fmt.Errorf("terminal unavailable")
	} else {
		defer func() {
			if err := t.ioTerm.Restore(); err != nil {
				common.Debugf("Failed to restore terminal mode: %v", err)
			}
			if err := t.ioTerm.Close(); err != nil {
				common.Debugf("Failed to close terminal input: %v", err)
			}
		}()
	}

	if err := term.RawMode(t.ioTerm); err != nil {
		return nil, common.Errorf("Failed to access terminal RawMode: %v", err)
	}
	bs := make([]byte, 5)

	if err := t.ioTerm.SetReadTimeout(2 * time.Second); err != nil {
		common.Warn("Failed to set read timeout")
	}

	var numRead int
	var err error
	keyInput := KeyInput{}
	if numRead, err = t.ioTerm.Read(bs); err != nil {
		if err.Error() != "EOF" {
			if err = t.ioTerm.Restore(); err != nil {
				err = common.Errorf("Failed to restore terminal mode: %v", err)
			}
			if err = t.ioTerm.Close(); err != nil {
				err = common.Errorf("Failed to close terminal input: %v", err)
			}
		}
		return nil, err
	} else if numRead == 3 && bs[0] == 27 && bs[1] == 91 {
		// Three-character control sequence, beginning with "ESC-[".

		// Since there are no ASCII lines for arrow keys, we use
		// Javascript key lines.
		switch bs[2] {
		case 65:
			keyInput.KeyCode = 38 // Up
		case 66:
			keyInput.KeyCode = 40 // Down
		case 67:
			keyInput.KeyCode = 39 // Right
		case 68:
			keyInput.KeyCode = 37 // Left
		}
	} else if numRead == 3 && bs[0] == 0x1B && bs[1] == 0x4F {
		switch bs[2] {
		case 50:
			keyInput.KeyCode = 101 // Option+1
		case 51:
			keyInput.KeyCode = 102 // Option+2
		case 52:
			keyInput.KeyCode = 102 // Option+3
		case 53:
			keyInput.KeyCode = 103 // Option+4
		}
	} else if numRead == 1 {
		keyInput.Ascii = int(bs[0])
	} else {
		text := fmt.Sprintf("%d characters read.", numRead)
		for i := 0; i < numRead; i++ {
			text = fmt.Sprintf("%s %d:%s", text, i, display.HexData(bs[i]))
		}
		common.Warnf(text)
		// Two characters read??
	}
	return &keyInput, nil
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
	display.NewDisplay(0, 0, t.width, t.height, true).Cls()
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
