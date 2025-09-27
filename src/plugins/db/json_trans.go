package db

import (
	"fmt"
	"text/template"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
)

func init() {
	register(func() common.Butter {
		return new(JsonTransButter)
	})
}

type JsonTransButter struct {
	value *generated.ColumnOption
	tmpl  *template.Template
}

type JsonTransButterPack struct {
	MessageTypeName string
	FieldName       string
	FieldType       string
	FieldTypeRaw    string
	LazyLoad        bool
}

func (j *JsonTransButter) FromExtType() protoreflect.ExtensionType {
	return generated.E_ColumnOpt
}

func (j *JsonTransButter) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.FieldDescriptor)
	return ok
}

func (j *JsonTransButter) tmplName() string {
	return "json_trans.tmpl"
}

func (j *JsonTransButter) Init(value interface{}) {
	data, ok := value.(*generated.ColumnOption)
	if !ok {
		return
	}
	j.value = data
	j.tmpl = config.GetTemplate(j.tmplName())
}

func (j *JsonTransButter) AsTmplPack(ctx *common.GenContext) interface{} {
	nowField := ctx.GetNowField()
	fieldType := common.MapperGoTypeName(ctx, nowField.Desc)
	return &JsonTransButterPack{
		MessageTypeName: ctx.GetNowMessageTypeName(),
		FieldName:       nowField.GoName,
		FieldType:       fieldType,
		FieldTypeRaw:    common.UnwrapPointType(fieldType),
		LazyLoad:        j.value.JsonTrans.LazyLoad,
	}
}

func (j *JsonTransButter) Register(ctx *common.GenContext) error {
	if j.value == nil || j.value.GetJsonTrans() == nil {
		return nil
	}
	transOpt := j.value.GetJsonTrans()
	writers := ctx.GetWriters()
	nowField := ctx.GetNowField()
	if transOpt.LazyLoad {
		// need to register extra inner field into message-fields map
		extraField := &common.FieldPack{
			Name:   fmt.Sprintf("inner%s", nowField.GoName),
			GoType: common.MapperGoTypeName(ctx, nowField.Desc),
		}
		err := config.GetTemplate(config.CommonFieldTmpl).Execute(writers.FieldWriter(), extraField)
		if err != nil {
			return err
		}
	}
	_, err := writers.ImportWriter().Write([]byte("\"encoding/json\""))
	if err != nil {
		return err
	}
	packData := j.AsTmplPack(ctx)
	return config.GetTemplate(j.tmplName()).Execute(writers.MethodWriter(), packData)
}
