package plugins

import (
	"fmt"
	"log"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/butter"
	"github.com/ncuhome/cato/src/plugins/cheese"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/models/packs"
	"github.com/ncuhome/cato/src/plugins/utils"
)

const (
	repoModelImportAlias = "model"
	repoFetchOneFuncName = "FetchOne"
	repoFetchAllFuncName = "FetchAll"
)

type MessageWorker struct {
	message *protogen.Message
	gfs     []*models.GenerateFileDesc
}

func NewMessageWorker(msg *protogen.Message) *MessageWorker {
	mp := new(MessageWorker)
	mp.message = msg
	mp.gfs = make([]*models.GenerateFileDesc, 0)
	return mp
}

// RegisterContext because generate a file from a message, so a file-level writer for a message generates progress
func (mw *MessageWorker) RegisterContext(gc *common.GenContext) *common.GenContext {
	mc := cheese.NewMessageCheese()
	ctx := gc.WithMessage(mw.message, mc)
	return ctx
}

func (mw *MessageWorker) asBasicTmpl(ctx *common.GenContext) *packs.ModelContentTmplPack {
	mc := ctx.GetNowMessageContainer()
	return &packs.ModelContentTmplPack{
		ModelName: mw.message.GoIdent.GoName,
		Fields:    mc.GetField(),
		Methods:   mc.GetMethods(),
	}
}

func (mw *MessageWorker) GenerateFile() string {
	return fmt.Sprintf("%s.cato.go", mw.outputFileName())
}

func (mw *MessageWorker) outputFileName() string {
	patterns := utils.SplitCamelWords(mw.message.GoIdent.GoName)
	mapper := utils.GetStringsMapper(generated.FieldMapper_CATO_FIELD_MAPPER_SNAKE_CASE)
	return mapper(patterns)
}

func (mw *MessageWorker) GenerateContent(ctx *common.GenContext) string {
	sw := new(strings.Builder)
	tmpl := config.GetTemplate(config.ModelTmpl)
	err := tmpl.Execute(sw, mw.asBasicTmpl(ctx))
	if err != nil {
		log.Fatalln("[-] models plugger exec tmpl error, ", err)
	}
	return sw.String()
}

func (mw *MessageWorker) GenerateExtraContent(ctx *common.GenContext) string {
	sw := new(strings.Builder)
	tmpl := config.GetTemplate(config.TableExtendTmpl)
	err := tmpl.Execute(sw, mw.asBasicTmpl(ctx))
	if err != nil {
		log.Fatalln("[-] models plugger exec extend tmpl error, ", err)
	}
	return sw.String()
}

func (mw *MessageWorker) Active(ctx *common.GenContext) (bool, error) {
	descriptor := protodesc.ToDescriptorProto(mw.message.Desc)
	butters := butter.ChooseButter(mw.message.Desc)

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
	// for fields
	for _, field := range mw.message.Fields {
		fp := NewFieldCheese(field)
		fieldCtx := fp.RegisterContext(ctx)
		_, err := fp.Active(fieldCtx)
		if err != nil {
			return false, err
		}
		err = fp.Complete(fieldCtx)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (mw *MessageWorker) Complete(ctx *common.GenContext) error {
	err := mw.completeCols(ctx)
	if err != nil {
		return err
	}
	err = mw.completeKey(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (mw *MessageWorker) completeCols(ctx *common.GenContext) error {
	mc := ctx.GetNowMessageContainer()
	cols := mc.GetScopeCols()
	if len(cols) == 0 {
		return nil
	}
	tmpl := config.GetTemplate(config.TableColTmpl)
	pack := &packs.TableColTmplPack{
		MessageTypeName: ctx.GetNowMessageTypeName(),
		Cols:            cols,
	}
	return tmpl.Execute(mc.BorrowMethodsWriter(), pack)
}

func (mw *MessageWorker) completeKey(ctx *common.GenContext) error {
	keys := ctx.GetNowMessageContainer().GetScopeKeys()
	if len(keys) == 0 {
		return nil
	}
	mc := ctx.GetNowMessageContainer()
	for _, key := range keys {
		keyType := key.KeyType
		pack := &packs.RepoFuncTmplPack{
			KeyNameCombine:    key.GetFieldNameCombine(),
			ParamRaw:          strings.Join(key.GetParamsRaw(), ", "),
			ModelType:         fmt.Sprintf("%s.%s", repoModelImportAlias, ctx.GetNowMessageTypeName()),
			ModelPackage:      ctx.GetCatoPackage(),
			ModelPackageAlias: repoModelImportAlias,
		}
		packParams := make([]*packs.RepoFuncTmplPackParam, len(key.Fields))
		for index := range key.Fields {
			packParams[index] = &packs.RepoFuncTmplPackParam{
				FieldName: key.Fields[index].Name,
				ParamName: key.Fields[index].AsParamName(),
			}
		}
		pack.Params = packParams
		repoTmpls := []string{config.RepoFetchTmpl, config.RepoInsertTmpl}
		dbTmpls := []string{config.RdbFetchTmpl, config.RdbInsertTmpl}
		switch keyType {
		// unique and primary key will have FetchOne, UpdateBy, DeleteByMethod
		case generated.DBKeyType_CATO_DB_KEY_TYPE_PRIMARY, generated.DBKeyType_CATO_DB_KEY_TYPE_UNIQUE:
			pack.FetchReturnType = fmt.Sprintf("*%s", pack.ModelType)
			pack.FetchFuncName = repoFetchOneFuncName
			repoTmpls = append(repoTmpls, config.RepoUpdateTmpl, config.RepoDeleteTmpl)
			dbTmpls = append(dbTmpls, config.RdbUpdateTmpl, config.RdbDeleteTmpl)
		case generated.DBKeyType_CATO_DB_KEY_TYPE_COMBINE:
			pack.FetchReturnType = fmt.Sprintf("[]*%s", pack.ModelType)
			pack.FetchFuncName = repoFetchAllFuncName
		}
		// for fetch
		for _, repoTmpl := range repoTmpls {
			err := config.GetTemplate(repoTmpl).Execute(mc.BorrowRepoWriter(), pack)
			if err != nil {
				return err
			}
		}
		// for rdb
		for _, dbTmpl := range dbTmpls {
			err := config.GetTemplate(dbTmpl).Execute(mc.BorrowRdbWriter(), pack)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
