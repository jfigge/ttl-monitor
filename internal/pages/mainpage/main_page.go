package mainpage

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"ttl-monitor/internal/display"
	"ttl-monitor/internal/managers"
)

type MainPage struct {
	Y int
	x int
	managers.BasePage
	left       *managers.Area
	upperRight *managers.Area
	lowerRight *managers.Area
}

func NewMainPage() *MainPage {
	page := &MainPage{
		left:       managers.NewArea(0, 0, 20, 10, false, display.BoxDrawer(20, 10)),
		upperRight: managers.NewArea(20, 0, 20, 5, false, RandomDrawer(20, 5)),
		lowerRight: managers.NewArea(20, 5, 20, 5, false, RandomDrawer(20, 5)),
	}
	page.AddArea(page.left)
	page.AddArea(page.upperRight)
	page.AddArea(page.lowerRight)

	go page.redrawUpperRight(page.upperRight)
	go page.bound(page.left)
	return page
}

func RandomDrawer(w int, h int) func(d *display.Display) {
	return func(d *display.Display) {
		for i := 0; i < h; i++ {
			title := fmt.Sprintf(" %d x %d ", w, h)
			if i == 0 {
				left := w/2 - (len(title) / 2)
				right := w - len(title) - left
				if left > 1 || right > 1 {
					d.PrintAtf(1, 1, "%s%s%s", strings.Repeat("#", left), title, strings.Repeat("#", right))
				} else {
					d.PrintAtf(1, 1, "%s", strings.Repeat("#", w))
				}
			} else if i == h-1 {
				d.PrintAtf(1, i+1, "%s", strings.Repeat("#", w))
			} else {
				d.PrintAtf(1, i+1, "#%s#", strings.Repeat(" ", w-2))
			}
		}
		s := strings.Repeat("#", rand.Intn(w-2))
		d.PrintAtf(w/2-len(s)/2, h/2+1, "%s", s)
	}
}
func (p *MainPage) redrawUpperRight(area *managers.Area) {
	t := time.NewTicker(3 * time.Second)
	for true {
		select {
		case <-t.C:
			area.SetDirty(true)
		}
	}
}

var (
	x = 2
	y = 2
	v = 1
	h = 1
)

func (p *MainPage) bound(a *managers.Area) {
	counter := 0
	d := a.GetDisplay()
	t := time.NewTicker(200 * time.Millisecond)
	for true {
		select {
		case <-t.C:
			counter++
			d.PrintAtf(x, y, "*")
			x += h
			y += v

			if x >= d.Cols() {
				h = -1
				x += -2
			} else if x <= 1 {
				h = 1
				x += 2
			}

			if y >= d.Rows() {
				v = -1
				y += -2
			} else if y <= 1 {
				v = 1
				y += 2
			}
			if counter > 10 {
				counter = 0
				a.SetDirty(true)
			}
		}
	}
}
