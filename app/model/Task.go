/************************************************************
** @Description: model
** @Author: haodaquan
** @Date:   2018-06-11 21:26
** @Last Modified by:   Bee
** @Last Modified time: 2019-02-15 21:32
*************************************************************/
package model

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	TASK_SUCCESS = 0  // 任务执行成功
	TASK_ERROR   = -1 // 任务执行出错
	TASK_TIMEOUT = -2 // 任务执行超时
)

type Task struct {
	ID            int    `xorm:"id pk" json:"id"`
	GroupID       int    `xorm:"group_id" json:"group_id"`
	ServerIDs     string `xorm:"server_ids" json:"server_ids"`
	ServerType    int    `xorm:"server_type" json:"server_type"`
	TaskName      string `xorm:"task_name" json:"task_name"`
	Description   string `xorm:"description" json:"description"`
	CronSpec      string `xorm:"cron_spec" json:"cron_spec"`
	Concurrent    int    `xorm:"concurrent" json:"concurrent"`
	Command       string `xorm:"command" json:"command"`
	Timeout       int    `xorm:"timeout" json:"timeout"`
	ExecuteTimes  int    `xorm:"execute_times" json:"execute_times"`
	PrevTime      int64  `xorm:"prev_time" json:"prev_time"`
	Status        int    `xorm:"status" json:"status"`
	IsNotify      int    `xorm:"is_notify" json:"is_notify"`
	NotifyType    int    `xorm:"notify_type" json:"notify_type"`
	NotifyTplID   int    `xorm:"notify_tpl_id" json:"notify_tpl_id"`
	NotifyUserIds string `xorm:"notify_user_ids" json:"notify_user_ids"`
	CreatedID     int    `xorm:"create_id" json:"created_id"`
	UpdatedID     int    `xorm:"update_id" json:"updated_id"`
	CreatedAt     int64  `xorm:"create_time" json:"created_at"`
	UpdatedAt     int64  `xorm:"update_time" json:"created_at"`
}

func (t *Task) TableName() string {
	return TableName("task")
}

func (t *Task) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func TaskAdd(task *Task) (int64, error) {
	if task.TaskName == "" {
		return 0, fmt.Errorf("任务名称不能为空")
	}

	if task.CronSpec == "" {
		return 0, fmt.Errorf("时间表达式不能为空")
	}
	if task.Command == "" {
		return 0, fmt.Errorf("命令内容不能为空")
	}
	if task.CreatedAt == 0 {
		task.CreatedAt = time.Now().Unix()
	}
	return orm.NewOrm().Insert(task)
}

func TaskGetList(page, pageSize int, filters ...interface{}) ([]*Task, int64) {
	offset := (page - 1) * pageSize

	tasks := make([]*Task, 0)

	query := orm.NewOrm().QueryTable(TableName("task"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	query.OrderBy("-id", "-status", "task_name").Limit(pageSize, offset).All(&tasks)

	return tasks, total
}

func TaskResetGroupId(groupId int) (int64, error) {
	return orm.NewOrm().QueryTable(TableName("task")).Filter("group_id", groupId).Update(orm.Params{
		"group_id": 0,
	})
}

func TaskGetById(id int) (*Task, error) {
	task := &Task{
		ID: id,
	}

	err := orm.NewOrm().Read(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

//修改为逻辑删除
func TaskDel(id int) (int64, error) {
	return orm.NewOrm().QueryTable(TableName("task")).Filter("id", id).Update(orm.Params{
		"status": -1,
	})

	// _, err := orm.NewOrm().QueryTable(TableName("task")).Filter("id", id).Delete()
	// return err
}

//运行总次数
func TaskTotalRunNum() (int64, error) {

	res := make(orm.Params)
	_, err := orm.NewOrm().Raw("select sum(execute_times) as num,task_name from pp_task").RowsToMap(&res, "num", "task_name")

	if err != nil {
		return 0, err
	}

	for k, _ := range res {
		i64, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return 0, err
		}

		return i64, nil

	}
	return 0, nil
}
