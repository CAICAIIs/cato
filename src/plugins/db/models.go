package db

import (
	"fmt"
	"io"
	"log"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/utils"
)

type ModelsPlugger struct {
	message *protogen.Message
	parent  *protogen.File

	fields  map[string]*FieldsPlugger
	imports []*strings.Builder
	methods []*strings.Builder
	extra   []*strings.Builder
	tmpl    *template.Template
}

type ModelsPluggerPack struct {
	PackageName string
	Imports     []string
	ModelName   string
	Fields      []string
	Methods     []string
}

func (mp *ModelsPlugger) LoadContext(message *protogen.Message, file *protogen.File) {
	mp.message = message
	mp.parent = file
	mp.fields = make(map[string]*FieldsPlugger)
	mp.imports = make([]*strings.Builder, 0)
	mp.methods = make([]*strings.Builder, 0)
	mp.extra = make([]*strings.Builder, 0)
}

func (mp *ModelsPlugger) findField(name string) (*protogen.Field, bool) {
	for _, field := range mp.message.Fields {
		if string(field.Desc.Name()) == name {
			return field, true
		}
	}
	return nil, false
}

func (mp *ModelsPlugger) BorrowFieldsWriter(name string) (io.Writer, bool) {
	_, ok := mp.fields[name]
	if !ok {
		fieldDesc, ok := mp.findField(name)
		if !ok {
			return nil, false
		}
		mp.fields[name] = &FieldsPlugger{fieldDesc, make([]*strings.Builder, 0)}
	}
	return mp.fields[name].BorrowWriter(), true
}

func (mp *ModelsPlugger) BorrowMethodsWriter() io.Writer {
	mp.methods = append(mp.methods, new(strings.Builder))
	return mp.methods[len(mp.methods)-1]
}

func (mp *ModelsPlugger) BorrowImportsWriter() io.Writer {
	mp.imports = append(mp.imports, new(strings.Builder))
	return mp.imports[len(mp.imports)-1]
}

func (mp *ModelsPlugger) BorrowExtraWriter() io.Writer {
	mp.extra = append(mp.extra, new(strings.Builder))
	return mp.extra[len(mp.extra)-1]
}

func (mp *ModelsPlugger) GetExtensionType() protoreflect.ExtensionType {
	return generated.E_DbOpt
}

func (mp *ModelsPlugger) GetMessageName() string {
	return mp.message.GoIdent.GoName
}

func (mp *ModelsPlugger) AsTmplPack() *ModelsPluggerPack {
	imports := make([]string, len(mp.imports))
	for i, imp := range mp.imports {
		imports[i] = imp.String()
	}
	fields := make([]string, len(mp.fields))
	fieldsIndex := 0
	for _, field := range mp.fields {
		value := field.GetContent()
		fields[fieldsIndex] = value
	}
	methods := make([]string, len(mp.methods))
	for index, method := range mp.methods {
		methods[index] = method.String()
	}
	return &ModelsPluggerPack{
		PackageName: utils.GetGoPackageName(mp.parent.GoImportPath),
		Imports:     imports,
		ModelName:   mp.message.GoIdent.GoName,
		Fields:      fields,
		Methods:     methods,
	}
}

func (mp *ModelsPlugger) GetTemplateName() string {
	return "models.tmpl"
}

func (mp *ModelsPlugger) Init(template *template.Template) {
	mp.tmpl = template
}

func (mp *ModelsPlugger) GenerateFile() string {
	return fmt.Sprintf("%s.cato.go", strings.ToLower(mp.message.GoIdent.GoName))
}

func (mp *ModelsPlugger) GenerateContent() string {
	sw := new(strings.Builder)
	err := mp.tmpl.Execute(sw, mp.AsTmplPack())
	if err != nil {
		log.Fatalln("[-] models plugger exec tmpl error, ", err)
	}
	return sw.String()
}
