package src

import (
	"log"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/db"
)

type DbGenerator struct {
	req  *pluginpb.CodeGeneratorRequest
	resp *pluginpb.CodeGeneratorResponse
}

func NewDBGenerator(req *pluginpb.CodeGeneratorRequest) *DbGenerator {
	return &DbGenerator{req: req}
}

func (g *DbGenerator) Generate(resp *pluginpb.CodeGeneratorResponse) *pluginpb.CodeGeneratorResponse {
	genOption, err := protogen.Options{}.New(g.req)
	if err != nil {
		log.Fatalln(err)
	}
	for _, file := range genOption.Files {
		for _, message := range file.Messages {
			mp := new(db.ModelsPlugger)
			mp.Init(config.GetTemplate(mp.GetTemplateName()))
			mp.LoadContext(message)
			descriptor := g.messageToDescriptorProto(message)
			if proto.HasExtension(descriptor, generated.E_TableOpt) {
				tableExt := new(db.TableMessageEx)
				tableExt.Init(config.GetTemplate(tableExt.GetTmplFileName()))
				tableExt.LoadPlugger(mp)
				err = tableExt.Register()
				if err != nil {
					log.Fatalln(err)
				}
			}
			fileName := mp.GenerateFile()
			content := mp.GenerateContent()
			resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
				Name:    &fileName,
				Content: &content,
			})
		}
	}
	return nil
}

func (g *DbGenerator) messageToDescriptorProto(msg *protogen.Message) *descriptorpb.DescriptorProto {
	if msg.Desc != nil {
		return msg.Desc.(interface{ ProtoReflect() protoreflect.Message }).
			ProtoReflect().Interface().(*descriptorpb.DescriptorProto)
	}
	return nil
}
