package db

import (
	"io"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

type FieldsPlugger struct {
	fieldValue *protogen.Field
	fields     []*strings.Builder
}

func (fp *FieldsPlugger) BorrowWriter() io.Writer {
	fp.fields = append(fp.fields, &strings.Builder{})
	return fp.fields[len(fp.fields)-1]
}

func (fp *FieldsPlugger) GetName() string {
	return fp.fieldValue.GoName
}

func (fp *FieldsPlugger) GetGoType() string {
	// todo: type map
	return ""
}

func (fp *FieldsPlugger) GetContent() string {
	ss := make([]string, len(fp.fields))
	for i, field := range fp.fields {
		ss[i] = field.String()
	}
	return strings.Join(ss, " ")
}
