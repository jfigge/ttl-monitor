package main

import (
	"os"
	"time"
	"ttl-monitor/internal/common"
	"ttl-monitor/internal/managers"
	"ttl-monitor/internal/pages/helppage"
	"ttl-monitor/internal/pages/historypage"
	"ttl-monitor/internal/pages/mainpage"
)

var (
	mainPage    managers.PageId
	historyPage managers.PageId
	helpPage    managers.PageId
)

func main() {
	t, err := managers.NewTerminal(40, 10)
	if err != nil {
		os.Exit(1)
	}

	mainPage = t.AddPage(mainpage.NewMainPage())
	historyPage = t.AddPage(historypage.NewHistoryPage())
	helpPage = t.AddPage(helppage.NewHeLpPage())

	go jibberish(t)
	go input(t)
	t.ShowPage(mainPage)
	t.Wait()
	t.Terminate()
}

func jibberish(t *managers.Terminal) {
	counter := 0
	for true {
		time.Sleep(1 * time.Second)
		counter++
		if counter > 10 {
			counter = 0
			t.Render(true)
		} else {
			t.Render(false)
		}
	}
	return
}

func input(t *managers.Terminal) {
	for true {
		keys, err := t.ReadChar()
		if err != nil {
			continue
		}
		pageMeta := t.ActivePage()
		if pageMeta != nil {
			pageMeta.ProcessInput(keys)
		} else {
			common.Debugf("Lost input.  Keycode: %d / Ascii key: %d", keys.KeyCode, keys.Ascii)
		}
	}
}
