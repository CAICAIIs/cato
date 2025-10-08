package packs

import (
	"github.com/ncuhome/cato/src/plugins/models"
)

type TableColTmplPack struct {
	MessageTypeName string
	Cols            []*models.Col
}
