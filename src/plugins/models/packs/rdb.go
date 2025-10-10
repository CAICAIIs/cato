package packs

type RdbTmplPack struct {
	RdbRepoPackage        string
	IsModelAnotherPackage bool
	ModelPackageAlias     string
	ModelPackage          string
	RdbRepoFuncs          []string
	IsRepoAnotherPackage  bool
	RepoPackageAlias      string
	RepoPackage           string
	ModelType             string
	FetchOneReturnType    string
	FetchAllReturnType    string
}
