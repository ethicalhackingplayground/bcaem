package scope

import (
	"fmt"
	"strings"
)

type ScopeElement struct {
	Target      string
	Category    string
	Description string
}

type ProgramData struct {
	Url        string
	InScope    []ScopeElement
	OutOfScope []ScopeElement
}

func PrintProgramScope(programScope ProgramData, delimiter string) {
	lines := ""
	for _, scopeElement := range programScope.InScope {
		var line string
		line += scopeElement.Target + delimiter
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
