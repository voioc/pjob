package service

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/worker"
)

type TaskLogService struct {
	common.Base
}

// TaskS instance
func TaskLogS(c *gin.Context) *TaskLogService {
	return &TaskLogService{Base: common.Base{C: c}}
}

// func (s *TaskLogService) LogList(page, pageSize int, filters ...interface{}) ([]*model.TaskLog, int64, error) {
// 	offset := (page - 1) * pageSize
// 	data := make([]*model.TaskLog, 0)

// 	// query := model.GetDB()
// 	// var count int
// 	condition := " 1 = 1 "
// 	if len(filters) > 0 {
// 		for k := 0; k < len(filters); k += 2 {
// 			condition = fmt.Sprintf("%s and %s %v", condition, filters[k].(string), filters[k+1])
// 		}
// 	}

// 	total, err := model.GetDB().Where(condition).Count(&model.TaskLog{})
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

	type dc struct {
		Days  string
		Count int
	}

	tmp := make([]dc, 0)
	// key := m[status]

	// if RunNumCache.IsExist(key) {
	// 	json.Unmarshal(RunNumCache.Get(key).([]byte), &res)
	// 	logs.Info("cache")
	// 	return res
	// }
	if err := model.GetDB().SQL("SELECT FROM_UNIXTIME(create_time,'%Y-%m-%d') days,COUNT(id) count FROM pp_task_log WHERE status in(?) GROUP BY days ORDER BY days DESC limit ?;",
		status, limit).Find(&tmp); err != nil {
		return nil, err
	}

	data := map[string]int{}
	for _, row := range tmp {
		data[row.Days] = row.Count
	}

	// data, err := json.Marshal(res)
	// if err != nil {
	// 	return nil
	// }
	// RunNumCache.Put(key, data, 2*time.Hour)

	return data, nil

}

// func (s *TaskLogService) LogByID(id []int) (map[int]*model.TaskLog, error) {
// 	logs := make([]*model.TaskLog, 0)

// 	if err := model.GetDB().In("id", id).Find(&logs); err != nil {
// 		return nil, err
// 	}

// 	data := map[int]*model.TaskLog{}
// 	for _, row := range logs {
// 		data[row.ID] = row
// 	}

// 	return data, nil
// }

// func (s *TaskLogService) LogDelID(ids interface{}) error {
// 	_, flag1 := ids.([]int)
// 	_, flag2 := ids.([]string)

// 	if flag1 || flag2 {
// 		if _, err := model.GetDB().In("id", ids).Delete(&model.TaskLog{}); err != nil {
// 			return err
// 		}
// 	}

// 	return fmt.Errorf("record not found")
// }

func (s *TaskLogService) LogDelTaskID(ids interface{}) error {
	_, flag1 := ids.([]int)
	_, flag2 := ids.([]string)

	if flag1 || flag2 {
		if _, err := model.GetDB().In("task_id", ids).Delete(&model.TaskLog{}); err != nil {
			return err
		}
	}

	return nil
}

func (s *TaskLogService) TaskLogFunc(job *worker.Job, result *worker.JobResult) int {
	log := model.TaskLog{
		TaskID:      job.TaskID,
		ServerID:    job.ServerID,
		ServerName:  job.ServerName,
		Output:      result.OutMsg,
		Error:       result.ErrMsg,
		ProcessTime: int(time.Since(job.StartAt) / time.Millisecond),
		CreatedAt:   job.StartAt.Unix(),
	}

	timeout := time.Duration(time.Hour * 24)
	if job.Timeout > 0 {
		timeout = time.Second * time.Duration(job.Timeout)
	}

	if result.IsTimeout {
		log.Status = model.TASK_TIMEOUT
		log.Error = fmt.Sprintf("任务执行超过 %d 秒\n----------------------\n%s\n", int(timeout/time.Second), result.ErrMsg)
	} else if !result.IsOk {
		log.Status = model.TASK_ERROR
		log.Error = "ERROR:" + result.ErrMsg
	}

	if err := model.Add(log); err != nil {
		fmt.Println(err.Error())
	}
	return log.ID
}
