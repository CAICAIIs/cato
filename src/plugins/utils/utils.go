package utils

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

func GetGoPackageName(importPath protogen.GoImportPath) string {
	patterns := strings.Split(strings.Trim(importPath.String(), "\""), "/")
	if len(patterns) == 0 || patterns[0] == "." {
		return ""
	}
	return patterns[len(patterns)-1]
}
