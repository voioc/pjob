package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/worker"
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

	if err := model.GetDB().Where(condition).OrderBy("field(status, 1, 2, 3, 0), id desc").Limit(pageSize, offset).Find(&tasks); err != nil {
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

func (s *TaskService) CreateJob(task *model.Task) ([]*worker.Job, error) {
	jobs := make([]*worker.Job, 0)

	ServerIDs := strings.Split(task.ServerIDs, ",")
	for _, serverID := range ServerIDs {
		var job *worker.Job
		server := model.TaskServer{}

		// 本地执行
		if serverID == "0" {
			job = worker.NewCommandJob(task.ID, 0, task.TaskName, task.Command)
			server.ServerName = "本地服务器"
			// job.Task = task
			// job.Concurrent = false
			// if task.Concurrent == 1 {
			// 	job.Concurrent = true
			// }

			// job.ServerID = 0
			// job.ServerName = "本地服务器"
			// jobs = append(jobs, job)

		} else { // 远程执行
			sid, _ := strconv.Atoi(serverID)
			if err := model.DataByID(&server, sid); err != nil {
				fmt.Println(err.Error())
			}

			if server.Status == 2 {
				fmt.Println("服务器" + serverID + "已禁用")
				continue
			}

			if server.ConnectionType == 1 { // 1 ssh 2 telnet 3 agent
				if server.Type == 1 { // 1 密码 2 私钥
					// 密码验证登录服务器
					job = worker.RemoteCommandJobByPassword(task.ID, sid, task.TaskName, task.Command, &server)
					// if task.Concurrent == 1 {
					// 	job.Concurrent = true
					// }

					// job.ServerName = server.ServerName

					// jobs = append(jobs, job)
				} else {
					job = worker.RemoteCommandJob(task.ID, sid, task.TaskName, task.Command, &server)
					// job.Task = task
					// if task.Concurrent == 1 {
					// 	job.Concurrent = true

					// }
					// job.ServerName = server.ServerName
					// jobs = append(jobs, job)
				}
			} else if server.ConnectionType == 2 {
				if server.Type == 1 {
					//密码验证登录服务器
					job = worker.RemoteCommandJobByTelnetPassword(task.ID, sid, task.TaskName, task.Command, &server)
					// job.Task = task
					// job.Concurrent = false
					// if task.Concurrent == 1 {
					// 	job.Concurrent = true
					// }

					// job.ServerName = server.ServerName
					// jobs = append(jobs, job)
				}
			} else if server.ConnectionType == 3 {
				//密码验证登录服务器
				job = worker.RemoteCommandJobByAgentPassword(task.ID, sid, task.TaskName, task.Command, &server)
				// if task.Concurrent == 1 {
				// 	job.Concurrent = true
				// }

				//job.Concurrent = task.Concurrent == 1
				// job.ServerName = server.ServerName
				// jobs = append(jobs, job)

			}
		}

		// 设置回调函数
		job.LogFunc = func(result *worker.JobResult, t time.Time) int {
			log := model.TaskLog{
				TaskID:      job.ID,
				ServerID:    job.ServerID,
				ServerName:  job.ServerName,
				Output:      result.OutMsg,
				Error:       result.ErrMsg,
				ProcessTime: int(time.Since(t) / time.Millisecond),
				CreatedAt:   t.Unix(),
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

		if task.Concurrent == 1 {
			job.Concurrent = true
		}

		job.ServerName = server.ServerName

		jobs = append(jobs, job)
	}

	return jobs, nil
}

// 加载任务
func (s *TaskService) Loading() {
	list := make([]*model.Task, 0)
	if err := model.List(&list, 1, 1000000, "status =", 1); err != nil {
		fmt.Println(err.Error())
	}

	for _, task := range list {
		// 创建定时Job
		jobs, _ := s.CreateJob(task)

		// 开启任务
		for _, job := range jobs {
			if worker.AddJob(task.CronSpec, job) {
				// task.Status = 1
				// // task.Update()
				// if err := model.Update(task.ID, &task); err != nil {
				// 	fmt.Println(err.Error())
				// }
			}
		}
	}

}
