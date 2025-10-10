package cheese

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/utils"
)

type FileCheese struct {
	imports map[string]*models.Import
	// todo optimize as repo map
	catoPackage    *models.Import
	catoExtPackage *models.Import
	repoPackage    *models.Import
	rdbRepoPackage *models.Import
}

func NewFileCheese(file *protogen.File) *FileCheese {
	cheese := new(FileCheese)
	cheese.imports = make(map[string]*models.Import)

	desc := file.Desc
	for index := 0; index < desc.Imports().Len(); index++ {
		importFile := desc.Imports().Get(index)
		importPackage := string(importFile.FileDescriptor.Package())
		importCatoPath, ok := utils.GetCatoPackageFromFile(importFile.FileDescriptor)
		if !ok {
			continue
		}
		cheese.imports[importPackage] = new(models.Import).Init(importCatoPath)
	}
	return cheese
}

func (fc *FileCheese) GetImportPathAlias(path string) string {
	v, ok := fc.imports[path]
	if !ok {
		return ""
	}
	return v.Alias
}

func (fc *FileCheese) GetImports() []string {
	imports, index := make([]string, len(fc.imports)), 0
	for _, v := range fc.imports {
		imports[index] = v.GetPath()
		index++
	}
	return imports
}

func (fc *FileCheese) SetCatoPackage(packagePath string) {
	i := new(models.Import).Init(packagePath)
	fc.catoPackage = i
}

func (fc *FileCheese) GetCatoPackage() *models.Import {
	return fc.catoPackage
}

func (fc *FileCheese) SetCatoExtPackage(packagePath string) {
	i := new(models.Import).Init(packagePath)
	fc.catoExtPackage = i
}

func (fc *FileCheese) GetCatoExtPackage() *models.Import {
	return fc.catoExtPackage
}

func (fc *FileCheese) SetRepoPackage(packagePath string) {
	i := new(models.Import).Init(packagePath)
	fc.repoPackage = i
}

func (fc *FileCheese) GetRepoPackage() *models.Import {
	return fc.repoPackage
}

func (fc *FileCheese) SetRdbRepoPackage(packagePath string) {
	i := new(models.Import).Init(packagePath)
	fc.rdbRepoPackage = i
}

func (fc *FileCheese) GetRdbRepoPackage() *models.Import {
	return fc.rdbRepoPackage
}
