package main

import (
	"fmt"
	"os"
	"time"

	"ttl-monitor/internal/common"
	"ttl-monitor/internal/pages"
	"ttl-monitor/internal/pages/home"
	"ttl-monitor/internal/services"
	"ttl-monitor/internal/services/display"

	"gopkg.in/yaml.v3"
)

type App struct {
	Config      *common.Config
	mainPage    display.PageId
	historyPage display.PageId
	helpPage    display.PageId
	terminal    *display.Terminal
	keyboard    *services.Keyboard
	serial      *services.Serial
	terminate   chan struct{}
}

func main() {
	app := &App{
		terminate: make(chan struct{}),
	}
	app.initConfig()
	app.initServices()
	app.addPages()
	app.Run(app.mainPage)
	app.Terminate()
	fmt.Println("Done")
}

func (a *App) initConfig() {
	bs, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Sprintf("failed to load config file: %v\n", err)
		os.Exit(1)
	}

	a.Config = common.DefaultConfig()
	err = yaml.Unmarshal(bs, a.Config)
	if err != nil {
		fmt.Sprintf("failed to unmarshal config file: %v\n", err)
		os.Exit(1)
	}
}

func (a *App) initServices() {
	var err error
	a.terminal, err = display.NewTerminal(a.Config.Screen.Width, a.Config.Screen.Height)
	if err != nil {
		os.Exit(1)
	}
	time.Sleep(100 * time.Millisecond)

	a.keyboard = services.NewKeyboard(a.terminal)
	a.keyboard.Monitor()

	a.serial = services.NewSerial(a.Config)
}

func (a *App) addPages() {
	a.mainPage = a.terminal.AddPage(home.NewHomePage(a.serial))
	a.historyPage = a.terminal.AddPage(pages.NewHistoryPage(a.serial))
	a.helpPage = a.terminal.AddPage(pages.NewHelpPage(a.serial))
}

func (a *App) Run(page display.PageId) {
	a.terminal.ShowPage(page)
	go func() {
		a.terminal.Wait()
		a.terminate <- struct{}{}
	}()
	select {
	case <-a.terminate:
	}
}

func (a *App) Terminate() {
	a.terminal.Terminate()
	a.keyboard.Terminate()
}
