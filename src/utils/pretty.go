package utils

import (
	"fmt"
	"regexp"
	"strconv"
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

func GetIntBitSize(rawValue string) uint8 {

	size := uint8(32)

	if len(rawValue) > 10 {
		size = 64
	} else {
		// if number is out of range for 32-bit integer
		// then it is a 64-bit integer
		// But to avoid checking both positive and negative ranges, we just check the positive range by using the absolute value
		// if the absolute value is greater than 2,147,483,647 then it is a 64-bit integer
		number, _ := strconv.ParseInt(rawValue, 10, 32)
		if number < 0 {
			number = -number
		}
		if number > 2147483647 {
			size = 64
		}
	}

	return size
}

func GetFloatBitSize(rawValue string) uint8 {
	size := uint8(32)

	number, _ := strconv.ParseFloat(rawValue, 64)

	if number < 0 {
		number = -number
	}

	// check the floating point decimal size

	decimal := int64(number)

	//max size of a 32-bit floating point number is 7 digits
	if decimal > 9999999 {
		size = 64
	}

	return size
}