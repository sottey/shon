package pkg

import (
	"fmt"
	"strings"
)

// Print line if verbose is true
func DebugPrint(message string, verbose bool) {
	if verbose {
		fmt.Println(message)
	}
}

// Indent line according to indentation level specified
func IndentLine(level, spaces int, line string) string {
	return strings.Repeat(" ", level*spaces) + line
}
