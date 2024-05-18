package utils

import (
	"fmt"
	"regexp"
)

// ANSI color escape codes
const (
    RESET  		=  	"\033[0m"
    RED    		= 	"\033[31m"
	BOLD_RED 	= 	"\033[38;05;196m"
    GREEN  		= 	"\033[32m"
    YELLOW 		= 	"\033[33m"
	ORANGE 		= 	"\033[38;05;221m"
    BLUE   		= 	"\033[34m"
    PURPLE 		= 	"\033[35m"
    CYAN   		= 	"\033[36m"
    WHITE  		= 	"\033[37m"
	GREY   		= 	"\033[90m"
	BOLD   		= 	"\033[1m"
)

func Colorize(color, text string) string {

	//check the format
	//regex := regexp.MustCompile(`\033\[([0-9]+)([;]*)m`)

	regex := regexp.MustCompile(`\033\[[0-9]{1,3}[0-9;]*m`)

	//panic if the color is not valid
	if !regex.MatchString(color) {
		panic("Invalid color")
	}

	return fmt.Sprintf("%s%s%s", color, text, RESET)
}


func IF(conddition bool, a, b any) any {
	if conddition {
		return a
	}
	return b
}

