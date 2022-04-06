package model

import (
	"fmt"
	"strings"

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

// dataModel 指针类型数据结构
func List(dataModel interface{}, page, pageSize int, filters ...interface{}) error {
	offset := (page - 1) * pageSize
	// data := make([]*model.Role, 0)

	in := map[string]interface{}{}
	order := "id asc"
	condition := " 1 = 1 "
	if len(filters) > 0 {
		for k := 0; k < len(filters); k += 2 {
			// 如果是数组则单独筛出来
			if _, flag := filters[k+1].([]int); flag {
				in[filters[k].(string)] = filters[k+1]
			} else if strings.Trim(filters[k].(string), "") == "order" {
				order = filters[k+1].(string)
			} else {
				condition = fmt.Sprintf("%s and %s %v", condition, filters[k].(string), filters[k+1])
			}
		}
	}

	db := GetDB().Where("1=1")
	if len(in) > 0 {
		for col, v := range in {
			if col != "" {
				regex := strings.Split(col, " ")
				if len(regex) == 2 && regex[1] == "not" {
					db = db.NotIn(col, v)
				} else {
					db = db.In(col, v)
				}
			}
		}
	}

	if err := db.Where(condition).OrderBy(order).Limit(pageSize, offset).Find(dataModel); err != nil {
		return err
	}

	// if len(dataModel) < 1 {
	// 	return fmt.Errorf("record not found")
	// }

	return nil
}

// dataModel 指针类型数据结构
func ListCount(dataModel interface{}, filters ...interface{}) (int64, error) {
	// offset := (page - 1) * pageSize
	// data := make([]*model.Role, 0)

	in := map[string]interface{}{}
	condition := " 1 = 1 "
	if len(filters) > 0 {
		for k := 0; k < len(filters); k += 2 {
			// 如果是数组则单独筛出来
			if _, flag := filters[k+1].([]int); flag {
				in[filters[k].(string)] = filters[k+1]
			} else if strings.Trim(filters[k].(string), "") == "order" {
				continue
			} else {
				condition = fmt.Sprintf("%s and %s %v", condition, filters[k].(string), filters[k+1])
			}
		}
	}

	db := GetDB().Where("1=1")
	if len(in) > 0 {
		for col, v := range in {
			if col != "" {
				regex := strings.Split(col, " ")
				if len(regex) == 2 && regex[1] == "not" {
					db = db.NotIn(col, v)
				} else {
					db = db.In(col, v)
				}
			}
		}
	}

	return db.Where(condition).Count(dataModel)
}

func DataByID(m interface{}, id int) error {
	if _, err := GetDB().ID(id).Get(m); err != nil {
		return err
	}

	return nil
}

// m 数据类型 ids []int 或者 []string
func DataByIDs(m interface{}, ids interface{}, col ...string) error {
	colum := "id"
	if len(col) > 0 {
		colum = col[0]
	}

	if err := GetDB().In(colum, ids).Find(m); err != nil {
		return err
	}

	return nil
}

func Add(data interface{}) error {
	_, err := GetDB().Insert(data)
	return err
}

func Update(id int, data interface{}, args ...bool) error {
	if len(args) > 0 && args[0] { // 状态删除
		if _, err := GetDB().Cols("status").ID(id).Update(data); err != nil {
			return err
		}
	} else {
		if _, err := GetDB().ID(id).Update(data); err != nil {
			return err
		}
	}

	return nil
}

func Del(dataModel interface{}, ids interface{}) error {
	_, flag1 := ids.([]int)
	_, flag2 := ids.([]string)

	if flag1 || flag2 {
		_, err := GetDB().In("id", ids).Delete(dataModel)
		return err
	}

	return nil
}
