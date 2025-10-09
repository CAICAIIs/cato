package models

import (
	"fmt"
	"strings"
)

type Import struct {
	Alias      string
	ImportPath string
}

func (cip *Import) GetPath() string {
	value := fmt.Sprintf("\"%s\"", cip.ImportPath)
	if cip.Alias != "" {
		value = fmt.Sprintf("%s \"%s\"", cip.Alias, cip.ImportPath)
	}
	return value
}

func (cip *Import) IsSame(i *Import) bool {
	if i == nil {
		return cip == nil
	}
	return cip.ImportPath == i.ImportPath
}

func (cip *Import) IsEmpty() bool {
	return cip.ImportPath == ""
}

func (cip *Import) Init(importPath string) *Import {
	cip.ImportPath = importPath
	pattern := strings.Split(importPath, "/")
	length := len(pattern)
	if length <= 2 {
		return cip
	}
	cip.Alias = strings.Join([]string{pattern[length-2], pattern[length-1]}, "")
	return cip
}
