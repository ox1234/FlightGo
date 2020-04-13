package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	// Green is Success Logger.
	Green *log.Logger

	// Blue is Note Logger.
	Blue *log.Logger

	// Red is Error Logger.
	Red *log.Logger
)

const (
	red = uint8(iota + 91)
	green
	yellow
	blue
	magenta
)

func color(color uint8, str string) string {
	return fmt.Sprintf("\x1b[1;%dm%s\x1b[0m", color, str)
}

func init() {
	Blue = log.New(os.Stdout, color(blue, "[*] "), 0)
	Green = log.New(os.Stdout, color(green, "[+] "), 0)
	Red = log.New(os.Stdout, color(red, "[!] "), 0)
}
