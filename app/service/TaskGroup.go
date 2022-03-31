package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/common"
)

type TaskGroupService struct {
	common.Base
}

// TaskS instance
func TaskGroupS(c *gin.Context) *TaskGroupService {
	return &TaskGroupService{Base: common.Base{C: c}}
}

func (s *TaskGroupService) GroupList(page, pageSize int, filters ...interface{}) ([]*model.TaskGroup, int64, error) {
	offset := (page - 1) * pageSize
	data := make([]*model.TaskGroup, 0)

	// query := model.GetDB()
	// var count int
	condition := " 1 = 1 "
	if len(filters) > 0 {
		for k := 0; k < len(filters); k += 2 {
			condition = fmt.Sprintf("%s and %s %s", condition, filters[k].(string), filters[k+1])
		}
	}

	total, err := model.GetDB().Where(condition).Count(&model.TaskGroup{})
	if err != nil {
		return nil, 0, err
	}

	if err := model.GetDB().Where(condition).Limit(pageSize, offset).Find(&data); err != nil {
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
