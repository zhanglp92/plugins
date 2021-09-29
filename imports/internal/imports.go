package internal

import (
	"golang.org/x/tools/imports"
	"strings"
)

func init() {
	updateImportToGroupHandlers()
}

// Process ...
func Process(data []byte) ([]byte, error) {
	out, err := imports.RemoveImportSpaces(data)
	if err != nil {
		return data, err
	}
	return imports.Process("", out, &imports.Options{
		Comments:   true,
		TabIndent:  true,
		FormatOnly: true,
		TabWidth:   8,
	})
}

func updateImportToGroupHandlers() {
	imports.AppendHeadImportToGroupHandlers(
		func(_, importPath string) (num int, ok bool) {
			firstComponent := strings.Split(importPath, "/")[0]
			if firstComponent == "gitlab.myteksi.net" {
				return 100, true
			}
			return
		},
	)
}
