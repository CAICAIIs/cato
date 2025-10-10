package plugins

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"

	"github.com/ncuhome/cato/src/plugins/butter"
	"github.com/ncuhome/cato/src/plugins/cheese"
	"github.com/ncuhome/cato/src/plugins/common"
)

type FileWorker struct {
	file    *protogen.File
	context *common.GenContext
}

func NewFileWorker(file *protogen.File) *FileWorker {
	fc := new(FileWorker)
	fc.file = file
	return fc
}

func (fc *FileWorker) RegisterContext(gc *common.GenContext) *common.GenContext {
	f := cheese.NewFileCheese(fc.file)
	ctx := gc.WithFile(fc.file, f)
	return ctx
}

func (fc *FileWorker) Active(ctx *common.GenContext) (bool, error) {
	descriptor := protodesc.ToFileDescriptorProto(fc.file.Desc)
	butters := butter.ChooseButter(fc.file.Desc)
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
	return true, nil
}

func (fc *FileWorker) Complete(_ *common.GenContext) error {
	return nil
}
