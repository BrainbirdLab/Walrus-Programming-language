package parser

import (
	"fmt"
	"os"
	"strings"
	"walrus/frontend/lexer"
	"walrus/utils"
)

type htype string

const (
	TEXT_HINT htype = "text_hint"
	CODE_HINT htype = "code_hint"
)

type HintType struct {
	HType htype
	HText string
}

type ErrorMessage struct {
	Message string
	hints   []HintType
}

func (e *ErrorMessage) AddHint(htext string, htype htype) *ErrorMessage {
	e.hints = append(e.hints, HintType{
		HText: htext,
		HType: htype,
	})
	return e
}

func (e *ErrorMessage) Display() {
	fmt.Print(e.Message)
	// hints
	for i, hint := range e.hints {
		if i == 0 {
			fmt.Print(utils.Colorize(utils.ORANGE, "Hint: "))
		}
		//fmt.Print(utils.Colorize( utils.ORANGE, (fmt.Sprintf("Hint: %s\n", hint.HText))))
		if hint.HType == TEXT_HINT {
			fmt.Print(utils.Colorize(utils.ORANGE, hint.HText))
		} else {
			fmt.Print(lexer.Highlight(hint.HText))
		}
	}
	fmt.Println("")
	os.Exit(1)
}

func makePadding(width, line int) string {
	return fmt.Sprintf("%*d | ", width, line)
}

func MakeError(p *Parser, lineNo int, filePath string, startPos lexer.Position, endPos lexer.Position, errMsg string) *ErrorMessage {

	// decorate the error with ~~~~~ under the error line

	var errStr string

	var prvLines []string
	var nxtLines []string
	line := (*p.Lines)[lineNo-1]
	maxWidth := len(fmt.Sprintf("%d", len(*p.Lines)))

	if lineNo-1 > 0 {
		// add the padding to each line
		prvLines = (*p.Lines)[lineNo-2 : lineNo-1]

		for i, l := range prvLines {
			prvLines[i] = utils.Colorize(utils.GREY, makePadding(maxWidth, lineNo-1+i) + lexer.Highlight(l))
		}
	}

	if lineNo+1 < len(*p.Lines) {
		nxtLines = (*p.Lines)[lineNo : lineNo+1]

		for i, l := range nxtLines {
			nxtLines[i] = utils.Colorize(utils.GREY, makePadding(maxWidth, lineNo+1+i)) + lexer.Highlight(l)
		}
	}

	errStr += fmt.Sprintf("\nIn file: %s:%d:%d\n", filePath, startPos.Line, startPos.Column)

	padding := makePadding(maxWidth, startPos.Line)

	errStr += strings.Join(prvLines, "\n") + "\n"
	errStr += utils.Colorize(utils.GREY, padding) + lexer.Highlight(line[0:startPos.Column-1]) + utils.Colorize(utils.RED, line[startPos.Column-1 : endPos.Column-1]) + lexer.Highlight(line[endPos.Column-1 : ]) + "\n"
	errStr += strings.Repeat(" ", (startPos.Column-1)+len(padding))
	errStr += fmt.Sprint(utils.Colorize(utils.BOLD_RED, fmt.Sprintf("%s%s\n", "^", strings.Repeat("~", (endPos.Column - startPos.Column) - 1))))
	errStr += strings.Join(nxtLines, "\n") + "\n"
	errStr += fmt.Sprint(utils.Colorize(utils.RED, fmt.Sprintf("Error: %s\n", errMsg)))

	return &ErrorMessage{
		Message: errStr,
	}
}
