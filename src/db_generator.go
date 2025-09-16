package src

import (
	"log"
	"path/filepath"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/db"
	"github.com/ncuhome/cato/src/plugins/utils"
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
		goPackageName := utils.GetGoPackageName(file.GoImportPath)
		if goPackageName == "" {
			continue
		}
		for _, message := range file.Messages {
			mp := new(db.ModelsPlugger)
			mp.Init(config.GetTemplate(mp.GetTemplateName()))
			mp.LoadContext(message, file)
			descriptor := protodesc.ToDescriptorProto(message.Desc)
			if proto.HasExtension(descriptor.Options, generated.E_TableOpt) {
				tableExt := new(db.TableMessageEx)
				value := proto.GetExtension(descriptor.Options, generated.E_TableOpt).(*generated.TableOption)
				tableExt.Init(config.GetTemplate(tableExt.GetTmplFileName()), value)
				tableExt.LoadPlugger(mp)
				err = tableExt.Register()
				if err != nil {
					log.Fatalln(err)
				}
				fileName := filepath.Join(mp.GenerateFile())
				content := mp.GenerateContent()
				resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
					Name:    &fileName,
					Content: &content,
				})
			}
		}
	}
	return nil
}
