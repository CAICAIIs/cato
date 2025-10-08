package packs

type RepoFuncTmplPackParam struct {
	FieldName string
	ParamName string
}

type RepoFuncTmplPack struct {
	KeyNameCombine    string
	ParamRaw          string
	ModelType         string
	Params            []*RepoFuncTmplPackParam
	FetchFuncName     string
	FetchReturnType   string
	ModelPackage      string
	ModelPackageAlias string
}
