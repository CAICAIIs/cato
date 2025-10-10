package cheese

import (
	"io"
	"strings"
)

type FieldCheese struct {
	tags      []*strings.Builder
	jsonTrans bool
}

func NewFieldCheese() *FieldCheese {
	cheese := &FieldCheese{}
	cheese.tags = make([]*strings.Builder, 0)
	return cheese
}

func (fp *FieldCheese) BorrowTagWriter() io.Writer {
	fp.tags = append(fp.tags, new(strings.Builder))
	return fp.tags[len(fp.tags)-1]
}

func (fp *FieldCheese) GetTags() []string {
	data := make([]string, len(fp.tags))
	for i, tag := range fp.tags {
		data[i] = tag.String()
	}
	return data
}

func (fp *FieldCheese) SetJsonTrans(b bool) {
	fp.jsonTrans = b
}

func (fp *FieldCheese) IsJsonTrans() bool {
	return fp.jsonTrans
}
