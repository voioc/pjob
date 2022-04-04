package service

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/common"
)

type NotifyService struct {
	common.Base
}

// NotifyS instance
func NotifyS(c *gin.Context) *NotifyService {
	return &NotifyService{Base: common.Base{C: c}}
}

func (s *NotifyService) NotifyTypeList(ntype int) ([]*model.NotifyTpl, error) {
	data := make([]*model.NotifyTpl, 0)
	if err := model.GetDB().Where("tpl_type = ? and status = 1", ntype).Find(&data); err != nil {
		return nil, err

	}
	return data, nil
}

func (s *NotifyService) NotifyListIDs(ids []int) ([]*model.NotifyTpl, error) {
	data := make([]*model.NotifyTpl, 0)
	if err := model.GetDB().Where("status = 1").In("id", ids).Find(&data); err != nil {
		return nil, err

	}
	return data, nil
}

func (s *NotifyService) NotifyList(page, pageSize int, filters ...interface{}) ([]*model.NotifyTpl, int64, error) {
	offset := (page - 1) * pageSize
	data := make([]*model.NotifyTpl, 0)

	db := model.GetDB().Where("1=1")

	in := map[string]interface{}{}
	condition := " 1 = 1 "
	if len(filters) > 0 {
		for k := 0; k < len(filters); k += 2 {
			// 如果是数组则单独筛出来
			if _, flag := filters[k+1].([]int); flag {
				in[filters[k].(string)] = filters[k+1]
			} else {
				condition = fmt.Sprintf("%s and %s %v", condition, filters[k].(string), filters[k+1])
			}
		}
	}

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

	total, err := db.Where(condition).Count(&model.NotifyTpl{})
	if err != nil {
		return nil, 0, err
	}

	if err := db.Where(condition).Limit(pageSize, offset).Find(&data); err != nil {
		return nil, 0, err
	}

	// query := orm.NewOrm().QueryTable(TableName("task"))
	// if len(filters) > 0 {
	// 	l := len(filters)
	// 	for k := 0; k < l; k += 2 {
	// 		query = query.Filter(filters[k].(string), filters[k+1])
	// 	}
	// }

	return data, total, nil
}

// func (s *NotifyService) NotifyByID(id int) (*model.NotifyTpl, error) {
// 	data := &model.NotifyTpl{}

// 	if _, err := model.GetDB().Where("id = ?", id).Get(data); err != nil {
// 		return nil, err
// 	}

// 	if data.ID == 0 {
// 		return nil, fmt.Errorf("record not found")
// 	}

// 	return data, nil
// }

// func (s *NotifyService) Add(data *model.NotifyTpl) (int, error) {
// 	_, err := model.GetDB().Insert(data)
// 	return data.ID, err
// }

// func (s *NotifyService) Update(data *model.NotifyTpl, args ...bool) error {
// 	if len(args) > 0 && args[0] {
// 		if _, err := model.GetDB().Cols("status").Where("id = ?", data.ID).Update(data); err != nil {
// 			return err
// 		}
// 	} else {
// 		if _, err := model.GetDB().Where("id = ?", data.ID).Update(data); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (s *NotifyService) Del(ids interface{}) error {
// 	_, flag1 := ids.([]int)
// 	_, flag2 := ids.([]string)

// 	if flag1 || flag2 {
// 		_, err := model.GetDB().In("id", ids).Delete(&model.NotifyTpl{})
// 		return err
// 	}

// 	return nil
// }
