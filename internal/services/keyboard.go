package services

import (
	"fmt"
	"github.com/pkg/term"
	"sync"
	"time"
	"ttl-monitor/internal/common"
	"ttl-monitor/internal/models"
	"ttl-monitor/internal/services/display"
)

type Keyboard struct {
	waitGroup  sync.WaitGroup
	terminated bool
	terminal   *display.Terminal
	ioTerm     *term.Term
}

func NewKeyboard(t *display.Terminal) *Keyboard {
	return &Keyboard{
		terminal:   t,
		terminated: true,
	}
}

func (k *Keyboard) Monitor() {
	if k.terminated {
		k.terminated = false
		k.waitGroup.Add(1)
		go func() {
			for !k.terminated {
				keys, err := k.ReadChar()
				if err != nil {
					continue
				}
				pageMeta := k.terminal.ActivePage()
				if pageMeta != nil {
					pageMeta.ProcessInput(keys)
				} else {
					common.Debugf("Lost input.  Keycode: %d / Ascii key: %d", keys.KeyCode, keys.Ascii)
				}
			}
			k.waitGroup.Done()
		}()
	}
}

func (k *Keyboard) ReadChar() (*models.KeyInput, error) {
	k.ioTerm, _ = term.Open("/dev/tty")
	if k.ioTerm == nil {
		return nil, fmt.Errorf("terminal unavailable")
	} else {
		defer func() {
			if err := k.ioTerm.Restore(); err != nil {
				common.Debugf("Failed to restore terminal mode: %v", err)
			}
			if err := k.ioTerm.Close(); err != nil {
				common.Debugf("Failed to close terminal input: %v", err)
			}
		}()
	}

	if err := term.RawMode(k.ioTerm); err != nil {
		return nil, common.Errorf("Failed to access terminal RawMode: %v", err)
	}
	bs := make([]byte, 5)

	if err := k.ioTerm.SetReadTimeout(100 * time.Millisecond); err != nil {
		common.Warn("Failed to set read timeout")
	}

	var numRead int
	var err error
	keyInput := models.KeyInput{}
	if numRead, err = k.ioTerm.Read(bs); err != nil {
		if err.Error() != "EOF" {
			if err = k.ioTerm.Restore(); err != nil {
				err = common.Errorf("Failed to restore terminal mode: %v", err)
			}
			if err = k.ioTerm.Close(); err != nil {
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
			text = fmt.Sprintf("%s %d:%s", text, i, common.HexData(bs[i]))
		}
		common.Warnf(text)
		// Two characters read??
	}
	return &keyInput, nil
}

func (k *Keyboard) Terminate() {
	if !k.terminated {
		k.terminated = true
		k.waitGroup.Wait()
	}
}
