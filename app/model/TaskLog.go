/*
* @Author: haodaquan
* @Date:   2017-06-21 12:23:22
* @Last Modified by:   haodaquan
* @Last Modified time: 2017-06-22 14:57:13
 */

package model

import (
	"github.com/astaxie/beego/cache"
)

type TaskLog struct {
	ID          int    `xorm:"id pk" json:"id"`
	TaskID      int    `xorm:"task_id" json:"task_id"`
	ServerID    int    `xorm:"server_id" json:"server_id"`
	ServerName  string `xorm:"server_name" json:"server_name"`
	Output      string `xorm:"output" json:"output"`
	Error       string `xorm:"error" json:"error"`
	Status      int    `xorm:"status" json:"status"`
	ProcessTime int    `xorm:"process_time" json:"process_time"`
	CreatedAt   int64  `xorm:"create_time" json:"created_at"`
}

var RunNumCache, _ = cache.NewCache("memory", `{"interval":60}`)

func (t *TaskLog) TableName() string {
	return TableName("task_log")
}

// var TaskLogFunc = func(job *worker.Job, result *worker.JobResult) int {
// 	log := TaskLog{
// 		TaskID:      job.ID,
// 		ServerID:    job.ServerID,
// 		ServerName:  job.ServerName,
// 		Output:      result.OutMsg,
// 		Error:       result.ErrMsg,
// 		ProcessTime: int(time.Since(job.StartAt) / time.Millisecond),
// 		CreatedAt:   job.StartAt.Unix(),
// 	}

// 	timeout := time.Duration(time.Hour * 24)
// 	if job.Timeout > 0 {
// 		timeout = time.Second * time.Duration(job.Timeout)
// 	}

// 	if result.IsTimeout {
// 		log.Status = TASK_TIMEOUT
// 		log.Error = fmt.Sprintf("任务执行超过 %d 秒\n----------------------\n%s\n", int(timeout/time.Second), result.ErrMsg)
// 	} else if !result.IsOk {
// 		log.Status = TASK_ERROR
// 		log.Error = "ERROR:" + result.ErrMsg
// 	}

// 	if err := Add(log); err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	return log.ID
// }

// func TaskLogAdd(t *TaskLog) (int64, error) {
// 	return orm.NewOrm().Insert(t)
// }

// func TaskLogGetList(page, pageSize int, filters ...interface{}) ([]*TaskLog, int64) {
// 	offset := (page - 1) * pageSize

// 	logs := make([]*TaskLog, 0)

// 	query := orm.NewOrm().QueryTable(TableName("task_log"))
// 	if len(filters) > 0 {
// 		l := len(filters)
// 		for k := 0; k < l; k += 2 {
// 			query = query.Filter(filters[k].(string), filters[k+1])
// 		}
// 	}

// 	total, _ := query.Count()
// 	query.OrderBy("-id").Limit(pageSize, offset).All(&logs)

// 	return logs, total
// }

// func TaskLogGetById(id int) (*TaskLog, error) {
// 	obj := &TaskLog{
// 		ID: id,
// 	}

// 	err := orm.NewOrm().Read(obj)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return obj, nil
// }

// func TaskLogDelById(id int) error {
// 	_, err := orm.NewOrm().Delete(&TaskLog{ID: id})
// 	return err
// }

// func TaskLogDelByTaskId(taskId int) (int64, error) {
// 	return orm.NewOrm().QueryTable(TableName("task_log")).Filter("task_id", taskId).Delete()
// }

// func GetLogNum(status int) (int64, error) {
// 	return orm.NewOrm().QueryTable(TableName("task_log")).Filter("status", status).Count()
// }

// type SumDays struct {
// 	Day string
// 	Sum int
// }

// func SumByDays(limit int, status string) orm.Params {

// 	var m = map[string]string{
// 		"0":  "okNum",
// 		"-1": "errNum",
// 		"-2": "expiredRun"}

// 	res := make(orm.Params)
// 	key := m[status]

// 	if RunNumCache.IsExist(key) {
// 		json.Unmarshal(RunNumCache.Get(key).([]byte), &res)
// 		logs.Info("cache")
// 		return res
// 	}
// 	_, err := orm.NewOrm().Raw("SELECT FROM_UNIXTIME(create_time,'%Y-%m-%d') days,COUNT(id) count FROM pp_task_log WHERE status in(?) GROUP BY days ORDER BY days DESC limit ?;",
// 		status, limit).RowsToMap(&res, "days", "count")

// 	if err != nil {
// 		return nil
// 	}

// 	data, err := json.Marshal(res)
// 	if err != nil {
// 		return nil
// 	}
// 	RunNumCache.Put(key, data, 2*time.Hour)
// 	return res

// }
