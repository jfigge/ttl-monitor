package old

// import (
// "fmt"
// "github.com/pkg/term"
// "github.td.teradata.com/sandbox/logic-ctl/internal/config"
// "github.td.teradata.com/sandbox/logic-ctl/internal/services/common"
// "github.td.teradata.com/sandbox/logic-ctl/internal/services/display"
// "github.td.teradata.com/sandbox/logic-ctl/internal/services/logging"
// "github.td.teradata.com/sandbox/logic-ctl/internal/services/serial"
// "github.td.teradata.com/sandbox/logic-ctl/internal/services/status"
// "os"
// "strings"
// "sync"
// "time"
// )
//
// type Driver struct {
// 	page         common.Page
// 	connected    bool
//
// 	state        *status.State
// 	serial       *serial.Serial
// 	log          *logging.Log
//
// 	monitorChan  chan bool
// 	inputChan    chan common.Input
// 	xTerm        *term.Term
// 	keyIntercept []common.Intercept
// 	wg           *sync.WaitGroup
// }
//
// func New() *Driver {
//
// 	screen, err := display.Initialize()
// 	if err != nil {
// 		fmt.Printf("%sunable to initialize terminal: %v%s", common.Red, err, common.Reset)
// 		os.Exit(1)
// 	}
//
// 	if screen.Rows() < 33 + config.CLIConfig.Log.Height || screen.Cols() < 100 {
// 		fmt.Printf("%sMinimum console size must be 100x38.  Currently at %dx%d%s\n", common.Red, screen.Cols(), screen.Rows(), common.Reset)
// 		os.Exit(1)
// 	}
//
// 	d := Driver{
// 		wg:          &sync.WaitGroup{},
// 		page:        screen.CreatePage(),
// 		monitorChan: make(chan bool),
// 		inputChan:   make(chan common.Input),
// 	}
// 	d.page.SetInputProcessor(d.Process)
// 	logRegion := d.page.CreateRegion(0, screen.Rows() - config.CLIConfig.Log.Height, screen.Cols(), config.CLIConfig.Log.Height)
// 	//d.page.CreateRegion(screen.Cols()/3,2, screen.Cols()/3, 2).SetRenderer(d.StatusDrawer)
//
// 	d.log          = logging.New(logRegion)
// 	//d.state        = status.NewState(d.log, d.redraw)
// 	d.serial       = serial.New(d.log, d.state, d.connectionStatus, d.wg)
// 	return &d
// }
//
// func (d *Driver) Run() {
// 	time.Sleep(1 * time.Second)
// 	d.page.Activate()
// 	go d.input(d.wg)
//
// 	for display.Screen.ActivePage() != nil {
// 		if a, k, e := d.ReadChar(); e != nil {
// 			d.log.Warn(e.Error())
// 		} else if a != 0 || k != 0 {
// 			d.inputChan <- common.Input{Ascii: a, KeyCode: k, Connected: d.connected}
// 		}
// 	}
//
// 	close(d.monitorChan)
//
// 	d.wg.Wait()
// 	fmt.Println("Wait Done")
// }
//
// func (d *Driver) input(wg *sync.WaitGroup) {
// 	wg.Add(1)
// 	defer func() {
// 		fmt.Println("Input Done")
// 		wg.Done()
// 	}()
// 	input, connected, ok := common.Input{}, false, true
//
// 	// By monitoring all state changes in one thread we can
// 	// ensure all updates are serialized
// 	for ok {
// 		select {
// 		case connected, ok = <-d.monitorChan:
// 			d.connected = connected
// 			//d.redraw(true)
//
// 		case input, ok = <-d.inputChan:
// 			page := display.Screen.ActivePage()
// 			if page != nil && page.ProcessesInput() && page.ProcessInput(input) {
// 				display.Screen.DeactivatePage(page)
// 			}
// 		}
// 	}
// 	d.serial.Terminate()
// }
//
// func (d *Driver) ReadChar() (ascii int, keyCode int, err error) {
// 	d.xTerm, _ = term.Open("/dev/tty")
// 	if d.xTerm == nil {
// 		return 0, 0, fmt.Errorf("terminal unavailable")
// 	} else {
// 		defer func() {
// 			if err := d.xTerm.Restore(); err != nil {
// 				d.log.Errorf("Failed to restore terminal mode: %v", err)
// 			}
// 			if err := d.xTerm.Close(); err != nil {
// 				d.log.Errorf("Failed to close terminal input: %v", err)
// 			}
// 		}()
// 	}
//
// 	if err = term.RawMode(d.xTerm); err != nil {
// 		str := fmt.Sprintf("Failed to access terminal RawMode: %v", err)
// 		d.log.Error(str)
// 		return
// 	}
// 	bs := make([]byte, 5)
//
// 	if err := d.xTerm.SetReadTimeout(2 * time.Second); err != nil {
// 		d.log.Warn("Failed to set read timeout")
// 	}
// 	if numRead, err := d.xTerm.Read(bs); err != nil {
// 		if err.Error() != "EOF" {
// 			if err := d.xTerm.Restore(); err != nil {
// 				d.log.Errorf("Failed to restore terminal mode: %v", err)
// 			}
// 			if err := d.xTerm.Close(); err != nil {
// 				d.log.Errorf("Failed to close terminal input: %v", err)
// 			}
// 		}
// 		return 0, 0, nil
// 	} else if numRead == 3 && bs[0] == 27 && bs[1] == 91 {
// 		// Three-character control sequence, beginning with "ESC-[".
//
// 		// Since there are no ASCII lines for arrow keys, we use
// 		// Javascript key lines.
// 		switch bs[2] {
// 		case 65: keyCode = 38 // Up
// 		case 66: keyCode = 40 // Down
// 		case 67: keyCode = 39 // Right
// 		case 68: keyCode = 37 // Left
// 		}
// 	} else if numRead == 3 && bs[0] == 0x1B && bs[1] == 0x4F {
// 		switch bs[2] {
// 		case 50: keyCode = 101 // Option+1
// 		case 51: keyCode = 102 // Option+2
// 		case 52: keyCode = 102 // Option+3
// 		case 53: keyCode = 103 // Option+4
// 		}
// 	} else if numRead == 1 {
// 		ascii = int(bs[0])
// 	} else {
// 		text := fmt.Sprintf("%d characters read.", numRead)
// 		for i := 0; i < numRead; i++ {
// 			text = fmt.Sprintf("%s %d:%s", text, i, display.HexData(bs[i]))
// 		}
// 		d.log.Warnf(text)
// 		// Two characters read??
// 	}
// 	return
// }
//
// func (d *Driver) connectionStatus(connected bool) {
// 	// Remove any previous value before it can be read
// 	select {
// 	case <-d.monitorChan:
// 	default: // don't block if there is no value
// 	}
//
// 	// Push the current connected status
// 	d.monitorChan <- connected
// }
// func (d *Driver) Draw() {
// }
//
// func (d *Driver) Process(input common.Input) bool {
// 	if input.KeyCode != 0 {
// 		switch input.KeyCode {
// 		default:
// 			d.log.Warnf("Unknown code: [%v]", input.KeyCode)
// 		}
// 	} else {
// 		switch input.Ascii {
// 		case 'q':
// 			return true
// 		case 'h':
// 			d.log.HistoryPage().Activate()
// 		case 'p':
// 			d.serial.PortPage().Activate()
// 		default:
// 			d.log.Warnf("Unmapped ascii code: [%c]", input.Ascii)
// 		}
// 	}
// 	return false
// }
//
// func (d *Driver) Redraw() {
//
// }
//
// func (d *Driver) StatusDrawer(display common.Display, cols, rows int) {
// 	line := fmt.Sprintf("%s%s%s", common.BGRed, strings.Repeat(" ", cols), common.Reset)
// 	for row := 1; row <= rows; row++ {
// 		display.PrintAt(1, row, line)
// 	}
// }
//

// func Initialize() (common.Screen, error) {
// 	if term == nil {
// 		if !xterm.IsTerminal(int(os.Stdin.Fd())) {
// 			return nil, errors.New("not a terminal")
// 		}
// 		term = &terminal{
// 			fd: int(os.Stdin.Fd()),
// 			activePage: make([]common.Page, 0),
// 			fullScreen: true,
// 		}
// 		for term.cols == 0 {
// 			if w, h, e := xterm.GetSize(int(os.Stdin.Fd())); e != nil {
// 				return nil, e
// 			} else {
// 				term.x    = 1
// 				term.y    = 1
// 				term.cols = w
// 				term.rows = h
// 			}
// 		}
//
// 		if s, e := xterm.GetState(int(os.Stdin.Fd())); e != nil {
// 			return nil, e
// 		} else {
// 			term.state = s
// 		}
// 		term.HideCursor()
// 		Screen = term
// 		return term, nil
// 	}
// 	return nil, errors.New("already initialized")
// }
