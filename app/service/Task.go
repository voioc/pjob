package service

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/common"
)

type TaskService struct {
	common.Base
}

// TaskS instance
func TaskS(c *gin.Context) *TaskService {
	return &TaskService{Base: common.Base{C: c}}
}

func (s *TaskService) TaskGetList(page, pageSize int, filters ...interface{}) ([]*model.Task, int64, error) {
	offset := (page - 1) * pageSize
	tasks := make([]*model.Task, 0)

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

	total, err := db.Where(condition).Count(&model.Task{})
	if err != nil {
		return nil, 0, err
	}

	if err := model.GetDB().Where(condition).OrderBy("field(status, 1, 2, 3, 0), id desc ").Limit(pageSize, offset).Find(&tasks); err != nil {
		return nil, 0, err
	}

	// query := orm.NewOrm().QueryTable(TableName("task"))
	// if len(filters) > 0 {
	// 	l := len(filters)
	// 	for k := 0; k < l; k += 2 {
	// 		query = query.Filter(filters[k].(string), filters[k+1])
	// 	}
	// }

	return tasks, total, nil
}

// 运行总次数
func (s *TaskService) TaskTotalRunNum() (int64, error) {

	// res := make(orm.Params)
	// _, err := orm.NewOrm().Raw("select sum(execute_times) as num,task_name from pp_task").RowsToMap(&res, "num", "task_name")

	return model.GetDB().SumInt(&model.Task{}, "execute_times")
}

func (s *TaskService) TaskByID(id int) (*model.Task, error) {
	task := &model.Task{}

	if _, err := model.GetDB().Where("id = ?", id).Get(task); err != nil {
		return nil, err
	}

	if task.ID == 0 {
		return nil, fmt.Errorf("task not found")
	}

	return task, nil
}

func (s *TaskService) Add(task *model.Task) (int, error) {
	_, err := model.GetDB().Insert(task)
	return task.ID, err
}

func (s *TaskService) Update(task *model.Task, args ...bool) error {
	if len(args) > 0 && args[0] {
		if _, err := model.GetDB().Cols("status").Where("id = ?", task.ID).Update(task); err != nil {
			return err
		}
	} else {
		if _, err := model.GetDB().Where("id = ?", task.ID).Update(task); err != nil {
			return err
		}
	}

	return nil
}

func (s *TaskService) Del(ids interface{}) error {
	_, flag1 := ids.([]int)
	_, flag2 := ids.([]string)

	if flag1 || flag2 {
		_, err := model.GetDB().In("id", ids).Delete(&model.Task{})
		return err
	}

	return nil
}
