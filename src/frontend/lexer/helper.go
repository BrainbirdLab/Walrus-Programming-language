package lexer

import (
	"strings"

	"walrus/utils"
)

func Highlight(line string) string {
	words := strings.Split(line, " ")

	for i, word := range words {
		if IsKeyword(TOKEN_KIND(word)) || IsBuiltInType(TOKEN_KIND(word)) {
			words[i] = utils.Colorize(utils.PURPLE, word)
		} else if IsNumber(TOKEN_KIND(word)) {
			words[i] = utils.Colorize(utils.ORANGE, word)
		}
	}

	return strings.Join(words, " ")
}