package db

import (
	"text/template"

	"github.com/ncuhome/cato/generated"
)

type TableMessageEx struct {
	message *ModelsPlugger

	value *generated.TableOption

	tmpl *template.Template
}

type TableMessageExTmplPack struct {
	MessageTypeName string
	TableName       string
	Comment         string
}

func (t *TableMessageEx) GetTmplFileName() string {
	return "table_name.tmpl"
}

func (t *TableMessageEx) Init(tmpl *template.Template, value *generated.TableOption) {
	t.tmpl = tmpl
	t.value = value
}

func (t *TableMessageEx) LoadPlugger(message *ModelsPlugger) {
	t.message = message
}

func (t *TableMessageEx) AsTmplPack() interface{} {
	nameOpt := t.value.GetNameOption()
	if nameOpt == nil || nameOpt.GetLazyName() || nameOpt.GetSimpleName() == "" {
		return nil
	}
	return &TableMessageExTmplPack{
		MessageTypeName: t.message.GetMessageName(),
		TableName:       nameOpt.GetSimpleName(),
		Comment:         t.value.GetComment(),
	}
}

func (t *TableMessageEx) Register() error {
	if t.value == nil || t.message == nil {
		return nil
	}
	pack := &TableMessageExTmplPack{
		MessageTypeName: t.message.GetMessageName(),
		Comment:         t.value.GetComment(),
	}
	// check if the table name is simple
	if t.value.NameOption.GetSimpleName() != "" {
		pack.TableName = t.value.NameOption.GetSimpleName()
		return t.tmpl.Execute(t.message.BorrowMethodsWriter(), pack)
	}
	// empty table name will impl in an extra file
	return t.tmpl.Execute(t.message.BorrowExtraWriter(), pack)
}
