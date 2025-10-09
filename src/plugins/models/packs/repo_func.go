package packs

type RepoFuncTmplPackParam struct {
	FieldName string
	ParamName string
}

type RepoFuncTmplPack struct {
	KeyNameCombine        string
	ModelType             string
	Params                []*RepoFuncTmplPackParam
	FetchFuncName         string
	FetchReturnType       string
	ModelPackage          string
	ModelPackageAlias     string
	IsModelAnotherPackage bool

	Tmpls       []string
	IsUniqueKey bool
}

func (pack *RepoFuncTmplPack) Copy() *RepoFuncTmplPack {
	params := make([]*RepoFuncTmplPackParam, len(pack.Params))
	for index := range pack.Params {
		params[index] = &RepoFuncTmplPackParam{
			FieldName: pack.Params[index].FieldName,
			ParamName: pack.Params[index].ParamName,
		}
	}
	return &RepoFuncTmplPack{
		KeyNameCombine:        pack.KeyNameCombine,
		ModelType:             pack.ModelType,
		Params:                params,
		FetchFuncName:         pack.FetchFuncName,
		FetchReturnType:       pack.FetchReturnType,
		ModelPackage:          pack.ModelPackage,
		ModelPackageAlias:     pack.ModelPackageAlias,
		IsModelAnotherPackage: pack.IsModelAnotherPackage,
		Tmpls:                 pack.Tmpls,
		IsUniqueKey:           pack.IsUniqueKey,
	}
}
