// https://www.lihaoyi.com/post/BuildyourownCommandLinewithANSIescapecodes.html#colors
package display

// import (
//	"fmt"
//	"strings"
//	"sync"
//	"ttl-monitor/internal/common"
// )
//
// const (
// )
//
// var (
//	HEX = [16]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
// )
//
// func BinData(data uint8) string {
//	return fmt.Sprintf("%08b", data)
// }
// func HexData(data uint8) string {
//	return fmt.Sprintf("%s%s", HEX[data>>4], HEX[data&15])
// }
// func HexAddress(address uint16) string {
//	return fmt.Sprintf("%s%s%s%s", HEX[address>>12], HEX[address>>8&15], HEX[address>>4&15], HEX[address&15])
// }
//
// func (t *terminal) At(col int, row int) bool {
//	str := Bell
//	if col >= 1 && col <= t.cols && row >= 1 && row <= t.rows {
//		str = fmt.Sprintf(setPosition, row+t.y, col+t.x)
//		t.col = col
//		t.row = row
//	}
//	fmt.Printf(str)
//	return str != Bell
// }
// func (t *terminal) Start() {
//	fmt.Printf(//func (t *terminal) Start() {
// //	fmt.Printf(setColumn, t.x)
// //	t.col = t.x
// //}
// //func (t *terminal) Home() {
// //	fmt.Printf(setPosition, t.x, t.y)
// //	t.col = t.x
// //	t.row = t.y
// //}, t.x)
//	t.col = t.x
// }
// func (t *terminal) Home() {
//	fmt.Printf(setPosition, t.x, t.y)
//	t.col = t.x
//	t.row = t.y
// }
// func (t *terminal) PrintAtf(col int, row int, text string, a ...interface{}) bool {
//	return t.PrintAt(col, row, fmt.Sprintf(text, a...))
// }
// func (t *terminal) PrintAt(col int, row int, text string) bool {
//	ok := t.At(col, row)
//	if ok {
//		t.Print(text)
//	}
//	return ok
// }
// func (t *terminal) Printf(text string, a ...interface{}) {
//	t.Print(fmt.Sprintf(text, a...))
// }
// func (t *terminal) Print(text string) {
//	tokens, clean := t.Split(text)
//	max := t.cols - t.col + 1
//	if len(clean) > max {
//		text = ""
//		length := 0
//		for _, tkn := range tokens {
//			if length+tkn.length < max {
//				text += tkn.text
//				length += tkn.length
//			} else {
//				text += tkn.text[:max-length]
//				break
//			}
//		}
//	}
//	fmt.Printf("%s%s", text, common.Reset)
// }
// func (t *terminal) Bell() {
//	fmt.Printf(Bell)
// }
// func (t *terminal) Cls() {
//	if t.fullScreen {
//		fmt.Printf(clearScreen)
//	} else {
//		line := fmt.Sprintf("%s%s", common.BGBlue, strings.Repeat(" ", t.cols))
//		for row := 1; row <= t.rows; row++ {
//			t.PrintAt(1, row, line)
//		}
//	}
//	t.Home()
// }
// func (t *terminal) HideCursor() {
//	fmt.Printf(hide)
// }
// func (t *terminal) ShowCursor() {
//	fmt.Printf(show)
// }
// func (t *terminal) Row() int {
//	return t.row
// }
// func (t *terminal) Col() int {
//	return t.col
// }
// func (t *terminal) Rows() int {
//	return t.rows
// }
// func (t *terminal) Cols() int {
//	return t.cols
// }
// func (t *terminal) CreatePage() common.Page {
//	page := PageImpl{
//		display: t,
//		regions: []common.Region{},
//	}
//	Screen.lock.Lock()
//	defer t.lock.Unlock()
//	t.pages = append(t.pages, &page)
//	if len(t.activePage) == 0 {
//		t.activePage = append(t.activePage, &page)
//	}
//	return &page
// }
// func (t *terminal) IsActivePage() bool {
//	return t.ActivePage() == t.page
// }
// func (t *terminal) ActivePage() common.Page {
//	t.lock.Lock()
//	defer t.lock.Unlock()
//	if len(t.activePage) > 0 {
//		return t.activePage[len(t.activePage)-1]
//	}
//	return nil
// }
// func (t *terminal) ActivatePage(newPage common.Page) {
//	var unlockOnce sync.Once
//	t.lock.Lock()
//	defer unlockOnce.Do(t.lock.Unlock)
//	alreadyActive := false
//	for index, page := range t.activePage {
//		if page == newPage {
//			t.activePage = t.activePage[:index+1]
//			alreadyActive = true
//			break
//		}
//	}
//	if !alreadyActive {
//		t.activePage = append(t.activePage, newPage)
//	}
//	unlockOnce.Do(t.lock.Unlock)
//	newPage.Render()
// }
// func (t *terminal) DeactivatePage(oldPage common.Page) {
//	newActivePage := -1
//	t.lock.Lock()
//	t.lock.Unlock()
//	var newActivePages []common.Page
//
//	for index, page := range t.activePage {
//		if page != oldPage {
//			newActivePages = append(newActivePages, page)
//		} else if index == len(t.activePage)-1 {
//			newActivePage = index
//		}
//	}
//	t.activePage = newActivePages
//	if newActivePage == len(t.activePage) {
//		t.activePage[newActivePage].Render()
//	}
// }
// func (t *terminal) Activate() {
//	t.ActivatePage(t.page)
// }
// func (t *terminal) Render() {
//	if Screen.IsActivePage() {
//		t.lock.Lock()
//		defer t.lock.Unlock()
//		if t.renderer == nil {
//			t.renderer = defaultRenderer
//		}
//		t.renderer(t, t.cols, t.rows)
//	} else {
//		fmt.Sprintf("")
//	}
// }
// func (t *terminal) SetRenderer(renderer common.Renderer) {
//	t.renderer = renderer
// }
//
// type PageImpl struct {
//	display   common.Display
//	regions   []common.Region
//	processor common.ProcessInput
//	lock      sync.Mutex
// }
//
// func (p *PageImpl) CreateRegion(x, y, cols, rows int) common.Region {
//	if x < 1 {
//		x = 1
//	} else if x > p.display.Cols() {
//		x = p.display.Cols()
//	}
//	if y < 1 {
//		y = 1
//	} else if y > p.display.Rows() {
//		y = p.display.Rows()
//	}
//	if cols < 1 {
//		cols = 1
//	} else if cols+x-1 > p.display.Cols() {
//		cols = p.display.Cols() - x + 1
//	}
//	if rows < 1 {
//		rows = 1
//	} else if rows+y-1 > p.display.Rows() {
//		rows = p.display.Rows() - y + 1
//	}
//	region := &terminal{
//		x:          x - 1,
//		y:          y - 1,
//		col:        1,
//		row:        1,
//		cols:       cols,
//		rows:       rows,
//		fullScreen: x == 0 && y == 0 && cols == p.display.Cols() && rows == p.display.Rows(),
//		page:       p,
//	}
//	p.lock.Lock()
//	defer p.lock.Unlock()
//	p.regions = append(p.regions, region)
//	return region
// }
// func (p *PageImpl) CreatePageRegion() common.Region {
//	return p.CreateRegion(1, 1, p.display.Cols(), p.display.Rows())
// }
// func (p *PageImpl) Activate() {
//	p.display.Activate()
// }
// func (p *PageImpl) Deactivate() {
//	Screen.DeactivatePage(p)
// }
// func (p *PageImpl) ProcessesInput() bool {
//	return p.processor != nil
// }
// func (p *PageImpl) Render() {
//	if p.display.IsActivePage() {
//		p.lock.Lock()
//		defer p.lock.Unlock()
//		p.display.Cls()
//		for _, region := range p.regions {
//			region.Render()
//		}
//	}
// }
// func (p *PageImpl) SetInputProcessor(inputProcessor common.ProcessInput) {
//	p.processor = inputProcessor
// }
// func (p *PageImpl) ProcessInput(input common.Input) bool {
//	if p.processor != nil {
//		return p.processor(input)
//	}
//	return false
// }
//
// type token struct {
//	text   string
//	length int
// }
//
// func (t *terminal) Split(text string) ([]token, string) {
//	var result []token
//	clean := make([]byte, 0, len(text))
//	last := 0
//	for _, entry := range rex.FindAllStringIndex(text, -1) {
//		if last < entry[0] {
//			t := text[last:entry[0]]
//			result = append(result, token{text: t, length: len(t)})
//			clean = append(clean, t...)
//		}
//		t := text[entry[0]:entry[1]]
//		result = append(result, token{text: t, length: 0})
//		last = entry[1]
//	}
//	if last < len(text) {
//		t := text[last:len(text)]
//		result = append(result, token{text: t, length: len(t)})
//		clean = append(clean, t...)
//	}
//	return result, string(clean)
// }
