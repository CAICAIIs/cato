package models

import (
	"fmt"
	"strings"

	"github.com/ncuhome/cato/generated"
)

type Field struct {
	Name   string
	GoType string
}

func (f *Field) AsParamName() string {
	return strings.ToLower(f.Name)
}

type Key struct {
	// this represents from field and type
	KeyName string
	KeyType generated.DBKeyType
	Fields  []*Field
}

func (k *Key) GetFieldNameCombine() string {
	filedNames := make([]string, len(k.Fields))
	for i, f := range k.Fields {
		filedNames[i] = f.Name
	}
	return strings.Join(filedNames, "And")
}

func (k *Key) GetParamsRaw() []string {
	data := make([]string, len(k.Fields))
	for index, field := range k.Fields {
		data[index] = fmt.Sprintf("%s %s", field.AsParamName(), field.GoType)
	}
	return data
}

type Col struct {
	ColName string
	*Field
}
