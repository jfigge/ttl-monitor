package display

import (
	"fmt"
	"strings"
)

type Renderer func(d *Display)

func BoxDrawer(w int, h int) func(d *Display) {
	return func(d *Display) {
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
	}
}

func ClearArea(w int, h int) func(d *Display) {
	return func(d *Display) {
		blank := strings.Repeat(" ", w)
		for i := 0; i < h; i++ {
			d.PrintAtf(1, i+1, blank)
		}
	}
}
