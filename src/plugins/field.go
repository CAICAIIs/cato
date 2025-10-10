package plugins

import (
	"fmt"
	"io"
	"log"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/src/plugins/butter"
	"github.com/ncuhome/cato/src/plugins/cheese"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/models/packs"
	"github.com/ncuhome/cato/src/plugins/utils"
)

type FieldWorker struct {
	field       *protogen.Field
	tags        []*strings.Builder
	DefaultTags []*models.Kv
}

func NewFieldCheese(field *protogen.Field) *FieldWorker {
	return &FieldWorker{
		field: field,
		tags:  make([]*strings.Builder, 0),
	}
}

func (fw *FieldWorker) RegisterContext(gc *common.GenContext) *common.GenContext {
	fc := cheese.NewFieldCheese()
	ctx := gc.WithField(fw.field, fc)
	return ctx
}

func (fw *FieldWorker) borrowTagWriter() io.Writer {
	fw.tags = append(fw.tags, new(strings.Builder))
	return fw.tags[len(fw.tags)-1]
}

func (fw *FieldWorker) AsTmplPack(fieldType string) *packs.FieldPack {
	pack := &packs.FieldPack{
		Field: &models.Field{
			Name:   fw.field.GoName,
			GoType: fieldType,
		},
	}
	tags := make([]string, len(fw.tags))
	tagMap := make(map[string]struct{})
	for index := range fw.tags {
		raw := fw.tags[index].String()
		tagKey := utils.GetTagKey(raw)
		_, hasTag := tagMap[tagKey]
		if tagKey == "" || hasTag {
			continue
		}
		tags[index] = fw.tags[index].String()
		tagMap[tagKey] = struct{}{}
	}
	pack.Tags = strings.Join(tags, " ")
	return pack
}

func (fw *FieldWorker) Active(ctx *common.GenContext) (bool, error) {
	butters := butter.ChooseButter(fw.field.Desc)
	descriptor := protodesc.ToFieldDescriptorProto(fw.field.Desc)
	for index := range butters {
		if !proto.HasExtension(descriptor.Options, butters[index].FromExtType()) {
			continue
		}
		value := proto.GetExtension(descriptor.Options, butters[index].FromExtType())
		butters[index].Init(value)
		err := butters[index].Register(ctx)
		if err != nil {
			return false, err
		}
	}
	// need register tags in ctx
	for _, scopeTag := range ctx.GetNowMessageContainer().GetScopeTags() {
		if scopeTag.KV == nil {
			continue
		}
		target := fw.borrowTagWriter()
		tagData := fmt.Sprintf("%s:\"%s\"", scopeTag.KV.Key, scopeTag.GetTagValue(fw.field.GoName))
		_, err := target.Write([]byte(tagData))
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (fw *FieldWorker) Complete(ctx *common.GenContext) error {
	wr := ctx.GetNowMessageContainer().BorrowFieldWriter()
	// register into field writer
	fieldType := common.MapperGoTypeName(ctx, fw.field.Desc)
	if ctx.GetNowFieldContainer().IsJsonTrans() {
		fieldType = "string"
		mc := ctx.GetNowMessageContainer()
		mc.SetScopeColType(fw.field.GoName, fieldType)
	}
	pack := fw.AsTmplPack(fieldType)
	return config.GetTemplate(config.FieldTmpl).Execute(wr, pack)
}
