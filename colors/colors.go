package colors

import (
	"fmt"
	"runtime"
)

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var purple = "\033[35m"
var cyan = "\033[36m"
var gray = "\033[37m"
var white = "\033[97m"

func init() {
	// Windows is special, as usual
	if runtime.GOOS == "windows" {
		reset = ""
		red = ""
		green = ""
		yellow = ""
		blue = ""
		purple = ""
		cyan = ""
		gray = ""
		white = ""
	}
}

func colorString(s string, c string) string {
	return fmt.Sprintf("%s%s%s", c, s, reset)
}

// Red wraps colorString
func Red(s string) string {
	return colorString(s, red)
}

// Green wraps colorString
func Green(s string) string {
	return colorString(s, green)
}

// Yellow wraps colorString
func Yellow(s string) string {
	return colorString(s, yellow)
}

// Blue wraps colorString
func Blue(s string) string {
	return colorString(s, blue)
}

// Purple wraps colorString
func Purple(s string) string {
	return colorString(s, purple)
}

// Cyan wraps colorString
func Cyan(s string) string {
	return colorString(s, cyan)
}

// Gray wraps colorString
func Gray(s string) string {
	return colorString(s, gray)
}

// White wraps colorString
func White(s string) string {
	return colorString(s, white)
}
