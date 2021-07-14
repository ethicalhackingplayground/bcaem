package scope

import (
	"fmt"
	"log"
	"strings"
)

type ScopeElement struct {
	Target string
}

type ProgramData struct {
	Url string
}

func PrintProgramScope(programScope ProgramData, outputFlags string, delimiter string) {
	lines := ""
	for _, scopeElement := range programScope.InScope {
		var line string
		line += programScope.Url + delimiter
		line = strings.TrimSuffix(line, delimiter)
		if len(line) > 0 {
			lines += line + "\n"
		}
	}

	lines = strings.TrimSuffix(lines, "\n")

	if len(lines) > 0 {
		fmt.Println(lines)
	}
}
