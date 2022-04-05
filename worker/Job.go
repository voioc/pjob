package worker

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"runtime/debug"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/libs"
	"github.com/voioc/coco/logzap"

	"strconv"
	"strings"

	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/voioc/cjob/notify"
)

type Job struct {
	JobKey      int // jobId = id*10000+serverId
	ID          int // taskID
	TaskID      int
	TaskName    string
	LogID       int                            // 日志记录ID
	ServerID    int                            // 执行器信息
	ServerName  string                         // 执行器名称
	ServerType  int                            // 执行器类型，1-ssh 2-telnet 3-agent
	Name        string                         // 任务名称
	Task        *model.Task                    // 任务对象
	RunFunc     func(time.Duration) *JobResult // 执行函数
	SuffixFunc  func(*Job, *JobResult)         // 任务执行完成
	Timeout     int                            // 超时时间:秒
	Status      int                            // 任务状态，大于0表示正在执行中
	Concurrent  bool                           // 同一个任务是否允许并行执行
	StartAt     time.Time                      // 开始时间
	ProcessTime time.Duration                  // 花费时间
}

func (j *Job) GetStatus() int {
	return j.Status
}

func (j *Job) GetName() string {
	return j.Name
}

func (j *Job) GetID() int {
	return j.ID
}

func (j *Job) GetLogID() int {
	return j.LogID
}

func (j *Job) agentRun() (reply *JobResult) {
	// server, _ := model.TaskServerGetById(j.ServerId)
	server := model.TaskServer{}
	if err := model.DataByID(&server, j.ServerID); err != nil {
		fmt.Println(err.Error())
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server.ServerIP, server.Port))
	reply = new(JobResult)
	if err != nil {
		logs.Error("Net error:", err)
		reply.IsOk = false
		reply.ErrMsg = "Net error:" + err.Error()
		reply.IsTimeout = false
		reply.OutMsg = ""
		return reply
	}

	defer conn.Close()
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	defer client.Close()
	reply = new(JobResult)

	task := j.Task
	err = client.Call("RpcTask.RunTask", task, &reply)
	if err != nil {
		reply.IsOk = false
		reply.ErrMsg = "Net error:" + err.Error()
		reply.IsTimeout = false
		reply.OutMsg = ""
		return reply
	}
	return
}

func (j *Job) Run() {
	// 执行策略 1 同时执行 2 轮询
	if j.ServerType == 2 {
		if !PollServer(j) {
			return
		} else {
			SetCounter(strconv.Itoa(j.Task.ID))
		}
	}

	if !j.Concurrent && j.Status > 0 {
		logzap.Wx(nil, "Worker", fmt.Sprintf("任务[%d]上一次执行尚未结束，本次被忽略。", j.JobKey))
		return
	}

	defer func() {
		if err := recover(); err != nil {
			logzap.Ex(nil, "Worker", "%+v\n", string(debug.Stack()))
		}
	}()

	if workPool != nil {
		workPool <- true
		defer func() {
			<-workPool
		}()
	}

	fmt.Println(fmt.Sprintf("开始执行任务: %d", j.JobKey))

	j.Status++
	defer func() {
		j.Status--
	}()

	t := time.Now()
	j.StartAt = time.Now()

	timeout := time.Duration(time.Hour * 24)
	if j.Timeout > 0 {
		timeout = time.Second * time.Duration(j.Timeout)
	}

	var jobResult = new(JobResult)

	// agent
	if j.ServerType == 2 {
		jobResult = j.agentRun()
	} else {
		jobResult = j.RunFunc(timeout)
	}

	// if j.LogFunc != nil {
	// 	j.LogID = j.LogFunc(jobResult, j.StartAt)
	// }

	ut := time.Now().Sub(j.StartAt) / time.Millisecond
	j.ProcessTime = ut

	if j.SuffixFunc != nil {
		j.SuffixFunc(j, jobResult)
	}

	// 插入日志
	log := new(model.TaskLog)
	log.TaskID = j.ID
	log.ServerID = j.ServerID
	log.ServerName = j.ServerName
	log.Output = jobResult.OutMsg
	log.Error = jobResult.ErrMsg
	log.ProcessTime = int(ut)
	log.CreatedAt = j.StartAt.Unix()

	if jobResult.IsTimeout {
		log.Status = model.TASK_TIMEOUT
		log.Error = fmt.Sprintf("任务执行超过 %d 秒\n----------------------\n%s\n", int(timeout/time.Second), jobResult.ErrMsg)
	} else if !jobResult.IsOk {
		log.Status = model.TASK_ERROR
		log.Error = "ERROR:" + jobResult.ErrMsg
	}

	if log.Status < 0 && j.Task.IsNotify == 1 {
		if j.Task.NotifyUserIDs != "0" && j.Task.NotifyUserIDs != "" {
			adminInfo := AllAdminInfo(j.Task.NotifyUserIDs)
			phone := make(map[string]string, 0)
			dingtalk := make(map[string]string, 0)
			wechat := make(map[string]string, 0)
			toEmail := ""
			for _, v := range adminInfo {
				if v.Phone != "0" && v.Phone != "" {
					phone[v.Phone] = v.Phone
				}
				if v.Email != "0" && v.Email != "" {
					toEmail += v.Email + ";"
				}
				if v.Dingtalk != "0" && v.Dingtalk != "" {
					dingtalk[v.Dingtalk] = v.Dingtalk
				}
				if v.Wechat != "0" && v.Wechat != "" {
					wechat[v.Wechat] = v.Wechat
				}
			}
			toEmail = strings.TrimRight(toEmail, ";")

			TextStatus := []string{
				"超时",
				"错误",
				"正常",
			}

			status := log.Status + 2

			title, content, taskOutput, errOutput := "", "", "", ""

			// notifyTpl, err := model.NotifyTplGetById(j.Task.NotifyTplID)
			notifyTpl := model.NotifyTpl{}
			if err := model.DataByID(&notifyTpl, j.Task.NotifyTplID); err != nil {
				notifyTpl, err := model.NotifyTplGetByTplType(j.Task.NotifyType, model.NotifyTplTypeSystem)
				if err == nil {
					title = notifyTpl.Title
					content = notifyTpl.Content
				}
			} else {
				title = notifyTpl.Title
				content = notifyTpl.Content
			}

			taskOutput = strings.Replace(log.Output, "\n", " ", -1)
			taskOutput = strings.Replace(taskOutput, "\"", "\\\"", -1)
			errOutput = strings.Replace(log.Error, "\n", " ", -1)
			errOutput = strings.Replace(errOutput, "\"", "\\\"", -1)

			if title != "" {
				title = strings.Replace(title, "{{TaskId}}", strconv.Itoa(j.Task.ID), -1)
				title = strings.Replace(title, "{{ServerId}}", strconv.Itoa(j.ServerID), -1)
				title = strings.Replace(title, "{{TaskName}}", j.Task.TaskName, -1)
				title = strings.Replace(title, "{{ExecuteCommand}}", j.Task.Command, -1)
				title = strings.Replace(title, "{{ExecuteTime}}", beego.Date(time.Unix(log.CreatedAt, 0), "Y-m-d H:i:s"), -1)
				title = strings.Replace(title, "{{ProcessTime}}", strconv.FormatFloat(float64(log.ProcessTime)/1000, 'f', 6, 64), -1)
				title = strings.Replace(title, "{{ExecuteStatus}}", TextStatus[status], -1)
				title = strings.Replace(title, "{{TaskOutput}}", taskOutput, -1)
				title = strings.Replace(title, "{{ErrorOutput}}", errOutput, -1)
			}

			if content != "" {
				content = strings.Replace(content, "{{TaskId}}", strconv.Itoa(j.Task.ID), -1)
				content = strings.Replace(content, "{{ServerId}}", strconv.Itoa(j.ServerID), -1)
				content = strings.Replace(content, "{{TaskName}}", j.Task.TaskName, -1)
				content = strings.Replace(content, "{{ExecuteCommand}}", strings.Replace(j.Task.Command, "\"", "\\\"", -1), -1)
				content = strings.Replace(content, "{{ExecuteTime}}", beego.Date(time.Unix(log.CreatedAt, 0), "Y-m-d H:i:s"), -1)
				content = strings.Replace(content, "{{ProcessTime}}", strconv.FormatFloat(float64(log.ProcessTime)/1000, 'f', 6, 64), -1)
				content = strings.Replace(content, "{{ExecuteStatus}}", TextStatus[status], -1)
				content = strings.Replace(content, "{{TaskOutput}}", taskOutput, -1)
				content = strings.Replace(content, "{{ErrorOutput}}", errOutput, -1)
			}

			if j.Task.NotifyType == 0 && toEmail != "" {
				//邮件
				mailtype := "html"
				ok := notify.SendToChan(toEmail, title, content, mailtype)
				if !ok {
					fmt.Println("发送邮件错误", toEmail)
				}
			} else if j.Task.NotifyType == 1 && len(phone) > 0 {
				//信息
				param := make(map[string]string)
				err := json.Unmarshal([]byte(content), &param)
				if err != nil {
					fmt.Println("发送信息错误", err)
					return
				}

				ok := notify.SendSmsToChan(phone, param)
				if !ok {
					fmt.Println("发送信息错误", phone)
				}
			} else if j.Task.NotifyType == 2 && len(dingtalk) > 0 {
				//钉钉
				param := make(map[string]interface{})

				err := json.Unmarshal([]byte(content), &param)
				if err != nil {
					fmt.Println("发送钉钉错误", err)
					return
				}

				ok := notify.SendDingtalkToChan(dingtalk, param)
				if !ok {
					fmt.Println("发送钉钉错误", dingtalk)
				}
			} else if j.Task.NotifyType == 3 && len(wechat) > 0 {
				//微信
				param := make(map[string]string)
				err := json.Unmarshal([]byte(content), &param)
				if err != nil {
					fmt.Println("发送微信错误", err)
					return
				}

				ok := notify.SendWechatToChan(phone, param)
				if !ok {
					fmt.Println("发送微信错误", phone)
				}
			}
		}
	}

	// j.LogId, _ = model.TaskLogAdd(log)

	// if err := model.Add(log); err != nil {
	// 	fmt.Println(err.Error())
	// }
	// j.LogID = int64(log.ID)

	// 更新上次执行时间
	j.Task.PrevTime = t.Unix()
	j.Task.ExecuteTimes++
	// j.Task.Update("PrevTime", "ExecuteTimes")
	if err := model.Update(j.Task.ID, j.Task); err != nil {
		fmt.Println(err.Error())
	}
}

type JobResult struct {
	OutMsg    string
	ErrMsg    string
	IsOk      bool
	IsTimeout bool
}

//调度计数器
var Counter sync.Map

func GetCounter(key string) int {
	if v, ok := Counter.LoadOrStore(key, 0); ok {
		n := v.(int)
		return n
	}
	return 0
}

func SetCounter(key string) {
	if v, ok := Counter.Load(key); ok {
		n := v.(int)
		m := n + 1
		if n > 1000 {
			m = 0
		}
		Counter.Store(key, m)
	}
}

func NewJobFromTask(task *model.Task) ([]*Job, error) {
	if task.ID < 1 {
		return nil, fmt.Errorf("ToJob: 缺少id")
	}

	if task.ServerIDs == "" {
		return nil, fmt.Errorf("任务执行失败，找不到执行的服务器")
	}

	TaskServerIdsArr := strings.Split(task.ServerIDs, ",")
	jobArr := make([]*Job, 0)
	for _, server_id := range TaskServerIdsArr {
		if server_id == "0" {
			//本地执行
			job := NewCommandJob(task.ID, 0, task.TaskName, task.Command)
			// job.Task = task
			job.Concurrent = false
			if task.Concurrent == 1 {
				job.Concurrent = true
			}
			// job.Concurrent = task.Concurrent == 1
			job.ServerID = 0
			job.ServerName = "本地服务器"
			jobArr = append(jobArr, job)
		} else {
			server_id_int, _ := strconv.Atoi(server_id)
			// 远程执行
			// server, _ := model.TaskServerGetById(server_id_int)
			server := model.TaskServer{}
			if err := model.DataByID(&server, server_id_int); err != nil {
				fmt.Println(err.Error())
			}

			if server.Status == 2 {
				fmt.Println("服务器已禁用")
				continue
			}

			if server.ConnectionType == 0 { // 0 ssh 1 telnet
				if server.Type == 0 { // 0 密码 1 私钥
					// 密码验证登录服务器
					job := RemoteCommandJobByPassword(task.ID, server_id_int, task.TaskName, task.Command, &server)
					// job.Task = task
					job.Concurrent = false
					if task.Concurrent == 1 {
						job.Concurrent = true
					}
					//job.Concurrent = task.Concurrent == 1
					// job.ServerId = server_id_int
					job.ServerName = server.ServerName
					jobArr = append(jobArr, job)
				} else {
					job := RemoteCommandJob(task.ID, server_id_int, task.TaskName, task.Command, &server)
					// job.Task = task
					job.Concurrent = false
					if task.Concurrent == 1 {
						job.Concurrent = true
					}
					//job.Concurrent = task.Concurrent == 1
					// job.ServerId = server_id_int
					job.ServerName = server.ServerName
					jobArr = append(jobArr, job)
				}
			} else if server.ConnectionType == 1 {
				if server.Type == 0 {
					//密码验证登录服务器
					job := RemoteCommandJobByTelnetPassword(task.ID, server_id_int, task.TaskName, task.Command, &server)
					// job.Task = task
					job.Concurrent = false
					if task.Concurrent == 1 {
						job.Concurrent = true
					}
					//job.Concurrent = task.Concurrent == 1
					// job.ServerId = server_id_int
					job.ServerName = server.ServerName
					jobArr = append(jobArr, job)
				}
			} else if server.ConnectionType == 2 {
				//密码验证登录服务器
				job := RemoteCommandJobByAgentPassword(task.ID, server_id_int, task.TaskName, task.Command, &server)
				// job.Task = task
				job.Concurrent = false
				if task.Concurrent == 1 {
					job.Concurrent = true
				}
				//job.Concurrent = task.Concurrent == 1
				// job.ServerId = server_id_int
				job.ServerName = server.ServerName
				jobArr = append(jobArr, job)

			}
		}
	}

	return jobArr, nil
}

type RpcResult struct {
	Status  int
	Message string
}

func TestServer(server *model.TaskServer) error {
	if server.ConnectionType == 0 {
		switch server.Type {
		case 0:
			//密码登录
			return libs.RemoteCommandByPassword(server)
		case 1:
			//密钥登录
			return libs.RemoteCommandByKey(server)
		default:
			return errors.New("未知的登录方式")

		}
	} else if server.ConnectionType == 1 {
		if server.Type == 0 {
			//密码登录]
			return libs.RemoteCommandByTelnetPassword(server)
		} else {
			return errors.New("Telnet方式暂不支持密钥登陆！")
		}

	} else if server.ConnectionType == 2 {
		return libs.RemoteAgent(server)
	}

	return errors.New("未知错误")
}

func PollServer(j *Job) bool {
	//判断是否是当前执行器执行
	TaskServerIdsArr := strings.Split(j.Task.ServerIDs, ",")
	num := len(TaskServerIdsArr)

	if num == 0 {
		return false
	}

	count := GetCounter(strconv.Itoa(j.Task.ID))
	index := count % num
	pollServerId, _ := strconv.Atoi(TaskServerIdsArr[index])

	if j.ServerID != pollServerId {
		return false
	}

	//本地服务器
	if pollServerId == 0 {
		return true
	}

	//判断执行器或者服务器是否存活
	// server, _ := model.TaskServerGetById(pollServerId)
	server := model.TaskServer{}
	if err := model.DataByID(&server, pollServerId); err != nil {
		fmt.Println(err.Error())
	}

	if server.Status != 0 {
		return false
	}

	if err := TestServer(&server); err != nil {
		server.Status = 1
		if err := model.Update(server.ID, server); err != nil {
			fmt.Println(err.Error())
		}
		return false
	} else {
		server.Status = 0
		if err := model.Update(server.ID, server, true); err != nil {
			fmt.Println(err.Error())
		}
	}

	return true
}

// 冗余代码
type adminInfo struct {
	Id       int
	Email    string
	Phone    string
	Dingtalk string
	Wechat   string
	RealName string
}

func AllAdminInfo(adminIds string) []*adminInfo {
	Filters := make([]interface{}, 0)
	Filters = append(Filters, "status", 1)
	//Filters = append(Filters, "id__gt", 1)
	var notifyUserIds []int
	if adminIds != "0" && adminIds != "" {
		notifyUserIdsStr := strings.Split(adminIds, ",")
		for _, v := range notifyUserIdsStr {
			i, _ := strconv.Atoi(v)
			notifyUserIds = append(notifyUserIds, i)
		}
		Filters = append(Filters, "id__in", notifyUserIds)
	}
	// Result, _ := model.AdminGetList(1, 1000, Filters...)
	Result := make([]model.Admin, 0)
	if err := model.List(&Result, 1, 1000, Filters...); err != nil {
		fmt.Println(err.Error())
	}

	adminInfos := make([]*adminInfo, 0)
	for _, v := range Result {
		ai := adminInfo{
			Id:       v.ID,
			Email:    v.Email,
			Phone:    v.Phone,
			Dingtalk: v.Dingtalk,
			Wechat:   v.Wechat,
			RealName: v.RealName,
		}
		adminInfos = append(adminInfos, &ai)
	}

	return adminInfos
}
