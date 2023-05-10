package display

import (
	"fmt"
	"math/rand"
	"strings"
)

type Renderer func(d *Canvas)

func NoDraw(d *Canvas) {
}

func BoxDrawer(d *Canvas) {
	w := d.Cols()
	h := d.Rows()
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
		} else if i != 0 {
			d.PrintAtf(1, i+1, "#%s#", strings.Repeat(" ", w-2))
		}
	}
}

func RandomDrawer(d *Canvas) {
	w := d.Cols()
	h := d.Rows()
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
