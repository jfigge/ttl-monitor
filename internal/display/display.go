package display

import (
	"fmt"
	"regexp"
	"ttl-monitor/internal/tty"
)

const (
	setPosition = "\u001b[%d;%dH" // moves cursor to row n column m
	clearScreen = "\u001b[2J"     // clears entire Screen
	setColumn   = "\u001b[%dG"    // moves cursor to column n
)

var (
	HEX = [16]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
	rex = regexp.MustCompile("\u001b\\[[0-9]{1,2}m")
	// lock sync.Mutex
)

type Display struct {
	col       int
	row       int
	offsetCol int
	offsetRow int
	cols      int
	rows      int
	cls       func()
}

func NewDisplay(offsetCol int, offsetRow int, cols int, rows int, full bool) *Display {
	d := &Display{
		offsetCol: offsetCol,
		offsetRow: offsetRow,
		cols:      cols,
		rows:      rows,
	}
	if full {
		d.cls = func() { fmt.Printf(clearScreen) }
	} else {
		d.cls = func() { ClearArea(cols, rows)(d) }
	}
	return d
}

func (d *Display) Cols() int {
	return d.cols
}

func (d *Display) Rows() int {
	return d.rows
}

func (d *Display) Bell() {
	fmt.Printf(tty.Bell)
}

func (d *Display) Start() {
	d.At(1, d.row)
}

func (d *Display) Home() {
	d.At(1, 1)
}

func (d *Display) Cls() {
	d.cls()
	d.Home()
}

func (d *Display) At(col int, row int) bool {
	str := tty.Bell
	if col >= 1 && col <= d.cols && row >= 1 && row <= d.rows {
		str = fmt.Sprintf(setPosition, row+d.offsetRow, col+d.offsetCol)
		d.col = col
		d.row = row
	}
	fmt.Printf(str)
	return str != tty.Bell
}

func (d *Display) PrintAtf(col int, row int, text string, a ...interface{}) bool {
	return d.PrintAt(col, row, fmt.Sprintf(text, a...))
}
func (d *Display) PrintAt(col int, row int, text string) bool {
	ok := d.At(col, row)
	if ok {
		d.Print(text)
	}
	return ok
}
func (d *Display) Printf(text string, a ...interface{}) {
	d.Print(fmt.Sprintf(text, a...))
}
func (d *Display) Print(text string) {
	tokens, clean := d.Split(text)
	max := d.cols - d.col + 1
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

func (d *Display) Split(text string) ([]token, string) {
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

func BinData(data uint8) string {
	return fmt.Sprintf("%08b", data)
}
func HexData(data uint8) string {
	return fmt.Sprintf("%s%s", HEX[data>>4], HEX[data&15])
}
func HexAddress(address uint16) string {
	return fmt.Sprintf("%s%s%s%s", HEX[address>>12], HEX[address>>8&15], HEX[address>>4&15], HEX[address&15])
}
