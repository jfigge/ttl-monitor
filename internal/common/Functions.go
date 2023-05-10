package common

import "fmt"

var (
	HEX = [16]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
)

func BinData(data uint8) string {
	return fmt.Sprintf("%08b", data)
}
func HexData(data uint8) string {
	return fmt.Sprintf("%s%s", HEX[data>>4], HEX[data&15])
}
func HexAddress(address uint16) string {
	return fmt.Sprintf("%s%s%s%s", HEX[address>>12], HEX[address>>8&15], HEX[address>>4&15], HEX[address&15])
}
