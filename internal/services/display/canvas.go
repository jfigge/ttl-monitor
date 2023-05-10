package display

import (
	"fmt"
	"regexp"
	"strings"

	"ttl-monitor/internal/tty"
)

const (
	setPosition = "\u001b[%d;%dH" // moves cursor to row n column m
	clearScreen = "\u001b[2J"     // clears entire Screen
	setColumn   = "\u001b[%dG"    // moves cursor to column n
)

var (
	rex = regexp.MustCompile("\u001b\\[[0-9]{1,2}m")
)

type Canvas struct {
	col       int
	row       int
	offsetCol int
	offsetRow int
	cols      int
	rows      int
	cls       func()
}

func NewCanvas(offsetCol int, offsetRow int, cols int, rows int, full bool) *Canvas {
	d := &Canvas{
		offsetCol: offsetCol,
		offsetRow: offsetRow,
		cols:      cols,
		rows:      rows,
	}
	if full {
		d.cls = func() { fmt.Printf(clearScreen) }
	} else {
		d.cls = func() { ClearArea(d) }
	}
	return d
}

func ClearArea(d *Canvas) {
	w := d.Cols()
	h := d.Rows()
	blank := strings.Repeat(" ", w)
	for i := 0; i < h; i++ {
		d.PrintAtf(1, i+1, blank)
	}
}

func (c *Canvas) Cols() int {
	return c.cols
}

func (c *Canvas) Rows() int {
	return c.rows
}

func (c *Canvas) Bell() {
	fmt.Printf(tty.Bell)
}

func (c *Canvas) Start() {
	c.At(1, c.row)
}

func (c *Canvas) Home() {
	c.At(1, 1)
}

func (c *Canvas) Cls() {
	c.cls()
	c.Home()
}

func (c *Canvas) At(col int, row int) bool {
	str := tty.Bell
	if col >= 1 && col <= c.cols && row >= 1 && row <= c.rows {
		str = fmt.Sprintf(setPosition, row+c.offsetRow, col+c.offsetCol)
		c.col = col
		c.row = row
	}
	fmt.Printf(str)
	return str != tty.Bell
}

func (c *Canvas) PrintAtf(col int, row int, text string, a ...interface{}) bool {
	return c.PrintAt(col, row, fmt.Sprintf(text, a...))
}
func (c *Canvas) PrintAt(col int, row int, text string) bool {
	ok := c.At(col, row)
	if ok {
		c.Print(text)
	}
	return ok
}
func (c *Canvas) Printf(text string, a ...interface{}) {
	c.Print(fmt.Sprintf(text, a...))
}
func (c *Canvas) Print(text string) {
	tokens, clean := c.Split(text)
	max := c.cols - c.col + 1
	if len(clean) > max {
		text = ""
		length := 0
		for _, tkn := range tokens {
			if length+tkn.length < max {
				text += tkn.text
				length += tkn.length
			} else {
				text += tkn.text[:max-length]
				break
			}
		}
	}
	fmt.Printf("%s%s", text, tty.Reset)
}

type token struct {
	text   string
	length int
}

func (c *Canvas) Split(text string) ([]token, string) {
	var result []token
	clean := make([]byte, 0, len(text))
	last := 0
	for _, entry := range rex.FindAllStringIndex(text, -1) {
		if last < entry[0] {
			t := text[last:entry[0]]
			result = append(result, token{text: t, length: len(t)})
			clean = append(clean, t...)
		}
		t := text[entry[0]:entry[1]]
		result = append(result, token{text: t, length: 0})
		last = entry[1]
	}
	if last < len(text) {
		t := text[last:]
		result = append(result, token{text: t, length: len(t)})
		clean = append(clean, t...)
	}
	return result, string(clean)
}
