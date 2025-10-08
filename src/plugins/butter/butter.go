package butter

import (
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/src/plugins/common"
)

type Butter interface {
	FromExtType() protoreflect.ExtensionType
	WorkOn(desc protoreflect.Descriptor) bool
	Init(value interface{})
	Register(ctx *common.GenContext) error
}

var (
	factory     []func() Butter
	factoryOnce = new(sync.Once)
)

func Register(builder func() Butter) {
	factoryOnce.Do(func() {
		factory = make([]func() Butter, 0)
	})
	factory = append(factory, builder)
}

func ChooseButter(desc protoreflect.Descriptor) []Butter {
	chosen := make([]Butter, 0)
	for index := range factory {
		b := factory[index]()
		if b.WorkOn(desc) {
			chosen = append(chosen, b)
		}
	}
	return chosen
}
