package services

import (
	"fmt"
	"sync"
	"time"

	"ttl-monitor/internal/common"

	srl "go.bug.st/serial"
)

type Clock interface {
	High()
	Low()
}

type Status interface {
	Status(byte)
}

type Instruction interface {
	Instruction(byte)
}
type Lines interface {
	Lines([6]byte)
}

type Serial struct {
	portName     string
	buffer       chan byte
	mode         *srl.Mode
	port         srl.Port
	terminated   bool
	connected    bool
	clocks       []Clock
	statuses     []Status
	instructions []Instruction
	lines        []Lines
}

func NewSerial(cfg *common.Config) *Serial {
	s := &Serial{
		portName:     cfg.Serial.PortName,
		clocks:       make([]Clock, 0, 1),
		statuses:     make([]Status, 0, 1),
		instructions: make([]Instruction, 0, 1),
		lines:        make([]Lines, 0, 1),
		buffer:       make(chan byte),
		mode: &srl.Mode{
			DataBits: cfg.Serial.DataBits,
			BaudRate: cfg.Serial.BaudRate,
			StopBits: toStopBits(cfg.Serial.StopBits),
			Parity:   toParity(cfg.Serial.Parity),
		},
	}

	wg := &sync.WaitGroup{}
	go s.driver(wg)
	go s.portMonitor(wg)

	return s
}

func (s *Serial) RegisterClock(clock Clock) {
	s.clocks = append(s.clocks, clock)
}

func (s *Serial) RegisterStatus(status Status) {
	s.statuses = append(s.statuses, status)
}

func (s *Serial) RegisterInstruction(instruction Instruction) {
	s.instructions = append(s.instructions, instruction)
}

func (s *Serial) RegisterLines(lines Lines) {
	s.lines = append(s.lines, lines)
}

func toStopBits(value int) srl.StopBits {
	switch value {
	case 1:
		return srl.OneStopBit
	case 2:
		return srl.OnePointFiveStopBits
	case 3:
		return srl.TwoStopBits
	default:
		fmt.Println("Invalid StopBit")
		return srl.OneStopBit
	}
}

func toParity(value int) srl.Parity {
	switch value {
	case 0:
		return srl.NoParity
	case 1:
		return srl.OddParity
	case 2:
		return srl.EvenParity
	case 3:
		return srl.MarkParity
	case 4:
		return srl.SpaceParity
	default:
		fmt.Println("Invalid StopBit")
		return srl.NoParity
	}
}

func (s *Serial) portMonitor(wg *sync.WaitGroup) {
	wg.Add(1)
	defer func() {
		fmt.Println("PortMonitor Done")
		wg.Done()
	}()

	var err error
	tick := time.NewTicker(200 * time.Millisecond)
	for !s.terminated {
		select {
		case <-tick.C:
			if s.port == nil {
				if s.port, err = srl.Open(s.portName, s.mode); err != nil {
					s.port = nil
					if s.connected {
						s.connected = false
						// s.connStatus(false)
					}
				} else {
					// s.log.Infof("Opened port %s", config.CLIConfig.Serial.PortName)
					fmt.Printf("Opened port %s", s.portName)
					if !s.connected {
						// s.connStatus(true)
						s.connected = true
					}
					s.readPort()
				}
			}
		}
	}
	tick.Stop()
	close(s.buffer)
}

func (s *Serial) readPort() {
	bs := make([]byte, 100)
	defer func() {
		if r := recover(); r != nil {
			// s.log.Errorf("Recovered Read panic: %v", r)
			fmt.Printf("Recovered Read panic: %v", r)
		}
	}()
	for !s.terminated {
		if n, err := s.port.Read(bs); err != nil {
			_ = s.port.Close()
			s.port = nil
			// s.log.Infof("Lost port %s", portName)
			fmt.Printf("Lost port %s", s.portName)
			return
		} else {
			for i := 0; i < n; i++ {
				s.buffer <- bs[i]
				bs[i] = 0
			}
			// s.log.Warnf("Unexpected number of bytes received: Wanted 1, Received %d", n)
		}
	}
	close(s.buffer)
}

func (s *Serial) driver(wg *sync.WaitGroup) {
	wg.Add(1)
	defer func() {
		fmt.Println("Driver Done")
		wg.Done()
	}()
	b, ok := byte(0), true
	for {
		if b, ok = <-s.buffer; !ok {
			break
		}
		// s.log.Tracef("Inbound data: %s", string(b))
		fmt.Printf("Inbound data: %s", string(b))

		switch b {
		// case 'a':
		// 	s.address <- []byte{<-s.buffer, <-s.buffer}
		// case 'd':
		// 	s.data <- <-s.buffer // data
		// 	s.data <- <-s.buffer // error code
		// 	s.log.Tracef("Forwarding 'd' data complete")
		// case 'D':
		// 	s.data <- <-s.buffer // error code
		// 	s.log.Tracef("Forwarding 'D' data complete")
		// case 'o':
		// 	s.opCode <- <-s.buffer
		// 	s.log.Tracef("Forwarding 'o' data complete")
		case 's':
			ft := <-s.buffer
			i := <-s.buffer
			l := [6]byte{<-s.buffer, <-s.buffer, <-s.buffer, <-s.buffer, <-s.buffer, <-s.buffer}

			for _, status := range s.statuses {
				status.Status(ft)
			}
			for _, instruction := range s.instructions {
				instruction.Instruction(i)
			}
			for _, line := range s.lines {
				line.Lines(l)
			}
		case 'c':
			for _, clock := range s.clocks {
				clock.Low()
			}

		case 'C':
			for _, clock := range s.clocks {
				clock.High()
			}

		// case 'i':
		// 	s.irq.IrqLow()
		// case 'I':
		// 	s.irq.IrqHigh()
		// case 'n':
		// 	s.nmi.NmiLow()
		// case 'N':
		// 	s.nmi.NmiHigh()
		// case 'r':
		// 	s.reset.ResetLow()
		// case 'R':
		// 	s.reset.ResetHigh()
		default:
			// s.log.Warnf("Unknown byte: %v", display.HexData(b))
			fmt.Printf("Unknown byte: %v", common.HexData(b))
		}
	}
	fmt.Printf("Stopped receiving")
}
func (s *Serial) ResetChannels() {
	for {
		select {
		case <-s.buffer:
		// case <-s.address:
		// case <-s.opCode:
		// case <-s.data:
		default:
			if s.port != nil {
				_ = s.port.Close()
			}
			return
		}
	}
}
