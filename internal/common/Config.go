package common

const (
	defSerialBaudRate  = 115200
	defSerialDataBits  = 8
	defSerialStopBits  = 1
	defMinimumReadSize = 8
	defSerialParity    = 0

	defTerminalWidth  = 80
	defTerminalHeight = 50
)

func DefaultConfig() *Config {
	return &Config{
		Screen: &Screen{
			Width:  defTerminalWidth,
			Height: defTerminalHeight,
		},
		Serial: &Serial{
			BaudRate:        defSerialBaudRate,
			DataBits:        defSerialDataBits,
			StopBits:        defSerialStopBits,
			Parity:          defSerialParity,
			MinimumReadSize: defMinimumReadSize,
		},
	}
}

type Config struct {
	Screen *Screen `yaml:"screen"`
	Serial *Serial `yaml:"serial"`
}

type Screen struct {
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

type Serial struct {
	PortName        string `yaml:"port_name"`
	BaudRate        int    `yaml:"baud_rate"`
	DataBits        int    `yaml:"data_bits"`
	StopBits        int    `yaml:"stop_bits"`
	Parity          int    `yaml:"parity"`
	MinimumReadSize int    `yaml:"minimum_read_size"`
}
