package structs

import (
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/utils"
)

func init() {
	register(func() common.Butter {
		return new(FieldTagButter)
	})
}

type FieldTagButter struct {
	option *generated.StructOption
}

func (f *FieldTagButter) FromExtType() protoreflect.ExtensionType {
	return generated.E_StructOpt
}

func (f *FieldTagButter) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.MessageDescriptor)
	return ok
}

func (f *FieldTagButter) Init(value interface{}) {
	f.option = value.(*generated.StructOption)
}

func (f *FieldTagButter) Register(ctx *common.GenContext) error {
	if f.option == nil || len(f.option.GetFieldDefaultTags()) == 0 {
		return nil
	}
	tags := f.option.GetFieldDefaultTags()
	result := make([]*common.Tag, len(tags))
	// common tags will be load in message-work-on butter
	// so when load field, default tags will be loaded
	for index := range tags {
		result[index] = &common.Tag{
			KV: &common.Kv{
				Key:   tags[index].GetTagName(),
				Value: tags[index].GetTagValue(),
			},
			Mapper: utils.GetWordMapper(tags[index].Mapper),
		}
	}
	ctx.SetScopeTags(result)
	return nil
}
