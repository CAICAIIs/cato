package db

import (
	"log"
	"text/template"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/utils"
)

func init() {
	register(func() common.Butter {
		return new(ColFieldButter)
	})
}

type ColFieldButter struct {
	value *generated.ColumnOption
	tmpl  *template.Template
	tags  map[string]string
}

type ColFieldButterTmplPack struct {
	*common.FieldPack
	Tags []common.Kv
}

func (c *ColFieldButter) Init(value interface{}) {
	exValue, ok := value.(*generated.ColumnOption)
	if !ok {
		log.Fatalln("[-] cato ColFieldButter except ColumnOption")
	}
	c.value = exValue
	c.tmpl = config.GetTemplate(c.tmplName())
	c.tags = make(map[string]string)
}

func (c *ColFieldButter) tmplName() string {
	return config.CommonTagTmpl
}

func (c *ColFieldButter) AsTmplPack(_ *common.GenContext) interface{} {
	tags, index := make([]common.Kv, len(c.tags)), 0
	for k, v := range c.tags {
		tags[index] = common.Kv{
			Key:   k,
			Value: v,
		}
		index++
	}
	return tags
}

func (c *ColFieldButter) FromExtType() protoreflect.ExtensionType {
	return generated.E_ColumnOpt
}

func (c *ColFieldButter) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.FieldDescriptor)
	return ok
}

func (c *ColFieldButter) Register(ctx *common.GenContext) error {
	// self-tags has the highest priority
	selfTags := c.value.GetTags()
	if len(selfTags) == 0 {
		return nil
	}
	for _, tag := range selfTags {
		t := &common.Tag{
			KV:     &common.Kv{Key: tag.TagName, Value: tag.TagValue},
			Mapper: utils.GetWordMapper(tag.Mapper),
		}
		c.tags[t.KV.Key] = t.GetTagValue(ctx.GetNowField().GoName)
	}
	writers := ctx.GetWriters()
	// check if the value has a json-trans option
	packData := c.AsTmplPack(ctx)
	return c.tmpl.Execute(writers.TagWriter(), packData)
}
