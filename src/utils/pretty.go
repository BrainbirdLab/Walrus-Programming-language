package utils

import "fmt"

// ANSI color escape codes
const (
    RESET  =  "\033[0m"
    RED    = "\033[31m"
    GREEN  = "\033[32m"
    YELLOW = "\033[33m"
    BLUE   = "\033[34m"
    PURPLE = "\033[35m"
    CYAN   = "\033[36m"
    WHITE  = "\033[37m"
)

func Colorize(color, text string) string {
	switch color {
		case RED, GREEN, YELLOW, BLUE, PURPLE, CYAN, WHITE:
			return fmt.Sprintf("%s%s%s", color, text, RESET)
		default:
			panic("Invalid color")
	}
}



func IF(conddition bool, a, b any) any {
	if conddition {
		return a
	}
	return b
}

