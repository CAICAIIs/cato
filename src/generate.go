package src

import (
	"go/format"
	"log"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/ncuhome/cato/src/plugins"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/flags"
)

type CatoGenerator struct {
	req    *pluginpb.CodeGeneratorRequest
	params map[string]string
}

func NewCatoGenerator(req *pluginpb.CodeGeneratorRequest) *CatoGenerator {
	g := &CatoGenerator{req: req}
	g.params = make(map[string]string)
	if g.req.GetParameter() != "" {
		g.params = flags.ParseProtoOptFlag(g.req.GetParameter())
	}
	return g
}

func (g *CatoGenerator) Generate(resp *pluginpb.CodeGeneratorResponse) *pluginpb.CodeGeneratorResponse {
	genOption, err := protogen.Options{}.New(g.req)
	if err != nil {
		log.Fatalln(err)
	}
	for _, file := range genOption.Files {
		fc := plugins.NewFileCheese(file)
		ctx := fc.RegisterContext(new(common.GenContext))
		if ctx.GetCatoPackage() == "" {
			continue
		}
		for _, message := range file.Messages {
			// init message scope cheese
			mc := plugins.NewMessageWorker(message)
			mctx := mc.RegisterContext(ctx)
			ok, err := mc.Active(mctx)
			if err != nil || !ok {
				log.Fatalf("[-] cato could not activate message %s: %v\n", mctx.GetNowMessageTypeName(), err)
			}
		}
	}
	return nil
}

func (g *CatoGenerator) outputContent(filename, content string) *pluginpb.CodeGeneratorResponse_File {
	c, err := format.Source([]byte(content))
	if err != nil {
		log.Fatalf("[-] cato formatted %s file content error\n", filename)
	}
	formatted := string(c)
	return &pluginpb.CodeGeneratorResponse_File{
		Name:    &filename,
		Content: &formatted,
	}
}
