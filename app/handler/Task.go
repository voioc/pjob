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

	"github.com/astaxie/beego"
	cron "github.com/voioc/cjob/crons"
	"github.com/voioc/cjob/jobs"
)

type TaskController struct {
	BaseController
}

func (self *TaskController) List(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	data["pageTitle"] = "任务管理"
	data["taskGroup"] = taskGroupLists(self.taskGroups, self.userId)
	data["groupId"] = 0
	//arr := strings.Split(self.Ctx.GetCookie("groupid"), "|")
	//if len(arr) > 0 {
	//	self.Data["groupId"], _ = strconv.Atoi(arr[0])
	//}
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
	tg, _ := service.TaskGroups(uid, c.GetString("role_id"))
	data["pageTitle"] = "新增任务"
	data["taskGroup"] = taskGroupLists(tg, uid)
	data["serverGroup"] = serverLists(tg, uid)
	data["isAdmin"] = uid
	data["adminInfo"] = AllAdminInfo("")
	// self.display()

	c.HTML(http.StatusOK, "task/add.html", data)
}

func (self *TaskController) Edit(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	data["pageTitle"] = "编辑任务"

	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	task, err := model.TaskGetById(id)
	if err != nil {
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
	data["adminInfo"] = AllAdminInfo("")

	uid := c.GetInt("uid")
	tg, sg := service.TaskGroups(uid, c.GetString("role_id"))
	// 分组列表
	data["taskGroup"] = taskGroupLists(tg, uid)
	data["serverGroup"] = serverLists(sg, uid)
	data["isAdmin"] = uid

	var notifyUserIds []int
	if task.NotifyUserIds != "0" {
		notifyUserIdsStr := strings.Split(task.NotifyUserIds, ",")
		for _, v := range notifyUserIdsStr {
			i, _ := strconv.Atoi(v)
			notifyUserIds = append(notifyUserIds, i)
		}
	}

	data["notify_user_ids"] = notifyUserIds

	server_ids := strings.Split(task.ServerIds, ",")
	var server_ids_arr []int
	for _, sv := range server_ids {
		i, _ := strconv.Atoi(sv)
		server_ids_arr = append(server_ids_arr, i)
	}

	data["service_ids"] = server_ids_arr

	notifyTplList, _, err := model.NotifyTplGetByTplTypeList(task.NotifyType)
	tplList := make([]map[string]interface{}, len(notifyTplList))

	if err == nil {
		for k, v := range notifyTplList {
			row := make(map[string]interface{})
			row["id"] = v.Id
			row["tpl_name"] = v.TplName
			row["tpl_type"] = v.TplType
			tplList[k] = row
		}
	}

	data["notifyTpl"] = tplList

	// self.display()
	c.HTML(http.StatusOK, "task/edit.html", data)
}

func (self *TaskController) Copy(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "复制任务"
	data["adminInfo"] = AllAdminInfo("")

	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	task, err := model.TaskGetById(id)
	if err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	uid := c.GetInt("uid")

	//if task.Status == 1 {
	//	self.ajaxMsg("运行状态无法编辑任务，请先暂停任务", MSG_ERR)
	//}
	data["task"] = task

	data["adminInfo"] = AllAdminInfo("")

	// 分组列表
	tg, sg := service.TaskGroups(uid, c.GetString("role_id"))
	data["taskGroup"] = taskGroupLists(tg, uid)
	data["serverGroup"] = serverLists(sg, uid)
	data["isAdmin"] = uid

	var notifyUserIds []int
	if task.NotifyUserIds != "0" {
		notifyUserIdsStr := strings.Split(task.NotifyUserIds, ",")
		for _, v := range notifyUserIdsStr {
			i, _ := strconv.Atoi(v)
			notifyUserIds = append(notifyUserIds, i)
		}
	}

	data["notify_user_ids"] = notifyUserIds

	server_ids := strings.Split(task.ServerIds, ",")
	var server_ids_arr []int
	for _, sv := range server_ids {
		i, _ := strconv.Atoi(sv)
		server_ids_arr = append(server_ids_arr, i)
	}

	data["service_ids"] = server_ids_arr

	notifyTplList, _, err := model.NotifyTplGetByTplTypeList(task.NotifyType)
	tplList := make([]map[string]interface{}, len(notifyTplList))

	if err == nil {
		for k, v := range notifyTplList {
			row := make(map[string]interface{})
			row["id"] = v.Id
			row["tpl_name"] = v.TplName
			row["tpl_type"] = v.TplType
			tplList[k] = row
		}
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
	task, err := model.TaskGetById(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		// self.ajaxMsg(err.Error(), MSG_ERR)
		return
	}

	TextStatus := []string{
		"<font color='red'><i class='fa fa-minus-square'></i> 暂停</font>",
		"<font color='green'><i class='fa fa-check-square'></i> 运行中</font>",
		"<font color='orange'><i class='fa fa-question-circle'></i> 待审核</font>",
		"<font color='red'><i class='fa fa-times-circle'></i> 审核失败</font>",
	}

	data["TextStatus"] = TextStatus[task.Status]
	data["CreateTime"] = beego.Date(time.Unix(task.CreateTime, 0), "Y-m-d H:i:s")
	data["UpdateTime"] = beego.Date(time.Unix(task.UpdateTime, 0), "Y-m-d H:i:s")
	data["task"] = task

	tg, _ := service.TaskGroups(uid, c.GetString("role_ids"))

	// 分组列表
	data["taskGroup"] = taskGroupLists(tg, uid)

	serverName := ""
	if task.ServerIds == "0" {
		serverName = "本地服务器 <br>"
	} else {
		serverIdSli := strings.Split(task.ServerIds, ",")
		for _, v := range serverIdSli {
			if v == "0" {
				serverName = "本地服务器 <br>"
			}
		}
		servers, n := model.TaskServerGetByIds(task.ServerIds)
		if n > 0 {
			for _, server := range servers {
				fmt.Println(server.Status)
				if server.Status != 0 {
					serverName += server.ServerName + " <i class='fa fa-ban' style='color:#FF5722'></i> <br/> "
				} else {
					serverName += server.ServerName + " <br/> "
				}
			}
		} else {
			serverName += "服务器异常!!"
		}
	}

	//执行策略
	data["ServerType"] = "同时执行"
	if task.ServerType == 1 {
		data["ServerType"] = "轮询执行"
	}

	//任务分组
	groupName := "默认分组"
	if task.GroupId > 0 {
		group, err := model.GroupGetById(task.GroupId)
		if err == nil {
			groupName = group.GroupName
		}
	}

	data["GroupName"] = groupName

	//创建人和修改人
	createName := "未知"
	updateName := "未知"
	if task.CreateId > 0 {
		admin, err := model.AdminGetById(task.CreateId)
		if err == nil {
			createName = admin.RealName
		}
	}

	if task.UpdateId > 0 {
		admin, err := model.AdminGetById(task.UpdateId)
		if err == nil {
			updateName = admin.RealName
		}
	}

	//是否出错通知
	data["adminInfo"] = []*AdminInfo{}
	if task.NotifyUserIds != "0" && task.NotifyUserIds != "" {
		data["adminInfo"] = AllAdminInfo(task.NotifyUserIds)
	}
	data["CreateName"] = createName
	data["UpdateName"] = updateName
	data["serverName"] = serverName

	data["NotifyTplName"] = "未知"
	if task.IsNotify == 1 {
		notifyTpl, err := model.NotifyTplGetById(task.NotifyTplId)
		if err == nil {
			self.Data["NotifyTplName"] = notifyTpl.TplName
		}
	}

	c.HTML(http.StatusOK, "task/detail.html", data)
}

func (self *TaskController) Save(c *gin.Context) {
	uid := c.GetInt("uid")
	task_id, _ := strconv.Atoi(c.DefaultPostForm("id", ""))
	if task_id == 0 {
		task := new(model.Task)
		task.CreateId = c.GetInt("uid")
		task.GroupId, _ = strconv.Atoi(c.DefaultPostForm("group_id", "0"))
		task.TaskName = strings.TrimSpace(c.DefaultPostForm("task_name", ""))
		task.Description = strings.TrimSpace(c.DefaultPostForm("description", ""))
		task.Concurrent, _ = strconv.Atoi(c.DefaultPostForm("concurrent", "0"))
		task.ServerIds = strings.TrimSpace(c.DefaultPostForm("server_ids", ""))
		task.CronSpec = strings.TrimSpace(c.DefaultPostForm("cron_spec", ""))
		task.Command = strings.TrimSpace(c.DefaultPostForm("command", ""))
		task.Timeout, _ = strconv.Atoi(c.DefaultPostForm("timeout", "0"))
		task.IsNotify, _ = strconv.Atoi(c.DefaultPostForm("is_notify", "0"))
		task.ServerType, _ = strconv.Atoi(c.DefaultPostForm("server_type", "0"))

		task.NotifyType, _ = strconv.Atoi(c.DefaultPostForm("notify_type", "0"))
		task.NotifyTplId, _ = strconv.Atoi(c.DefaultPostForm("notify_tpl_id", "0"))
		task.NotifyUserIds = strings.TrimSpace(c.DefaultPostForm("notify_user_ids", ""))

		if task.IsNotify == 1 && task.NotifyTplId <= 0 {
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "请选择通知模板"))
			// self.ajaxMsg("请选择通知模板", MSG_ERR)
			return
		}

		msg, isBan := checkCommand(task.Command)
		if !isBan {
			// self.ajaxMsg("含有禁止命令："+msg, MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "含有禁止命令："+msg))
			return
		}

		task.CreateTime = time.Now().Unix()
		task.UpdateTime = time.Now().Unix()
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

		if _, err := model.TaskAdd(task); err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}

		// self.ajaxMsg("", MSG_OK)
		c.JSON(http.StatusOK, common.Success(c))
		return
	}

	task, _ := model.TaskGetById(task_id)
	//修改
	task.Id = task_id
	task.UpdateTime = time.Now().Unix()
	task.TaskName = strings.TrimSpace(c.DefaultPostForm("task_name", ""))
	task.Description = strings.TrimSpace(c.DefaultPostForm("description", ""))
	task.GroupId, _ = strconv.Atoi(c.DefaultPostForm("group_id", "0"))
	task.Concurrent, _ = strconv.Atoi(c.DefaultPostForm("concurrent", "0"))
	task.ServerIds = strings.TrimSpace(c.DefaultPostForm("server_ids", ""))
	task.CronSpec = strings.TrimSpace(c.DefaultPostForm("cron_spec", ""))
	task.Command = strings.TrimSpace(c.DefaultPostForm("command", ""))
	task.Timeout, _ = strconv.Atoi(c.DefaultPostForm("timeout", "0"))
	task.ServerType, _ = strconv.Atoi(c.DefaultPostForm("server_type", "0"))
	task.IsNotify, _ = strconv.Atoi(c.DefaultPostForm("is_notify", "0"))
	task.NotifyType, _ = strconv.Atoi(c.DefaultPostForm("notify_type", "0"))
	task.NotifyTplId, _ = strconv.Atoi(c.DefaultPostForm("notify_tpl_id", "0"))
	task.NotifyUserIds = strings.TrimSpace(c.DefaultPostForm("notify_user_ids", ""))
	task.UpdateId = uid
	task.Status = 2 //审核中,超级管理员不需要

	if task.IsNotify == 1 && task.NotifyTplId <= 0 {
		// self.ajaxMsg("请选择通知模板", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "请选择通知模板"))
		return
	}

	if self.userId == 1 {
		task.Status = 0
	}
	msg, isBan := checkCommand(task.Command)
	if !isBan {
		// self.ajaxMsg("含有禁止命令："+msg, MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "含有禁止命令："+msg))
		return
	}

	if _, err := cron.Parse(task.CronSpec); err != nil {
		// self.ajaxMsg("cron表达式无效", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "cron表达式无效"))
		return
	}

	if err := task.Update(); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Success(c))
}

//检查是否含有禁用命令
func checkCommand(command string) (string, bool) {

	filters := make([]interface{}, 0)
	filters = append(filters, "status", 0)
	ban, _ := model.BanGetList(1, 1000, filters...)
	for _, v := range ban {
		if strings.Contains(command, v.Code) {
			return v.Code, false
		}
	}
	return "", true
}

func (self *TaskController) Audit(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.DefaultPostForm("id", ""))

	// taskId, _ := self.GetInt("id")
	if taskID == 0 {
		// self.ajaxMsg("任务不存在", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "任务不存在"))
		return
	}

	res := changeStatus(taskID, 0, c.GetInt("uid"))
	if !res {
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

	res := changeStatus(taskID, 3, c.GetInt("uid"))
	if !res {
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

	task, err := model.TaskGetById(taskID)
	if err != nil {
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

	jobArr, err := jobs.NewJobFromTask(task)

	if err != nil {
		// self.ajaxMsg("创建任务失败", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "创建任务失败"))
		return
	}

	for _, job := range jobArr {
		if jobs.AddJob(task.CronSpec, job) {
			task.Status = 1
			task.Update()
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

	task, err := model.TaskGetById(taskID)
	if err != nil {
		// self.ajaxMsg("查不到该任务", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "查不到该任务"))
		return
	}

	//移出任务
	TaskServerIdsArr := strings.Split(task.ServerIds, ",")
	for _, server_id := range TaskServerIdsArr {
		server_id_int, _ := strconv.Atoi(server_id)
		jobKey := libs.JobKey(task.Id, server_id_int)
		jobs.RemoveJob(jobKey)
	}

	task.Status = 0
	task.Update()

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

// AjaxRun ss
func (self *TaskController) AjaxRun(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	task, err := model.TaskGetById(id)
	if err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	jobArr, err := jobs.NewJobFromTask(task)
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

		if task, err := model.TaskGetById(id); err == nil {
			jobArr, err := jobs.NewJobFromTask(task)
			if err == nil {
				for _, job := range jobArr {
					jobs.AddJob(task.CronSpec, job)
				}

				task.Status = 1
				task.Update()
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

		task, err := model.TaskGetById(id)
		fmt.Println(task)

		// 移出任务
		TaskServerIdsArr := strings.Split(task.ServerIds, ",")
		for _, server_id := range TaskServerIdsArr {
			server_id_int, _ := strconv.Atoi(server_id)
			jobKey := libs.JobKey(task.Id, server_id_int)
			jobs.RemoveJob(jobKey)
		}

		if err == nil {
			task.Status = 0
			task.Update()
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

		task, _ := model.TaskGetById(id)

		//移出任务
		TaskServerIdsArr := strings.Split(task.ServerIds, ",")

		for _, server_id := range TaskServerIdsArr {
			server_id_int, _ := strconv.Atoi(server_id)
			jobKey := libs.JobKey(task.Id, server_id_int)
			jobs.RemoveJob(jobKey)
		}
		model.TaskDel(id)
		model.TaskLogDelByTaskId(id)
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
		changeStatus(id, 0, self.userId)
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
		changeStatus(id, 3, self.userId)
	}

	c.JSON(http.StatusOK, common.Success(c))
}

func changeStatus(taskId, status, userId int) bool {
	if taskId == 0 {
		return false
	}

	task, _ := model.TaskGetById(taskId)
	//修改
	task.Id = taskId
	task.UpdateTime = time.Now().Unix()
	task.UpdateId = userId
	task.Status = status //0,1,2,3,9

	if err := task.Update(); err != nil {
		return false
	}
	return true
}

// AjaxDel ddd
func (self *TaskController) AjaxDel(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	task, _ := model.TaskGetById(id)

	uid := c.GetInt("uid")
	task.UpdateTime = time.Now().Unix()
	task.UpdateId = uid
	task.Status = -1
	task.Id = id

	//TODO 查询服务器是否用于定时任务
	if err := task.Update(); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *TaskController) AjaxNotifyType(c *gin.Context) {
	notifyType, _ := strconv.Atoi(c.DefaultPostForm("notify_type", "0"))
	result, count, _ := model.NotifyTplGetByTplTypeList(notifyType)

	list := make([]map[string]interface{}, len(result))

	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["tpl_name"] = v.TplName
		row["tpl_type"] = v.TplType
		list[k] = row
	}

	// self.ajaxList("成功", MSG_OK, count, list)
	ext := map[string]int{"total": int(count)}
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
	taskGroups, _ := service.TaskGroups(uid, "0")
	taskGroup := taskGroupLists(taskGroups, uid)
	self.pageSize = pagesize

	// 查询条件
	filters := make([]interface{}, 0)

	if status == 2 {
		//审核中，审核失败
		ids := []int{2, 3}
		filters = append(filters, "status__in", ids)
	} else {
		ids := []int{0, 1}
		filters = append(filters, "status__in", ids)
	}

	// 搜索全部
	if groupId == 0 {
		if uid != 1 {
			groups := strings.Split(taskGroups, ",")
			groupsIds := make([]int, 0)
			for _, v := range groups {
				id, _ := strconv.Atoi(v)
				groupsIds = append(groupsIds, id)
			}
			filters = append(filters, "group_id__in", groupsIds)
		}
	} else if groupId > 0 {
		filters = append(filters, "group_id", groupId)
	}

	if taskName != "" {
		filters = append(filters, "task_name__icontains", taskName)
	}

	result, count := model.TaskGetList(page, pagesize, filters...)
	list := make([]map[string]interface{}, len(result))

	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id

		groupName := "默认分组"

		if name, ok := taskGroup[v.GroupId]; ok {
			groupName = name
		}

		row["group_name"] = groupName
		row["task_name"] = StatusText[v.Status] + " " + groupName + "-" + "&nbsp;" + v.TaskName
		row["description"] = v.Description

		//row["status_text"] = StatusText[v.Status]
		row["status"] = v.Status
		row["pre_time"] = beego.Date(time.Unix(v.PrevTime, 0), "Y-m-d H:i:s")
		row["execute_times"] = v.ExecuteTimes
		row["cron_spec"] = v.CronSpec

		TaskServerIdsArr := strings.Split(v.ServerIds, ",")
		serverId := 0
		if len(TaskServerIdsArr) > 0 {
			serverId, _ = strconv.Atoi(TaskServerIdsArr[0])
		}
		jobskey := libs.JobKey(v.Id, serverId)
		e := jobs.GetEntryById(jobskey)

		if e != nil {
			row["next_time"] = beego.Date(e.Next, "Y-m-d H:i:s")
			row["prev_time"] = "-"
			if e.Prev.Unix() > 0 {
				row["prev_time"] = beego.Date(e.Prev, "Y-m-d H:i:s")
			} else if v.PrevTime > 0 {
				row["prev_time"] = beego.Date(time.Unix(v.PrevTime, 0), "Y-m-d H:i:s")
			}
			row["running"] = 1
		} else {
			row["next_time"] = "-"
			if v.PrevTime > 0 {
				row["prev_time"] = beego.Date(time.Unix(v.PrevTime, 0), "Y-m-d H:i:s")
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
		task.CreateId, _ = strconv.Atoi(c.DefaultPostForm("create_id", "0"))
		task.GroupId, _ = strconv.Atoi(c.DefaultPostForm("group_id", "0"))
		task.TaskName = strings.TrimSpace(c.DefaultPostForm("task_name", ""))
		task.Description = strings.TrimSpace(c.DefaultPostForm("description", ""))
		task.Concurrent, _ = strconv.Atoi(c.DefaultPostForm("concurrent", "0"))
		task.ServerIds = strings.TrimSpace(c.DefaultPostForm("server_ids", ""))
		task.CronSpec = strings.TrimSpace(c.DefaultPostForm("cron_spec", ""))
		task.Command = strings.TrimSpace(c.DefaultPostForm("command", ""))
		task.Timeout, _ = strconv.Atoi(c.DefaultPostForm("timeout", "0"))
		task.IsNotify, _ = strconv.Atoi(c.DefaultPostForm("is_notify", "0"))
		task.NotifyType, _ = strconv.Atoi(c.DefaultPostForm("notify_type", "0"))
		task.NotifyTplId, _ = strconv.Atoi(c.DefaultPostForm("notify_tpl_id", "0"))
		task.NotifyUserIds = strings.TrimSpace(c.DefaultPostForm("notify_user_ids", ""))

		if task.IsNotify == 1 && task.NotifyTplId <= 0 {
			// self.ajaxMsg("请选择通知模板", MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "请选择通知模板"))
			return
		}

		msg, isBan := checkCommand(task.Command)
		if !isBan {
			// self.ajaxMsg("含有禁止命令："+msg, MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "含有禁止命令："+msg))
			return
		}

		task.CreateTime = time.Now().Unix()
		task.UpdateTime = time.Now().Unix()
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

		if id, err = model.TaskAdd(task); err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}

		taskID = int(id)
		// self.ajaxMsg(task_id, MSG_OK)
		c.JSON(http.StatusOK, common.Success(c, map[string]int{"task_id": taskID}))
		return
	}

	task, _ := model.TaskGetById(taskID)

	if task.Status == 1 {
		// self.ajaxMsg("运行状态无法编辑任务，请先暂停任务", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "运行状态无法编辑任务，请先暂停任务"))
		return
	}

	//修改
	task.Id = taskID
	task.UpdateTime = time.Now().Unix()
	task.TaskName = strings.TrimSpace(c.PostForm("task_name"))
	task.Description = strings.TrimSpace(c.PostForm("description"))
	task.GroupId, _ = strconv.Atoi(c.DefaultPostForm("group_id", "0"))
	task.Concurrent, _ = strconv.Atoi(c.DefaultPostForm("concurrent", "0"))
	task.ServerIds = strings.TrimSpace(c.PostForm("server_ids"))
	task.CronSpec = strings.TrimSpace(c.PostForm("cron_spec"))
	task.Command = strings.TrimSpace(c.PostForm("command"))
	task.Timeout, _ = strconv.Atoi(c.DefaultPostForm("timeout", "0"))
	task.IsNotify, _ = strconv.Atoi(c.DefaultPostForm("is_notify", "0"))
	task.NotifyType, _ = strconv.Atoi(c.DefaultPostForm("notify_type", "0"))
	task.NotifyTplId, _ = strconv.Atoi(c.DefaultPostForm("notify_tpl_id", "0"))
	task.NotifyUserIds = strings.TrimSpace(c.PostForm("notify_user_ids"))
	task.UpdateId, _ = strconv.Atoi(c.DefaultPostForm("update_id", "0"))
	task.Status = 0 //接口不需要

	if task.IsNotify == 1 && task.NotifyTplId <= 0 {
		// self.ajaxMsg("请选择通知模板", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "请选择通知模板"))
		return
	}

	msg, isBan := checkCommand(task.Command)
	if !isBan {
		// self.ajaxMsg("含有禁止命令："+msg, MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "含有禁止命令："+msg))
		return
	}

	if _, err := cron.Parse(task.CronSpec); err != nil {
		// self.ajaxMsg("cron表达式无效", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "cron表达式无效"))
		return
	}

	if err := task.Update(); err != nil {
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

	task, err := model.TaskGetById(taskID)
	if err != nil {
		// self.ajaxMsg("查不到该任务", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "查不到该任务"))
		return
	}

	if task.Status != 0 {
		// self.ajaxMsg("任务状态有误", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "任务状态有误"))
		return
	}

	jobArr, err := jobs.NewJobFromTask(task)
	if err != nil {
		// self.ajaxMsg("创建任务失败", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "创建任务失败"))
		return
	}

	for _, job := range jobArr {
		if jobs.AddJob(task.CronSpec, job) {
			task.Status = 1
			task.Update()
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

	task, err := model.TaskGetById(taskID)
	if err != nil {
		// self.ajaxMsg("查不到该任务", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "查不到该任务"))
		return
	}

	//移出任务
	TaskServerIdsArr := strings.Split(task.ServerIds, ",")

	for _, server_id := range TaskServerIdsArr {
		server_id_int, _ := strconv.Atoi(server_id)
		jobKey := libs.JobKey(task.Id, server_id_int)
		jobs.RemoveJob(jobKey)
	}

	task.Status = 0
	task.Update()

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}
