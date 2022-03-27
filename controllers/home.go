/**********************************************
** @Des: This file ...
** @Author: haodaquan
** @Date:   2017-09-08 10:21:13
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-09 18:04:41
***********************************************/
package controllers

import (
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"text/template"
	"time"

	"github.com/astaxie/beego"
	"github.com/gin-gonic/gin"
	"github.com/voioc/pjob/jobs"
	"github.com/voioc/pjob/libs"
	"github.com/voioc/pjob/models"
	"github.com/voioc/pjob/service"
)

type HomeController struct {
	BaseController
}

// Index dk
func (self *HomeController) Index(c *gin.Context) {
	// self.Data["pageTitle"] = "系统首页"
	//self.display()
	// self.TplName = "public/main.html"

	uid := c.GetInt("uid")
	menu, _ := service.Menu(uid)
	c.HTML(http.StatusOK, "main.html", template.FuncMap{
		"siteName":  "系统首页",
		"SideMenu1": menu["SideMenu1"],
		"SideMenu2": menu["SideMenu2"],
	})
}

func (self *HomeController) Help() {
	self.Data["pageTitle"] = "Cron表达式说明"
	//self.display()
	self.TplName = "public/help.html"
}

func (self *HomeController) Start(c *gin.Context) {

	data := map[string]interface{}{}
	//总任务数量
	_, count := models.TaskGetList(1, 10)
	// self.Data["totalJob"] = count
	data["totalJob"] = count

	//日志总量
	_, totalLog := models.TaskLogGetList(1, 10)
	data["totalLog"] = totalLog

	//待审核任务数量
	_, totalAuditTask := models.TaskGetList(1, 10, "status", 2)
	data["totalAuditTask"] = totalAuditTask

	//失败
	errorNum, err := models.GetLogNum(-1)
	if err != nil {
		errorNum = 0
	}
	data["errorNum"] = errorNum

	//成功
	successNum, err := models.GetLogNum(0)
	if err != nil {
		successNum = 0
	}
	// self.Data["successNum"] = successNum
	data["successNum"] = successNum

	//用户数
	_, userNum := models.AdminGetList(1, 10, "status", 1)
	// self.Data["userNum"] = userNum
	data["userNum"] = userNum

	//累计运行总次数
	n, err := models.TaskTotalRunNum()
	if err != nil {
		n = 0
	}
	data["TaskTotalRunNum"] = n

	groups_map := serverGroupLists(self.serverGroups, self.userId)
	//计算总任务数量

	// 即将执行的任务
	entries := jobs.GetEntries(30)
	jobList := make([]map[string]interface{}, len(entries))
	startJob := 0 //即将执行的任务
	for k, v := range entries {
		row := make(map[string]interface{})
		job := v.Job.(*jobs.Job)
		task, _ := models.TaskGetById(job.GetId())
		row["task_id"] = job.GetId()
		row["task_name"] = job.GetName()
		row["task_group"] = groups_map[task.GroupId]
		row["next_time"] = beego.Date(v.Next, "Y-m-d H:i:s")
		jobList[k] = row
		startJob++
	}

	data["recentLogs"] = jobList

	// 最近执行失败的日志
	logs, _ := models.TaskLogGetList(1, 30, "status__lt", 0)
	errLogs := make([]map[string]interface{}, len(logs))

	for k, v := range logs {
		task, err := models.TaskGetById(v.TaskId)
		taskName := ""
		if err == nil {
			taskName = task.TaskName
		}

		row := make(map[string]interface{})
		row["task_name"] = taskName
		row["id"] = v.Id
		row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["process_time"] = float64(v.ProcessTime) / 1000
		row["ouput_size"] = libs.SizeFormat(float64(len(v.Output)))
		row["error"] = beego.Substr(v.Error, 0, 100)
		row["status"] = v.Status
		errLogs[k] = row

	}
	data["errLogs"] = errLogs
	data["startJob"] = startJob
	data["jobs"] = jobList

	//折线图
	okRun := models.SumByDays(30, "0")
	errRun := models.SumByDays(30, "-1")
	expiredRun := models.SumByDays(30, "-2")

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
		i, _ := strconv.ParseInt(v.(string), 10, 64)
		ss = append(ss, kv{k, i})
	}

	sort.Slice(ss, func(i, j int) bool {

		return ss[i].Key < ss[j].Key
	})

	for _, v := range ss {

		days = append(days, v.Key)
		okNum = append(okNum, v.Value)

		if _, ok := errRun[v.Key]; ok {
			i, _ := strconv.ParseInt(errRun[v.Key].(string), 10, 64)
			errNum = append(errNum, i)
		} else {
			errNum = append(errNum, 0)
		}

		if _, ok := expiredRun[v.Key]; ok {
			i, _ := strconv.ParseInt(expiredRun[v.Key].(string), 10, 64)
			expiredNum = append(expiredNum, i)
		} else {
			expiredNum = append(expiredNum, 0)
		}
	}

	data["days"] = days
	data["okNum"] = okNum
	data["errNum"] = errNum
	data["expiredNum"] = expiredNum

	data["cpuNum"] = runtime.NumCPU()

	//系统运行信息
	info := libs.SystemInfo(models.StartTime)
	data["sysInfo"] = info

	data["pageTitle"] = "系统概况"
	// self.display()

	c.HTML(http.StatusOK, "start.html", data)
}
