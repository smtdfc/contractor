package exception

import (
	"fmt"
	"strings"
)

func PrintException(ex IException, code string) {
	loc := ex.GetLoc()
	msg := ex.GetMsg()

	startLine, startCol := loc.GetStart()
	file := loc.GetFile()

	fmt.Printf("\033[31;1merror\033[0m: %s\n", msg)

	fmt.Printf("  \033[34;1m-->\033[0m %s:%d:%d\n", file, startLine, startCol)

	lines := strings.Split(code, "\n")

	if startLine > 0 && startLine <= len(lines) {
		errorLine := lines[startLine-1]

		fmt.Printf("\033[34;1m%5d |\033[0m %s\n", startLine, errorLine)

		indent := ""
		for i, char := range errorLine {
			if i+1 >= startCol {
				break
			}
			if char == '\t' {
				indent += "\t"
			} else {
				indent += " "
			}
		}

		fmt.Printf("      \033[34;1m|\033[0m %s\033[31;1m^\033[0m\n", indent)
	}
	fmt.Println()
}
