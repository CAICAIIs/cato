package defines

import (
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
)

type PackageButter struct {
	value *generated.CatoOptions
}

func (p *PackageButter) FromExtType() protoreflect.ExtensionType {
	return generated.E_CatoOpt
}

func (p *PackageButter) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.FileDescriptor)
	return ok
}

func (p *PackageButter) Init(value interface{}) {
	_, ok := value.(*generated.CatoOptions)
	if !ok {
		return
	}
	p.value = value.(*generated.CatoOptions)
}

func (p *PackageButter) Register(ctx *common.GenContext) error {
	fc := ctx.GetNowFileContainer()
	if p.value.GetCatoPackage() != "" {
		fc.SetCatoPackage(p.value.GetCatoPackage())
	}
	if p.value.GetRepoPackage() != "" {
		fc.SetRepoPackage(p.value.GetRepoPackage())
	}
	if p.value.GetRdbRepoPackage() != "" {
		fc.SetRdbRepoPackage(p.value.GetRdbRepoPackage())
	}
	return nil
}
