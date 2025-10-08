package plugins

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/ncuhome/cato/src/plugins/cheese"
	"github.com/ncuhome/cato/src/plugins/common"
)

type FileWorker struct {
	file    *protogen.File
	context *common.GenContext
}

func NewFileCheese(file *protogen.File) *FileWorker {
	fc := new(FileWorker)
	fc.file = file
	return fc
}

func (fc *FileWorker) RegisterContext(gc *common.GenContext) *common.GenContext {
	f := cheese.NewFileCheese(fc.file)
	ctx := gc.WithFile(fc.file, f)
	return ctx
}
