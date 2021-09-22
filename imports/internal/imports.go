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
	return imports.Process("", data, nil)
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