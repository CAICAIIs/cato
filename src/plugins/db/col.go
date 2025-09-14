package db

import (
	"errors"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"text/template"
)

type ColumnFieldEx struct {
	field  *FieldsPlugger
	parent *ModelsPlugger

	value *generated.ColumnOption

	tmpl *template.Template
	tags map[string]*common.Kv
}

type ColumnFieldExTmlPack struct {
	Name   string
	GoType string
	Tags   []common.Kv
}

func (c *ColumnFieldEx) GetTmplFileName() string {
	return "column_field.tmpl"
}

func (c *ColumnFieldEx) Init(tmpl *template.Template) {
	c.tmpl = tmpl
}

func (c *ColumnFieldEx) LoadPlugger(field *FieldsPlugger, message *ModelsPlugger) {
	c.field = field
	c.parent = message
}

func (c *ColumnFieldEx) AsTmplPack() interface{} {
	tags, index := make([]common.Kv, len(c.tags)), 0
	for k, v := range c.tags {
		tags[index] = common.Kv{
			Key:   k,
			Value: v.Value,
		}
		index++
	}
	return &ColumnFieldExTmlPack{
		Name:   c.field.GetName(),
		GoType: c.field.GetGoType(),
		Tags:   tags,
	}
}

func (c *ColumnFieldEx) Register() error {
	// self-tags has the highest priority
	selfTags := c.value.GetTags()
	wr, ok := c.parent.BorrowFieldsWriter(c.field.GetName())
	if !ok {
		return errors.New("could not create writer")
	}
	for _, tag := range selfTags {
		c.tags[tag.TagName] = &common.Kv{
			Key:   tag.TagName,
			Value: tag.TagValue,
		}
	}
	packData := c.AsTmplPack()
	return c.tmpl.Execute(wr, packData)
}
