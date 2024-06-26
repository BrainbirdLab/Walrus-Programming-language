package helpers

import (
	"fmt"
	"math"
	"strings"
	"walrus/utils"
)

func GetCodePreview(source string, start int, end int, line int) (string, string) {
	// Split the source code into lines
	lines := strings.Split(source, "\n")

	// Calculate the window range
	windowSize := 2
	startLine := utils.Max(0, line-windowSize-1)
	endLine := utils.Min(len(lines), line+windowSize)

	// Get the lines within the window range
	windowLines := lines[startLine:endLine]

	// Find the maximum line number length
	maxLineNumLen := len(fmt.Sprintf("%d", endLine))

	// Add line numbers and pipe separators to the window lines
	var previewLines []string
	for i, line := range windowLines {
		lineNum := startLine + i + 1
		lineNumStr := fmt.Sprintf("%*d", maxLineNumLen, lineNum)
		lineStr := fmt.Sprintf("%s | %s", lineNumStr, line)
		previewLines = append(previewLines, lineStr)
	}

	// Join the preview lines into a single string
	preview := strings.Join(previewLines, "\n")

	// Get the error line
	errorLine := lines[line-1]

	// Calculate the start position relative to the error line
	lineStart := strings.Index(source, errorLine)
	column := int(math.Abs(float64(start - lineStart)))

	println(column)
	// Create the cursor line
	cursorLine := "     " + strings.Repeat(" ", column) + strings.Repeat("^", end+1-start)

	// Combine the error line and cursor line
	errorPreview := fmt.Sprintf("%s\n%s", errorLine, cursorLine)

	return preview, errorPreview
}
