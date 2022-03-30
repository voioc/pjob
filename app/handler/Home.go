/**********************************************
** @Des: This file ...
** @Author: haodaquan
** @Date:   2017-09-08 10:21:13
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-09 18:04:41
***********************************************/
package handler

import (
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/app/service"
	"github.com/voioc/cjob/jobs"
	"github.com/voioc/cjob/libs"
	"github.com/voioc/cjob/utils"
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
	menu, _ := service.Menu(uid)
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
	//总任务数量
	_, count := model.TaskGetList(1, 10)
	// self.Data["totalJob"] = count
	data["totalJob"] = count

	//日志总量
	_, totalLog := model.TaskLogGetList(1, 10)
	data["totalLog"] = totalLog

	//待审核任务数量
	_, totalAuditTask := model.TaskGetList(1, 10, "status", 2)
	data["totalAuditTask"] = totalAuditTask

	//失败
	errorNum, err := model.GetLogNum(-1)
	if err != nil {
		errorNum = 0
	}
	data["errorNum"] = errorNum

	//成功
	successNum, err := model.GetLogNum(0)
	if err != nil {
		successNum = 0
	}
	// self.Data["successNum"] = successNum
	data["successNum"] = successNum

	//用户数
	_, userNum := model.AdminGetList(1, 10, "status", 1)
	// self.Data["userNum"] = userNum
	data["userNum"] = userNum

	//累计运行总次数
	n, err := model.TaskTotalRunNum()
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
		task, _ := model.TaskGetById(job.GetId())
		row["task_id"] = job.GetId()
		row["task_name"] = job.GetName()
		row["task_group"] = groups_map[task.GroupId]
		row["next_time"] = beego.Date(v.Next, "Y-m-d H:i:s")
		jobList[k] = row
		startJob++
	}

	data["recentLogs"] = jobList

	// 最近执行失败的日志
	logs, _ := model.TaskLogGetList(1, 30, "status__lt", 0)
	errLogs := make([]map[string]interface{}, len(logs))

	for k, v := range logs {
		task, err := model.TaskGetById(v.TaskId)
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
	okRun := model.SumByDays(30, "0")
	errRun := model.SumByDays(30, "-1")
	expiredRun := model.SumByDays(30, "-2")

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
	info := libs.SystemInfo(model.StartTime)
	data["sysInfo"] = info

	data["pageTitle"] = "系统概况"
	// self.display()

	c.HTML(http.StatusOK, "start.html", data)
}
