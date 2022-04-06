/************************************************************
** @Description: controllers
** @Author: haodaquan
** @Date:   2018-06-11 21:11
** @Last Modified by:   haodaquan
** @Last Modified time: 2018-06-11 21:11
*************************************************************/
package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/app/service"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/libs"
	"github.com/voioc/cjob/utils"
	"github.com/voioc/cjob/worker"

	cron "github.com/voioc/cjob/crons"
)

type TaskController struct {
	BaseController
}

func (self *TaskController) List(c *gin.Context) {

	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	data["pageTitle"] = "任务管理"

	tg, _ := service.AuthS(c).TaskGroups(c.GetInt("uid"), c.GetString("role_ids"))

	data["taskGroup"], _ = service.TaskGroupS(c).GroupIDName(tg)
	data["groupId"] = 0
	// arr := strings.Split(self.Ctx.GetCookie("groupid"), "|")
	// if len(arr) > 0 {
	// 	 self.Data["groupId"], _ = strconv.Atoi(arr[0])
	// }
	// self.display()

	c.HTML(http.StatusOK, "task/list.html", data)
}

func (self *TaskController) AuditList(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	data["pageTitle"] = "任务审核"
	// self.display()
	c.HTML(http.StatusOK, "task/auditlist.html", data)
}

func (self *TaskController) Add(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	uid := c.GetInt("uid")
	tg, sg := service.AuthS(c).TaskGroups(uid, c.GetString("role_id"))
	data["pageTitle"] = "新增任务"
	data["taskGroup"], _ = service.TaskGroupS(c).GroupIDName(tg) // taskGroupLists(tg, uid)
	data["serverGroup"], _ = service.ServerS(c).ServerLists(sg)  // serverLists(tg, uid)
	data["isAdmin"] = uid
	data["adminInfo"], _ = service.AdminS(c).AdminInfo([]int{}) // AllAdminInfo("")
	// self.display()

	c.HTML(http.StatusOK, "task/add.html", data)
}

func (self *TaskController) Edit(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	data["pageTitle"] = "编辑任务"

	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	// err := service.TaskS(c).TaskByID(id) // model.TaskGetById(id)
	task := model.Task{}
	if err := model.DataByID(&task, id); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	if task.Status == 1 {
		// self.ajaxMsg("运行状态无法编辑任务，请先暂停任务", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "运行状态无法编辑任务，请先暂停任务"))
		return
	}

	data["task"] = task
	data["adminInfo"], _ = service.AdminS(c).AdminInfo(nil) // AllAdminInfo("")

	uid := c.GetInt("uid")
	tg, sg := service.AuthS(c).TaskGroups(uid, c.GetString("role_id"))

	// 分组列表
	data["taskGroup"], _ = service.TaskGroupS(c).GroupIDName(tg) // taskGroupLists(tg, uid)
	data["serverGroup"], _ = service.ServerS(c).ServerLists(sg)  // serverLists(sg, uid)
	data["isAdmin"] = uid

	var notifyUserIds []int
	if task.NotifyUserIDs != "0" {
		notifyUserIdsStr := strings.Split(task.NotifyUserIDs, ",")
		for _, v := range notifyUserIdsStr {
			i, _ := strconv.Atoi(v)
			notifyUserIds = append(notifyUserIds, i)
		}
	}

	data["notify_user_ids"] = notifyUserIds

	server_ids := strings.Split(task.ServerIDs, ",")
	var server_ids_arr []int
	for _, sv := range server_ids {
		i, _ := strconv.Atoi(sv)
		server_ids_arr = append(server_ids_arr, i)
	}

	data["service_ids"] = server_ids_arr

	filters := []interface{}{"status =", 1, "tpl_type =", task.NotifyType}
	// notifyTplList, err := service.NotifyS(c).NotifyTypeList(task.NotifyType) // model.NotifyTplGetByTplTypeList(task.NotifyType)
	notifyTplList := make([]model.NotifyTpl, 0)
	if err := model.List(&notifyTplList, 1, 1000, filters...); err != nil {
		fmt.Println(err.Error())
	}

	tplList := make([]map[string]interface{}, len(notifyTplList))
	for k, v := range notifyTplList {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["tpl_name"] = v.TplName
		row["tpl_type"] = v.TplType
		tplList[k] = row
	}

	data["notifyTpl"] = tplList

	// self.display()
	c.HTML(http.StatusOK, "task/edit.html", data)
}

func (self *TaskController) Copy(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "复制任务"
	data["adminInfo"], _ = service.AdminS(c).AdminInfo([]int{}) // AllAdminInfo("")

	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	// task, err := service.TaskS(c).TaskByID(id) // model.TaskGetById(id)
	task := model.Task{}
	if err := model.DataByID(&task, id); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	uid := c.GetInt("uid")

	//if task.Status == 1 {
	//	self.ajaxMsg("运行状态无法编辑任务，请先暂停任务", MSG_ERR)
	//}
	data["task"] = task

	data["adminInfo"], _ = service.AdminS(c).AdminInfo([]int{}) // AllAdminInfo("")

	// 分组列表
	tg, sg := service.AuthS(c).TaskGroups(uid, c.GetString("role_id"))
	data["taskGroup"], _ = service.TaskGroupS(c).GroupIDName(tg) // taskGroupLists(tg, uid)
	data["serverGroup"], _ = service.ServerS(c).ServerLists(sg)  // serverLists(sg, uid)
	data["isAdmin"] = uid

	var notifyUserIds []int
	if task.NotifyUserIDs != "0" {
		notifyUserIdsStr := strings.Split(task.NotifyUserIDs, ",")
		for _, v := range notifyUserIdsStr {
			i, _ := strconv.Atoi(v)
			notifyUserIds = append(notifyUserIds, i)
		}
	}

	data["notify_user_ids"] = notifyUserIds

	server_ids := strings.Split(task.ServerIDs, ",")
	var server_ids_arr []int
	for _, sv := range server_ids {
		i, _ := strconv.Atoi(sv)
		server_ids_arr = append(server_ids_arr, i)
	}

	data["service_ids"] = server_ids_arr

	// notifyTplList, err := service.NotifyS(c).NotifyTypeList(task.NotifyType) // model.NotifyTplGetByTplTypeList(task.NotifyType)
	notifyTplList := make([]model.NotifyTpl, 0)
	filters := []interface{}{"status =", 1, "tpl_type =", task.NotifyType}
	if err := model.List(&notifyTplList, 1, 1000, filters...); err != nil {
		fmt.Println(err.Error())
	}

	tplList := make([]map[string]interface{}, len(notifyTplList))
	for k, v := range notifyTplList {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["tpl_name"] = v.TplName
		row["tpl_type"] = v.TplType
		tplList[k] = row
	}

	data["notifyTpl"] = tplList

	// self.display()
	c.HTML(http.StatusOK, "task/copy.html", data)
}

// Detail dd
func (self *TaskController) Detail(c *gin.Context) {
	uid := c.GetInt("uid")
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	data["pageTitle"] = "任务详细"

	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	// task, err := service.TaskS(c).TaskByID(id) // model.TaskGetById(id)
	task := model.Task{}
	if err := model.DataByID(&task, id); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	TextStatus := []string{
		"<font color='red'><i class='fa fa-minus-square'></i> 暂停</font>",
		"<font color='green'><i class='fa fa-check-square'></i> 运行中</font>",
		"<font color='orange'><i class='fa fa-question-circle'></i> 待审核</font>",
		"<font color='red'><i class='fa fa-times-circle'></i> 审核失败</font>",
	}

	data["TextStatus"] = TextStatus[task.Status]
	data["CreateTime"] = time.Unix(task.CreatedAt, 0).Format("2006-01-02 15:04:05") // beego.Date(time.Unix(task.CreatedAt, 0), "Y-m-d H:i:s")
	data["UpdateTime"] = time.Unix(task.UpdatedAt, 0).Format("2006-01-02 15:04:05") // beego.Date(time.Unix(task.UpdatedAt, 0), "Y-m-d H:i:s")
	data["task"] = task

	tg, _ := service.AuthS(c).TaskGroups(uid, c.GetString("role_ids"))

	// 分组列表
	data["taskGroup"], _ = service.TaskGroupS(c).GroupIDName(tg) //taskGroupLists(tg, uid)

	serverName := ""
	if task.ServerIDs == "0" {
		serverName = "本地服务器 <br>"
	} else {
		serverIdSli := strings.Split(task.ServerIDs, ",")
		for _, v := range serverIdSli {
			if v == "0" {
				serverName = "本地服务器 <br>"
			}
		}

		sids := []int{}
		for _, row := range strings.Split(task.ServerIDs, ",") {
			sid, _ := strconv.Atoi(row)
			if sid != 0 {
				sids = append(sids, sid)
			}
		}

		servers := make([]*model.TaskServer, 0)
		if err := model.DataByIDs(&servers, sids); err != nil || len(servers) == 0 {
			if err != nil {
				fmt.Println(err.Error())
			}

			serverName += "服务器异常!!"
		} else {

			// servers, n := model.TaskServerGetByIds(task.ServerIDs)
			// if n > 0 {
			for _, server := range servers {
				// fmt.Println(server.Status)
				if server.Status != 0 {
					serverName += server.ServerName + " <i class='fa fa-ban' style='color:#FF5722'></i> <br/> "
				} else {
					serverName += server.ServerName + " <br/> "
				}
			}
			// } else {

			// }
		}
	}

	//执行策略
	data["ServerType"] = "同时执行"
	if task.ServerType == 1 {
		data["ServerType"] = "轮询执行"
	}

	// 任务分组
	groupName := "默认分组"
	if task.GroupID > 0 {
		// group, err := service.TaskGroupS(c).GroupByID(task.GroupID) //model.GroupGetById(task.GroupID)
		group := model.TaskGroup{}
		if err := model.DataByID(&group, task.GroupID); err == nil {
			groupName = group.GroupName
		}
	}

	data["GroupName"] = groupName

	// 创建人和修改人
	createName := "未知"
	updateName := "未知"
	if task.CreatedID > 0 {
		// admin, err := service.AdminS(c).AdminGetByID(task.CreatedID) // model.AdminGetById(task.CreatedID)
		admin := model.Admin{}
		if err := model.DataByID(&admin, task.CreatedID); err == nil {
			createName = admin.RealName
		}
	}

	if task.UpdatedID > 0 {
		// admin, err := service.AdminS(c).AdminGetByID(task.UpdatedID) // model.AdminGetById(task.UpdatedID)
		admin := model.Admin{}
		if err := model.DataByID(&admin, task.UpdatedID); err == nil {
			updateName = admin.RealName
		}
	}

	//是否出错通知
	data["adminInfo"] = []*model.Admin{}
	if task.NotifyUserIDs != "0" && task.NotifyUserIDs != "" {
		ids := []int{}
		for _, row := range strings.Split(task.NotifyUserIDs, ",") {
			id, _ := strconv.Atoi(row)
			if id != 0 {
				ids = append(ids, id)
			}
		}

		data["adminInfo"], _ = service.AdminS(c).AdminInfo(ids) // AllAdminInfo(task.NotifyUserIds)
	}
	data["CreateName"] = createName
	data["UpdateName"] = updateName
	data["serverName"] = serverName

	data["NotifyTplName"] = "未知"
	if task.IsNotify == 1 {
		notifyTpl, err := service.NotifyS(c).NotifyListIDs([]int{task.NotifyTplID}) // model.NotifyTplGetById(task.NotifyTplID)
		if err == nil && len(notifyTpl) > 0 {
			data["NotifyTplName"] = notifyTpl[0].TplName
		}
	}

	c.HTML(http.StatusOK, "task/detail.html", data)
}

func (self *TaskController) Save(c *gin.Context) {
	uid := c.GetInt("uid")
	taskID, _ := strconv.Atoi(c.DefaultPostForm("id", ""))

	command, err := service.BanS(c).CheckCommand(c.DefaultPostForm("command", ""))
	if err != nil {
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	if command != "" {
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "含有禁止命令："+command))
		return
	}

	if taskID == 0 {
		task := new(model.Task)
		task.CreatedID = c.GetInt("uid")
		task.GroupID, _ = strconv.Atoi(c.DefaultPostForm("group_id", "0"))
		task.TaskName = strings.TrimSpace(c.DefaultPostForm("task_name", ""))
		task.Description = strings.TrimSpace(c.DefaultPostForm("description", ""))
		task.Concurrent, _ = strconv.Atoi(c.DefaultPostForm("concurrent", "0"))
		task.ServerIDs = strings.TrimSpace(c.DefaultPostForm("server_ids", ""))
		task.CronSpec = strings.TrimSpace(c.DefaultPostForm("cron_spec", ""))
		task.Command = strings.TrimSpace(c.DefaultPostForm("command", ""))
		task.Timeout, _ = strconv.Atoi(c.DefaultPostForm("timeout", "0"))
		task.IsNotify, _ = strconv.Atoi(c.DefaultPostForm("is_notify", "0"))
		task.ServerType, _ = strconv.Atoi(c.DefaultPostForm("server_type", "0"))

		task.NotifyType, _ = strconv.Atoi(c.DefaultPostForm("notify_type", "0"))
		task.NotifyTplID, _ = strconv.Atoi(c.DefaultPostForm("notify_tpl_id", "0"))
		task.NotifyUserIDs = strings.TrimSpace(c.DefaultPostForm("notify_user_ids", ""))

		if task.IsNotify == 1 && task.NotifyTplID <= 0 {
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "请选择通知模板"))
			// self.ajaxMsg("请选择通知模板", MSG_ERR)
			return
		}

		task.CreatedAt = time.Now().Unix()
		task.UpdatedAt = time.Now().Unix()
		task.Status = 2 //审核中

		if uid == 1 {
			task.Status = 0 //审核中,超级管理员不需要
		}

		if task.TaskName == "" || task.CronSpec == "" || task.Command == "" {
			// self.ajaxMsg("", MSG_ERR)
			// fmt.Println("11111")
			// fmt.Println(task.TaskName, task.CronSpec, task.Command)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "请填写完整信息!"))
			return
		}

		if _, err := cron.Parse(task.CronSpec); err != nil {
			// self.ajaxMsg("cron表达式无效", MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "cron表达式无效"))
			return
		}

		msg := ""
		if task.TaskName == "" {
			msg = "任务名称不能为空"
		}

		if task.CronSpec == "" {
			msg = "时间表达式不能为空"
		}
		if task.Command == "" {
			msg = "命令内容不能为空"
		}

		if msg != "" {
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, msg))
			return
		}

		if task.CreatedAt == 0 {
			task.CreatedAt = time.Now().Unix()
		}

		if err := model.Add(task); err != nil { // model.TaskAdd(task); err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}

		// self.ajaxMsg("", MSG_OK)
		c.JSON(http.StatusOK, common.Success(c))
		return
	}

	// task, _ := service.TaskS(c).TaskByID(taskID) // model.TaskGetById(task_id)
	task := model.Task{}
	if err := model.DataByID(&task, taskID); err != nil {
		fmt.Println(err.Error())
	}

	// 修改
	task.ID = taskID
	task.UpdatedAt = time.Now().Unix()
	task.TaskName = strings.TrimSpace(c.DefaultPostForm("task_name", ""))
	task.Description = strings.TrimSpace(c.DefaultPostForm("description", ""))
	task.GroupID, _ = strconv.Atoi(c.DefaultPostForm("group_id", "0"))
	task.Concurrent, _ = strconv.Atoi(c.DefaultPostForm("concurrent", "0"))
	task.ServerIDs = strings.TrimSpace(c.DefaultPostForm("server_ids", ""))
	task.CronSpec = strings.TrimSpace(c.DefaultPostForm("cron_spec", ""))
	task.Command = strings.TrimSpace(c.DefaultPostForm("command", ""))
	task.Timeout, _ = strconv.Atoi(c.DefaultPostForm("timeout", "0"))
	task.ServerType, _ = strconv.Atoi(c.DefaultPostForm("server_type", "0"))
	task.IsNotify, _ = strconv.Atoi(c.DefaultPostForm("is_notify", "0"))
	task.NotifyType, _ = strconv.Atoi(c.DefaultPostForm("notify_type", "0"))
	task.NotifyTplID, _ = strconv.Atoi(c.DefaultPostForm("notify_tpl_id", "0"))
	task.NotifyUserIDs = strings.TrimSpace(c.DefaultPostForm("notify_user_ids", ""))
	task.UpdatedID = uid
	task.Status = 2 //审核中,超级管理员不需要

	if task.IsNotify == 1 && task.NotifyTplID <= 0 {
		// self.ajaxMsg("请选择通知模板", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "请选择通知模板"))
		return
	}

	if uid == 1 {
		task.Status = 0
	}

	if _, err := cron.Parse(task.CronSpec); err != nil {
		// self.ajaxMsg("cron表达式无效", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "cron表达式无效"))
		return
	}

	if err := model.Update(task.ID, task); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Success(c))
}

// //检查是否含有禁用命令
// func checkCommand(command string) (string, bool) {

// 	filters := make([]interface{}, 0)
// 	filters = append(filters, "status", 0)
// 	ban, _ := model.BanGetList(1, 1000, filters...)
// 	for _, v := range ban {
// 		if strings.Contains(command, v.Code) {
// 			return v.Code, false
// 		}
// 	}
// 	return "", true
// }

func (self *TaskController) Audit(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.DefaultPostForm("id", ""))

	// taskId, _ := self.GetInt("id")
	if taskID == 0 {
		// self.ajaxMsg("任务不存在", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "任务不存在"))
		return
	}

	task := &model.Task{
		Status:    0,
		UpdatedID: c.GetInt("uid"),
		UpdatedAt: time.Now().Unix(),
	}

	if err := model.Update(task.ID, task, true); err != nil { // changeStatus(taskID, 0, c.GetInt("uid"))
		// self.ajaxMsg("审核失败", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "审核失败"))
		return
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *TaskController) AjaxNopass(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.DefaultPostForm("id", ""))
	if taskID == 0 {
		// self.ajaxMsg("任务不存在", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "任务不存在"))
		return
	}

	task := &model.Task{
		ID:        taskID,
		Status:    3,
		UpdatedID: c.GetInt("uid"),
		UpdatedAt: time.Now().Unix(),
	}

	if err := model.Update(task.ID, task); err != nil { //changeStatus(taskID, 3, c.GetInt("uid"))
		// if !res {
		// self.ajaxMsg("操作失败", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "操作失败"))
		return
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *TaskController) AjaxStart(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.DefaultPostForm("id", ""))
	if taskID == 0 {
		// self.ajaxMsg("任务不存在", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "任务不存在"))
		return
	}

	// task, err := service.TaskS(c).TaskByID(taskID) // model.TaskGetById(taskID)
	task := model.Task{}
	if err := model.DataByID(&task, taskID); err != nil {
		// self.ajaxMsg("查不到该任务", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "查不到该任务"))
		return
	}

	if task.Status != 0 {
		msg := "任务状态有误"
		if task.Status == 2 {
			msg = "任务正在审核中,不能启动"
		}

		// self.ajaxMsg("任务状态有误", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, msg))
		return
	}

	jobArr, err := service.TaskS(c).CreateJob(&task) // worker.NewJobFromTask(&task)

	if err != nil {
		// self.ajaxMsg("创建任务失败", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "创建任务失败"))
		return
	}

	for _, job := range jobArr {
		if worker.AddJob(task.CronSpec, job) {
			task.Status = 1
			if err := model.Update(taskID, task); err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *TaskController) AjaxPause(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.DefaultPostForm("id", ""))
	if taskID == 0 {
		// self.ajaxMsg("任务不存在", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "任务不存在"))
		return
	}

	// task, err := service.TaskS(c).TaskByID(taskID) // model.TaskGetById(taskID)
	task := model.Task{}
	if err := model.DataByID(&task, taskID); err != nil {
		// self.ajaxMsg("查不到该任务", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "查不到该任务"))
		return
	}

	//移出任务
	TaskServerIdsArr := strings.Split(task.ServerIDs, ",")
	for _, server_id := range TaskServerIdsArr {
		server_id_int, _ := strconv.Atoi(server_id)
		jobKey := libs.JobKey(task.ID, server_id_int)
		worker.RemoveJob(jobKey)
	}

	task.Status = 0
	if err := model.Update(task.ID, &task, true); err != nil {
		fmt.Println(err.Error())
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

// AjaxRun ss
func (self *TaskController) AjaxRun(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	// task, err := service.TaskS(c).TaskByID(id) // model.TaskGetById(id)
	task := model.Task{}
	if err := model.DataByID(&task, id); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	jobArr, err := service.TaskS(c).CreateJob(&task) //worker.NewJobFromTask(&task)
	if err != nil {
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	for _, job := range jobArr {
		job.Run()
	}

	c.JSON(http.StatusOK, common.Success(c))
}

// AjaxBatchStart sdf
func (self *TaskController) AjaxBatchStart(c *gin.Context) {
	ids := strings.Split(c.PostForm("ids"), ",")
	// ids := strings.Split(idStr, ",")
	if len(ids) < 1 {
		// self.ajaxMsg("请选择要操作的任务", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "请选择要操作的任务"))
		return
	}

	for _, v := range ids {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}

		// if task, err := service.TaskS(c).TaskByID(id); err == nil { // model.TaskGetById(id); err == nil {
		task := model.Task{}
		if err := model.DataByID(&task, id); err != nil {
			fmt.Println(err.Error())
		} else {
			jobArr, err := service.TaskS(c).CreateJob(&task) // worker.NewJobFromTask(&task)
			if err == nil {
				for _, job := range jobArr {
					worker.AddJob(task.CronSpec, job)
				}

				task.Status = 1
				// task.Update()
				if err := model.Update(task.ID, &task); err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

// AjaxBatchPause kkk
func (self *TaskController) AjaxBatchPause(c *gin.Context) {
	ids := strings.Split(c.PostForm("ids"), ",")
	// ids := strings.Split(idStr, ",")
	if len(ids) < 2 && ids[0] == "" {
		// self.ajaxMsg("请选择要操作的任务", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "请选择要暂停的任务"))
		return
	}

	for _, v := range ids {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}

		// task, err := service.TaskS(c).TaskByID(id) // model.TaskGetById(id)
		task := model.Task{}
		if err := model.DataByID(&task, id); err != nil {
			fmt.Println(err.Error())
		} else {

			// 移出任务
			TaskServerIdsArr := strings.Split(task.ServerIDs, ",")
			for _, server_id := range TaskServerIdsArr {
				server_id_int, _ := strconv.Atoi(server_id)
				jobKey := libs.JobKey(task.ID, server_id_int)
				worker.RemoveJob(jobKey)
			}

			task.Status = 0
			if err := model.Update(task.ID, task); err != nil {
				fmt.Println(err.Error())
			}
			// task.Update()
		}
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

// AjaxBatchDel kddd
func (self *TaskController) AjaxBatchDel(c *gin.Context) {
	ids := strings.Split(c.PostForm("ids"), ",")
	// ids := strings.Split(idStr, ",")
	if len(ids) < 1 {
		// self.ajaxMsg("请选择要操作的任务", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "请选择要删除的任务"))
		return
	}

	for _, v := range ids {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}

		// task, _ := service.TaskS(c).TaskByID(id) // model.TaskGetById(id)
		task := model.Task{}
		if err := model.DataByID(&task, id); err != nil {
			fmt.Println(err.Error())
		}

		//移出任务
		TaskServerIdsArr := strings.Split(task.ServerIDs, ",")

		for _, server_id := range TaskServerIdsArr {
			server_id_int, _ := strconv.Atoi(server_id)
			jobKey := libs.JobKey(task.ID, server_id_int)
			worker.RemoveJob(jobKey)
		}

		// service.TaskS(c).Del([]int{id})
		if err := model.Del(&model.Task{}, id); err != nil {
			fmt.Println(err.Error())
		}

		service.TaskLogS(c).LogDelTaskID([]int{id})
		// model.TaskDel(id)
		// model.TaskLogDelByTaskId(id)
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *TaskController) AjaxBatchAudit(c *gin.Context) {
	ids := strings.Split(c.PostForm("ids"), ",")
	// ids := strings.Split(idStr, ",")
	if len(ids) < 1 {
		// self.ajaxMsg("请选择要操作的任务", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "请选择要审核的任务"))
		return
	}

	for _, v := range ids {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}

		task := &model.Task{
			Status:    0,
			UpdatedID: c.GetInt("uid"),
			UpdatedAt: time.Now().Unix(),
		}

		// if err := service.TaskS(c).Update(task, true); err != nil { // changeStatus(taskID, 0, c.GetInt("uid"))
		// 	// changeStatus(id, 0, self.userId)
		// }
		if err := model.Update(task.ID, task, true); err != nil {
			fmt.Println(err.Error())
		}
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *TaskController) Reject(c *gin.Context) {
	ids := strings.Split(c.PostForm("ids"), ",")
	// ids := strings.Split(idStr, ",")
	if len(ids) < 1 {
		// self.ajaxMsg("请选择要操作的任务", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "请选择要操作的任务"))
		return
	}

	for _, v := range ids {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}

		task := &model.Task{
			Status:    0,
			UpdatedID: c.GetInt("uid"),
			UpdatedAt: time.Now().Unix(),
		}
		// if err := service.TaskS(c).Update(task, true); err != nil { // changeStatus(taskID, 0, c.GetInt("uid"))
		// 	// changeStatus(id, 3, self.userId)
		// }
		if err := model.Update(task.ID, task, true); err != nil {
			fmt.Println(err.Error())
		}
	}

	c.JSON(http.StatusOK, common.Success(c))
}

// func changeStatus(taskId, status, userId int) bool {
// 	if taskId == 0 {
// 		return false
// 	}

// 	task, _ := model.TaskGetById(taskId)
// 	//修改
// 	task.ID = taskId
// 	task.UpdatedAt = time.Now().Unix()
// 	task.UpdatedID = userId
// 	task.Status = status //0,1,2,3,9

// 	if err := task.Update(); err != nil {
// 		return false
// 	}
// 	return true
// }

// AjaxDel ddd
func (self *TaskController) AjaxDel(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	// task, _ := service.TaskS(c).TaskByID(id) // model.TaskGetById(id)
	task := model.Task{}
	if err := model.DataByID(&task, id); err != nil {
		fmt.Println(err.Error())
	}

	uid := c.GetInt("uid")
	task.UpdatedAt = time.Now().Unix()
	task.UpdatedID = uid
	task.Status = -1
	task.ID = id

	//TODO 查询服务器是否用于定时任务
	if err := model.Update(task.ID, task); err != nil { // task.Update(); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *TaskController) AjaxNotifyType(c *gin.Context) {
	notifyType, _ := strconv.Atoi(c.DefaultPostForm("notify_type", "0"))
	// result, _ := service.NotifyS(c).NotifyTypeList(notifyType) // model.NotifyTplGetByTplTypeList(notifyType)
	result := make([]model.NotifyTpl, 0)
	filters := []interface{}{"status =", 1, "tpl_type =", notifyType}
	if err := model.List(&result, 1, 1000, filters...); err != nil {
		fmt.Println(err.Error())
	}

	list := make([]map[string]interface{}, 0)
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["tpl_name"] = v.TplName
		row["tpl_type"] = v.TplType
		list[k] = row
	}

	// self.ajaxList("成功", MSG_OK, count, list)
	ext := map[string]int{"total": len(list)}
	c.JSON(http.StatusOK, common.Success(c, list, ext))
}

func (self *TaskController) Table(c *gin.Context) {

	pagesize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	// 列表
	// page, err := c.DefaultQuery("page")
	// if err != nil {
	// 	page = 1
	// }

	// limit, err := self.GetInt("limit")
	// if err != nil {
	// 	limit = 30
	// }

	// groupId, _ := self.GetInt("group_id", 0)
	groupId := 0

	//0-全部，-1如果存在，n,如果不存在，0

	//if groupId == -1 {
	//	groupId = 0
	//	arr := strings.Split(self.Ctx.GetCookie("groupid"), "|")
	//	if len(arr) > 0 {
	//		groupId, _ = strconv.Atoi(arr[0])
	//	}
	//}

	//if groupId > 0 {
	//	self.Ctx.SetCookie("groupid", strconv.Itoa(groupId)+"|job")
	//}

	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))
	taskName := strings.TrimSpace(c.DefaultQuery("task_name", ""))

	StatusText := []string{
		"<font color='red'><i class='fa fa-minus-square'></i></font>",
		"<font color='green'><i class='fa fa-check-square'></i></font>",
		"<font color='orange'><i class='fa fa-question-circle'></i></font>",
		"<font color='red'><i class='fa fa-times-circle'></i></font>",
	}

	uid := c.GetInt("uid")
	tg, _ := service.AuthS(c).TaskGroups(uid, c.GetString("role_id"))
	taskGroup, _ := service.TaskGroupS(c).GroupIDName(tg) //taskGroupLists(taskGroups, uid)
	// self.pageSize = pagesize

	// 查询条件
	filters := make([]interface{}, 0)

	if status == 2 {
		// 审核中，审核失败
		filters = append(filters, "status", []int{2, 3})
	} else {
		filters = append(filters, "status", []int{0, 1})
	}

	// 搜索全部
	if groupId == 0 {
		if uid != 1 {
			groups := strings.Split(tg, ",")
			groupsIds := make([]int, 0)
			for _, v := range groups {
				id, _ := strconv.Atoi(v)
				groupsIds = append(groupsIds, id)
			}

			filters = append(filters, "group_id", groupsIds)
		}
	} else if groupId > 0 {
		filters = append(filters, "group_id", groupId)
	}

	if taskName != "" {
		filters = append(filters, "task_name LIKE %"+taskName+"%", " ")
	}

	filters = append(filters, "order", "field(status, 1, 2, 3, 0), id desc")
	// result, count, _ := service.TaskS(c).TaskGetList(page, pagesize, filters...) // model.TaskGetList(page, pagesize, filters...)
	result := make([]model.Task, 0)
	if err := model.List(&result, page, pagesize, filters...); err != nil {
		fmt.Println(err.Error())
	}

	count, err := model.ListCount(&model.Task{}, filters...)
	if err != nil {
		fmt.Println(err.Error())
	}

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID

		groupName := "默认分组"

		if name, ok := taskGroup[v.GroupID]; ok {
			groupName = name
		}

		row["group_name"] = groupName
		row["task_name"] = StatusText[v.Status] + " " + groupName + "-" + "&nbsp;" + v.TaskName
		row["description"] = v.Description

		//row["status_text"] = StatusText[v.Status]
		row["status"] = v.Status
		row["pre_time"] = time.Unix(v.PrevTime, 0).Format("2006-01-02 15:04:05") // beego.Date(time.Unix(v.PrevTime, 0), "Y-m-d H:i:s")
		row["execute_times"] = v.ExecuteTimes
		row["cron_spec"] = v.CronSpec

		TaskServerIDsArr := strings.Split(v.ServerIDs, ",")
		serverID := 0
		if len(TaskServerIDsArr) > 0 {
			serverID, _ = strconv.Atoi(TaskServerIDsArr[0])
		}

		jobskey := libs.JobKey(v.ID, serverID)
		e, _ := worker.GetEntryByID(jobskey)

		if e != nil {
			row["next_time"] = e.Next.Format("2006-01-02 15:04:05")
			row["prev_time"] = "-"
			if e.Prev.Unix() > 0 {
				row["prev_time"] = e.Prev.Format("2006-01-02 15:04:05")
			} else if v.PrevTime > 0 {
				row["prev_time"] = time.Unix(v.PrevTime, 0).Format("2006-01-02 :15:04:05")
			}
			row["running"] = 1
		} else {
			row["next_time"] = "-"
			if v.PrevTime > 0 {
				row["prev_time"] = time.Unix(v.PrevTime, 0).Format("2006-01-02 15:04:05")
			} else {
				row["prev_time"] = "-"
			}
			row["running"] = 0
		}

		list[k] = row
	}

	// self.ajaxList("成功", MSG_OK, count, list)
	ext := map[string]int{"total": int(count)}
	c.JSON(http.StatusOK, common.Success(c, list, ext))
}

func (self *TaskController) ApiTask(c *gin.Context) {
	// uid := c.GetInt("uid")
	taskID, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	if taskID == 0 {
		task := new(model.Task)
		task.CreatedID, _ = strconv.Atoi(c.DefaultPostForm("create_id", "0"))
		task.GroupID, _ = strconv.Atoi(c.DefaultPostForm("group_id", "0"))
		task.TaskName = strings.TrimSpace(c.DefaultPostForm("task_name", ""))
		task.Description = strings.TrimSpace(c.DefaultPostForm("description", ""))
		task.Concurrent, _ = strconv.Atoi(c.DefaultPostForm("concurrent", "0"))
		task.ServerIDs = strings.TrimSpace(c.DefaultPostForm("server_ids", ""))
		task.CronSpec = strings.TrimSpace(c.DefaultPostForm("cron_spec", ""))
		task.Command = strings.TrimSpace(c.DefaultPostForm("command", ""))
		task.Timeout, _ = strconv.Atoi(c.DefaultPostForm("timeout", "0"))
		task.IsNotify, _ = strconv.Atoi(c.DefaultPostForm("is_notify", "0"))
		task.NotifyType, _ = strconv.Atoi(c.DefaultPostForm("notify_type", "0"))
		task.NotifyTplID, _ = strconv.Atoi(c.DefaultPostForm("notify_tpl_id", "0"))
		task.NotifyUserIDs = strings.TrimSpace(c.DefaultPostForm("notify_user_ids", ""))

		if task.IsNotify == 1 && task.NotifyTplID <= 0 {
			// self.ajaxMsg("请选择通知模板", MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "请选择通知模板"))
			return
		}

		command, err2 := service.BanS(c).CheckCommand(task.Command) // checkCommand(task.Command)
		if err2 != nil {
			// self.ajaxMsg("含有禁止命令："+msg, MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err2.Error()))
			return
		}

		if command != "" {
			// self.ajaxMsg("含有禁止命令："+msg, MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "含有禁止命令: "+command))
			return
		}

		task.CreatedAt = time.Now().Unix()
		task.UpdatedAt = time.Now().Unix()
		task.Status = 0 //接口不需要审核

		if task.TaskName == "" || task.CronSpec == "" || task.Command == "" {
			// self.ajaxMsg("请填写完整信息", MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "请填写完整信息"))
			return
		}

		var id int64
		var err error
		if _, err = cron.Parse(task.CronSpec); err != nil {
			// self.ajaxMsg("cron表达式无效", MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "cron表达式无效"))
			return
		}

		if err = model.Add(&task); err != nil { // model.TaskAdd(task); err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}

		// taskID = int(id)
		// self.ajaxMsg(task_id, MSG_OK)
		c.JSON(http.StatusOK, common.Success(c, map[string]int64{"task_id": id}))
		return
	}

	// task, _ := service.TaskS(c).TaskByID(taskID) // model.TaskGetById(taskID)
	task := model.Task{}
	if err := model.DataByID(&task, taskID); err != nil {
		fmt.Println(err.Error())
	}

	if task.Status == 1 {
		// self.ajaxMsg("运行状态无法编辑任务，请先暂停任务", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "运行状态无法编辑任务，请先暂停任务"))
		return
	}

	//修改
	task.ID = taskID
	task.UpdatedAt = time.Now().Unix()
	task.TaskName = strings.TrimSpace(c.PostForm("task_name"))
	task.Description = strings.TrimSpace(c.PostForm("description"))
	task.GroupID, _ = strconv.Atoi(c.DefaultPostForm("group_id", "0"))
	task.Concurrent, _ = strconv.Atoi(c.DefaultPostForm("concurrent", "0"))
	task.ServerIDs = strings.TrimSpace(c.PostForm("server_ids"))
	task.CronSpec = strings.TrimSpace(c.PostForm("cron_spec"))
	task.Command = strings.TrimSpace(c.PostForm("command"))
	task.Timeout, _ = strconv.Atoi(c.DefaultPostForm("timeout", "0"))
	task.IsNotify, _ = strconv.Atoi(c.DefaultPostForm("is_notify", "0"))
	task.NotifyType, _ = strconv.Atoi(c.DefaultPostForm("notify_type", "0"))
	task.NotifyTplID, _ = strconv.Atoi(c.DefaultPostForm("notify_tpl_id", "0"))
	task.NotifyUserIDs = strings.TrimSpace(c.PostForm("notify_user_ids"))
	task.UpdatedID, _ = strconv.Atoi(c.DefaultPostForm("update_id", "0"))
	task.Status = 0 //接口不需要

	if task.IsNotify == 1 && task.NotifyTplID <= 0 {
		// self.ajaxMsg("请选择通知模板", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "请选择通知模板"))
		return
	}

	command, err2 := service.BanS(c).CheckCommand(task.Command) // checkCommand(task.Command)
	if err2 != nil {
		// self.ajaxMsg("含有禁止命令："+msg, MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err2.Error()))
		return
	}

	if command != "" {
		// self.ajaxMsg("含有禁止命令："+msg, MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "含有禁止命令: "+command))
		return
	}

	if _, err := cron.Parse(task.CronSpec); err != nil {
		// self.ajaxMsg("cron表达式无效", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "cron表达式无效"))
		return
	}

	if err := model.Update(task.ID, &task); err != nil { // task.Update(); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	// self.ajaxMsg(task_id, MSG_OK)
	c.JSON(http.StatusOK, common.Success(c, map[string]int{"task_id": taskID}))
	return
}

func (self *TaskController) ApiStart(c *gin.Context) {
	// taskId, _ := self.GetInt("id")
	taskID, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	if taskID == 0 {
		// self.ajaxMsg("任务不存在", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "任务不存在"))
		return
	}

	// task, err := service.TaskS(c).TaskByID(taskID) // model.TaskGetById(taskID)
	task := model.Task{}
	if err := model.DataByID(&task, taskID); err != nil {
		// self.ajaxMsg("查不到该任务", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "查不到该任务"))
		return
	}

	if task.Status != 0 {
		// self.ajaxMsg("任务状态有误", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "任务状态有误"))
		return
	}

	// 创建定时Job
	jobs, err := service.TaskS(c).CreateJob(&task)
	if err != nil {
		// self.ajaxMsg("创建任务失败", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "创建任务失败"))
		return
	}

	// 开启任务
	for _, job := range jobs {
		if worker.AddJob(task.CronSpec, job) {
			task.Status = 1
			// task.Update()
			if err := model.Update(task.ID, &task); err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *TaskController) ApiPause(c *gin.Context) {
	// taskId, _ := self.GetInt("id")
	taskID, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	if taskID == 0 {
		// self.ajaxMsg("任务不存在", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "任务不存在"))
		return
	}

	// task, err := service.TaskS(c).TaskByID(taskID) // model.TaskGetById(taskID)
	task := model.Task{}
	if err := model.DataByID(&task, taskID); err != nil {
		// self.ajaxMsg("查不到该任务", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "查不到该任务"))
		return
	}

	//移出任务
	TaskServerIdsArr := strings.Split(task.ServerIDs, ",")

	for _, server_id := range TaskServerIdsArr {
		server_id_int, _ := strconv.Atoi(server_id)
		jobKey := libs.JobKey(task.ID, server_id_int)
		worker.RemoveJob(jobKey)
	}

	task.Status = 0
	// task.Update()
	if err := model.Update(task.ID, &task); err != nil {
		fmt.Println(err.Error())
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}
