package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/common"
)

type TaskLogService struct {
	common.Base
}

// TaskS instance
func TaskLogS(c *gin.Context) *TaskLogService {
	return &TaskLogService{Base: common.Base{C: c}}
}

func (s *TaskLogService) LogList(page, pageSize int, filters ...interface{}) ([]*model.TaskLog, int64, error) {
	offset := (page - 1) * pageSize
	data := make([]*model.TaskLog, 0)

	// query := model.GetDB()
	// var count int
	condition := " 1 = 1 "
	if len(filters) > 0 {
		for k := 0; k < len(filters); k += 2 {
			condition = fmt.Sprintf("%s and %s %s", condition, filters[k].(string), filters[k+1])
		}
	}

	total, err := model.GetDB().Where(condition).Count(&model.TaskLog{})
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

func (s *TaskLogService) GetLogNum(status int) (int64, error) {
	// return orm.NewOrm().QueryTable(TableName("task_log")).Filter("status", status).Count()

	return model.GetDB().Where("status = ?", status).Count(&model.TaskLog{})
}

func (s *TaskLogService) SumByDays(limit int, status string) (map[string]int, error) {

	// var m = map[string]string{
	// 	"0":  "okNum",
	// 	"-1": "errNum",
	// 	"-2": "expiredRun",
	// }

	res := make([]map[string]int, 0)
	// key := m[status]

	// if RunNumCache.IsExist(key) {
	// 	json.Unmarshal(RunNumCache.Get(key).([]byte), &res)
	// 	logs.Info("cache")
	// 	return res
	// }
	if err := model.GetDB().SQL("SELECT FROM_UNIXTIME(create_time,'%Y-%m-%d') days,COUNT(id) count FROM pp_task_log WHERE status in(?) GROUP BY days ORDER BY days DESC limit ?;",
		status, limit).Find(&res); err != nil {
		return nil, err
	}

	// data, err := json.Marshal(res)
	// if err != nil {
	// 	return nil
	// }
	// RunNumCache.Put(key, data, 2*time.Hour)
	return res[0], nil

}
