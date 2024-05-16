package utils

import (
	"regexp"
	"strconv"
)

func BitSizeFromString(value string) uint8 {
	//extract the number from the string
	regexp := regexp.MustCompile(`\d+`)
	match := regexp.FindString(value)
	if match == "" {
		panic("Invalid bit size")
	}
	size, err := strconv.Atoi(match)

	if err != nil {
		panic("Invalid bit size")
	}

	return uint8(size)
}