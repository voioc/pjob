package service

import (
	"fmt"
	"strings"

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

// func (s *TaskGroupService) GroupList(page, pageSize int, filters ...interface{}) ([]*model.TaskGroup, int64, error) {
// 	offset := (page - 1) * pageSize
// 	data := make([]*model.TaskGroup, 0)

// 	// query := model.GetDB()
// 	// var count int
// 	condition := " 1 = 1 "
// 	if len(filters) > 0 {
// 		for k := 0; k < len(filters); k += 2 {
// 			condition = fmt.Sprintf("%s and %s %v", condition, filters[k].(string), filters[k+1])
// 		}
// 	}

// 	total, err := model.GetDB().Where(condition).Count(&model.TaskGroup{})
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	if err := model.GetDB().Where(condition).Limit(pageSize, offset).Find(&data); err != nil {
// 		return nil, 0, err
// 	}

// 	// query := orm.NewOrm().QueryTable(TableName("task"))
// 	// if len(filters) > 0 {
// 	// 	l := len(filters)
// 	// 	for k := 0; k < l; k += 2 {
// 	// 		query = query.Filter(filters[k].(string), filters[k+1])
// 	// 	}
// 	// }

// 	return data, total, nil
// }

// func (s *TaskGroupService) Update(data *model.TaskGroup, args ...bool) error {
// 	if data.GroupName == "" {
// 		return fmt.Errorf("组名不能为空")
// 	}

// 	if len(args) > 0 && args[0] {
// 		if _, err := model.GetDB().ID(data.ID).Cols("status").Update(data); err != nil {
// 			return err
// 		}
// 	} else {
// 		if _, err := model.GetDB().Where("id = ?", data.ID).Update(data); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (s *TaskGroupService) GroupAdd(obj *model.TaskGroup) (int64, error) {
// 	if obj.GroupName == "" {
// 		return 0, fmt.Errorf("组名不能为空")
// 	}

// 	return model.GetDB().Insert(obj)
// }

// 根据任务组id获取对应的名字
func (s *TaskGroupService) GroupIDName(ids string) (map[int]string, error) {
	ids = strings.Trim(strings.Trim(ids, ","), "")
	gid := strings.Split(ids, ",")
	fmt.Println(gid)

	group := make([]*model.TaskGroup, 0)
	// err := model.GetDB().Where("status = 1").In("id", gid).Find(&group)
	err := model.GetDB().Where("status = 1").Find(&group)
	if err != nil {
		return nil, err
	}

	data := map[int]string{}
	for _, gv := range group {
		data[gv.ID] = gv.GroupName
	}
	return data, nil
}

// func (s *TaskGroupService) GroupByID(id int) (*model.TaskGroup, error) {
// 	obj := &model.TaskGroup{}
// 	if flag, err := model.GetDB().Where("id = ?", id).Get(obj); !flag || err != nil {
// 		if err == nil && !flag {
// 			err = fmt.Errorf("task group not found")
// 		}
// 		return obj, err
// 	}

// 	return obj, nil
// }

// func (s *TaskGroupService) GroupDel(ids []int) error {
// 	_, err := model.GetDB().In("id", ids).Delete(&model.TaskGroup{})
// 	return err
// }
