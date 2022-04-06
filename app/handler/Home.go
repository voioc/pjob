/**********************************************
** @Des: This file ...
** @Author: haodaquan
** @Date:   2017-09-08 10:21:13
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-09 18:04:41
***********************************************/
package handler

import (
	"fmt"
	"net/http"
	"runtime"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/app/service"
	"github.com/voioc/cjob/libs"
	"github.com/voioc/cjob/utils"
	"github.com/voioc/cjob/worker"
)

type HomeController struct {
	BaseController
}

// Index dk
func (self *HomeController) Index(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	data["siteName"] = "系统首页"
	data["loginUserName"] = "管理员"

	//self.display()
	// self.TplName = "public/main.html"

	uid := c.GetInt("uid")
	menu, _ := service.AuthS(c).Menu(uid)
	// fmt.Println(menu)
	data["SideMenu1"] = menu["SideMenu1"]
	data["SideMenu2"] = menu["SideMenu2"]

	c.HTML(http.StatusOK, "main.html", data)
}

func (self *HomeController) Help(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	data["pageTitle"] = "Cron表达式说明"

	//self.display()
	// self.TplName = "public/help.html"
	c.HTML(http.StatusOK, "public/help.html", data)
}

func (self *HomeController) Start(c *gin.Context) {

	data := map[string]interface{}{}
	// 总任务数量
	// _, count := model.TaskGetList(1, 10)
	// _, count, err := service.TaskS(c).TaskGetList(1, 10)
	count, err := model.ListCount(&model.Task{})
	// self.Data["totalJob"] = count
	data["totalJob"] = count

	//日志总量
	// _, totalLog, err := service.TaskLogS(c).LogList(1, 10)
	totalLog, err := model.ListCount(&model.TaskLog{})
	data["totalLog"] = totalLog

	// 待审核任务数量
	_, totalAuditTask, err := service.TaskS(c).TaskGetList(1, 10, "status = ", 2)
	data["totalAuditTask"] = totalAuditTask

	//失败
	errorNum, err := service.TaskLogS(c).GetLogNum(-1)
	if err != nil {
		errorNum = 0
	}
	data["errorNum"] = errorNum

	//成功
	successNum, err := service.TaskLogS(c).GetLogNum(0)
	if err != nil {
		successNum = 0
	}
	// self.Data["successNum"] = successNum
	data["successNum"] = successNum

	// 用户数
	// _, userNum, err := service.AdminS(c).AdminList(1, 10, "status = ", 1)
	userNum, err := model.ListCount(&model.Admin{}, "status = ", 1)
	data["userNum"] = userNum

	// 累计运行总次数
	n, _ := service.TaskS(c).TaskTotalRunNum()
	data["TaskTotalRunNum"] = n

	uid := c.GetInt("uid")
	_, sg := service.AuthS(c).TaskGroups(uid, c.GetString("role_id"))
	groups_map := service.ServerGroupS(c).ServerGroupLists(sg, uid)
	// 计算总任务数量

	// 即将执行的任务
	entries := worker.GetEntries(30)
	jobList := make([]map[string]interface{}, len(entries))
	startJob := 0 //即将执行的任务
	for k, v := range entries {
		job := v.Job.(*worker.Job)

		// task, _ := service.TaskS(c).TaskByID(job.GetId())
		task := &model.Task{}
		if err := model.DataByID(task, job.GetTaskID()); err != nil {
			fmt.Println(err.Error())
		}

		row := make(map[string]interface{})
		row["task_id"] = job.GetTaskID()
		row["task_name"] = job.GetName()
		row["task_group"] = groups_map[task.GroupID]
		row["next_time"] = v.Next.Format("2006-01-02 15:04:05")
		jobList[k] = row
		startJob++
	}

	data["recentLogs"] = jobList

	// 最近执行失败的日志
	// logs, _, _ := service.TaskLogS(c).LogList(1, 30, "status != ", 0)
	logs := make([]model.TaskLog, 0)
	if err := model.List(&logs, 1, 30, "status != ", 0); err != nil {
		fmt.Println(err.Error())
	}

	errLogs := make([]map[string]interface{}, len(logs))

	for k, v := range logs {
		// task, err := service.TaskS(c).TaskByID(v.TaskID)
		task := model.Task{}
		if err := model.DataByID(&task, v.TaskID); err != nil {
			fmt.Println(err.Error())
		}

		taskName := ""
		if err == nil {
			taskName = task.TaskName
		}

		row := make(map[string]interface{})
		row["task_name"] = taskName
		row["id"] = v.ID
		row["start_time"] = time.Unix(v.CreatedAt, 0).Format("2006-01-02 15:04:05")
		row["process_time"] = float64(v.ProcessTime) / 1000
		row["output_size"] = libs.SizeFormat(float64(len(v.Output)))
		row["error"] = v.Error // v.Error[0:100]
		row["status"] = v.Status
		errLogs[k] = row

	}
	data["errLogs"] = errLogs
	data["startJob"] = startJob
	data["jobs"] = jobList

	// 折线图
	okRun, _ := service.TaskLogS(c).SumByDays(30, "0")
	errRun, _ := service.TaskLogS(c).SumByDays(30, "-1")
	expiredRun, _ := service.TaskLogS(c).SumByDays(30, "-2")
	// fmt.Println(okRun)

	days := []string{}
	okNum := []int64{}
	errNum := []int64{}
	expiredNum := []int64{}

	type kv struct {
		Key   string
		Value int64
	}

	//排序
	var ss []kv
	for k, v := range okRun {
		// i, _ := strconv.ParseInt(v, 10, 64)
		ss = append(ss, kv{k, int64(v)})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Key < ss[j].Key
	})

	for _, v := range ss {
		days = append(days, v.Key)
		okNum = append(okNum, v.Value)

		value := 0
		if _, ok := errRun[v.Key]; ok {
			// i, _ := strconv.ParseInt(v.Key, 10, 64)
			value = errRun[v.Key]
		}
		errNum = append(errNum, int64(value))

		value = 0
		if _, ok := expiredRun[v.Key]; ok {
			// i, _ := strconv.ParseInt(expiredRun[v.Key], 10, 64)
			value = expiredRun[v.Key]
		}
		expiredNum = append(expiredNum, int64(value))
	}

	data["days"] = days
	data["okNum"] = okNum
	data["errNum"] = errNum
	data["expiredNum"] = expiredNum

	data["cpuNum"] = runtime.NumCPU()

	//系统运行信息
	info := libs.SystemInfo(model.StartTime)
	data["sysInfo"] = info

	data["pageTitle"] = "系统概况"
	// self.display()

	c.HTML(http.StatusOK, "start.html", data)
}
