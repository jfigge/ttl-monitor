package home

import (
	"fmt"

	"ttl-monitor/internal/common"
	"ttl-monitor/internal/models"
	"ttl-monitor/internal/services"
	"ttl-monitor/internal/services/display"
)

type HomePage struct {
	display.AbstractPage
}

func NewHomePage(s *services.Serial) *HomePage {
	page := &HomePage{}
	page.CreateRegion(0, 1, 52, 17, display.BoxDrawer) // Code
	page.CreateRegion(101, 0, 8, 1, display.BoxDrawer) // coords

	page.CreateRegion(54, 1, 21, 1, display.BoxDrawer) // Flags
	page.CreateRegion(83, 1, 3, 1, display.BoxDrawer)  // IRQ
	page.CreateRegion(93, 1, 3, 1, display.BoxDrawer)  // NMI

	page.CreateRegion(54, 4, 24, 1, display.BoxDrawer) // Steps
	page.CreateRegion(82, 4, 5, 1, display.BoxDrawer)  // Clock
	page.CreateRegion(92, 4, 5, 1, display.BoxDrawer)  // Reset

	page.CreateRegion(54, 7, 29, 11, display.BoxDrawer) // instruction
	page.CreateRegion(90, 7, 19, 4, display.BoxDrawer)  // Bus content
	page.CreateRegion(90, 13, 19, 4, display.BoxDrawer) // ALU

	page.CreateRegion(0, 20, 63, 15, display.BoxDrawer) // Lines
	page.CreateRegion(65, 20, 40, 7, display.BoxDrawer) // Active Lines

	page.CreateRegion(0, 35, 109, 7, display.BoxDrawer) // Logs
	return page
}

func (h *HomePage) ProcessInput(keyInput *models.KeyInput) bool {
	if keyInput.Ascii == 'q' {
		return true
	}
	fmt.Printf("Key Code: %d, Ascii Code: %d\n", keyInput.KeyCode, keyInput.Ascii)
	return false
}

func (h *HomePage) Draw(canvas *display.Canvas) {
	canvas.PrintAtf(
		6, 1, "%s0  1  2  3  4  5  6  7   8  9  A  B  C  D  E  F        Flags                  IRQ       NMI%s",
		common.Yellow,
		common.Reset,
	)
	canvas.PrintAtf(61, 4, "%sStep                  Clock     Reset%s", common.Yellow, common.Reset)
	canvas.PrintAtf(58, 7, "%sInstructions                    Bus Content%s", common.Yellow, common.Reset)
	canvas.PrintAtf(86, 8, "%s DB:%s", common.Yellow, common.Reset)
	canvas.PrintAtf(86, 9, "%sABH:%s", common.Yellow, common.Reset)
	canvas.PrintAtf(86, 10, "%sABL:%s", common.Yellow, common.Reset)
	canvas.PrintAtf(86, 11, "%s SB:%s", common.Yellow, common.Reset)

	canvas.PrintAtf(91, 13, "%sALU%s", common.Yellow, common.Reset)
	canvas.PrintAtf(86, 14, "%s  B:%s", common.Yellow, common.Reset)
	canvas.PrintAtf(86, 15, "%s  A:%s", common.Yellow, common.Reset)
	canvas.PrintAtf(86, 16, "%s Op:%s", common.Yellow, common.Reset)
	canvas.PrintAtf(86, 17, "%sDir:%s", common.Yellow, common.Reset)

	canvas.PrintAtf(1, 20, "%sControl Lines:%s", common.Yellow, common.Reset)
	canvas.PrintAtf(66, 20, "%sActive Lines:%s", common.Yellow, common.Reset)
}
