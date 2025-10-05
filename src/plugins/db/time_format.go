package db

import (
	"text/template"
	"time"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func init() {
	register(func() common.Butter {
		return new(TimeOptionButter)
	})
}

type TimeOptionButter struct {
	value      *generated.TimeOption
	timeFormat string
	tmpl       *template.Template
}

type TimeOptionButterTmplPack struct {
	MessageTypeName string
	FieldName       string
	Format          string
}

func (t *TimeOptionButter) FromExtType() protoreflect.ExtensionType {
	return generated.E_ColumnOpt
}

func (t *TimeOptionButter) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.FieldDescriptor)
	return ok
}

func (t *TimeOptionButter) tmplName() string {
	return "time_format.tmpl"
}

func (t *TimeOptionButter) Init(value interface{}) {
	colOpt, ok := value.(*generated.ColumnOption)
	if !ok {
		return
	}
	t.value = colOpt.TimeOption
	t.tmpl = config.GetTemplate(t.tmplName())
	t.timeFormat = time.RFC3339
}

func (t *TimeOptionButter) AsTmplPack(ctx *common.GenContext) interface{} {
	return &TimeOptionButterTmplPack{
		MessageTypeName: ctx.GetNowMessageTypeName(),
		FieldName:       ctx.GetNowField().GoName,
		Format:          t.timeFormat,
	}
}

func (t *TimeOptionButter) Register(ctx *common.GenContext) error {
	if t.value == nil {
		return nil
	}
	timeOpt := t.value
	if timeOpt.GetTimeFormat() != "" {
		t.timeFormat = timeOpt.GetTimeFormat()
	}
	writer := ctx.GetWriters()
	err := t.tmpl.Execute(writer.MethodWriter(), t.AsTmplPack(ctx))
	if err != nil {
		return err
	}
	_, err = writer.ImportWriter().Write([]byte("\"time\""))
	return err
}
