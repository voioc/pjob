package model

import (
	"github.com/voioc/coco/db"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

func GetDB() *xorm.EngineGroup {
	engine := db.GetMySQL()
	tbMapper := names.NewPrefixMapper(names.SnakeMapper{}, "pp_")
	engine.SetTableMapper(tbMapper)
	return engine
}
