/************************************************************
** @Description: controllers
** @Author: george hao
** @Date:   2018-07-05 16:43
** @Last Modified by:  george hao
** @Last Modified time: 2018-07-05 16:43
*************************************************************/
package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/libs"
	"github.com/voioc/cjob/models"
	"github.com/voioc/cjob/service"

	"strings"
	"time"
)

type TaskLogController struct {
	BaseController
}

func (self *TaskLogController) List(c *gin.Context) {
	taskId, _ := strconv.Atoi(c.DefaultQuery("task_id", "0"))
	// if err != nil {
	// 	taskId = 1
	// }

	task, err := models.TaskGetById(taskId)
	if err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	data := map[string]interface{}{}
	data["pageTitle"] = "日志管理 - " + task.TaskName + "(#" + strconv.Itoa(task.Id) + ")"
	data["task_id"] = task.Id

	c.HTML(http.StatusOK, "tasklog/list.html", data)
}

func (self *TaskLogController) Table(c *gin.Context) {
	//列表
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	// if err != nil {
	// 	page = 1
	// }

	pageSize, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	// if err != nil {
	// 	limit = 30
	// }

	// self.pageSize = limit
	//查询条件
	filters := make([]interface{}, 0)
	taskId, err := strconv.Atoi(c.DefaultQuery("task_id", "0"))
	// if err != nil {
	// 	taskId = 1
	// }

	TextStatus := []string{
		"<font color='orange'><i class='fa fa-question-circle'></i> 超时</font>",
		"<font color='red'><i class='fa fa-times-circle'></i> 错误</font>",
		"<font color='green'><i class='fa fa-check-square'></i> 正常</font>",
	}

	// Status, err := self.GetInt("status")
	Status, err := strconv.Atoi(c.DefaultQuery("status", "0"))
	if err == nil && Status != 9 {
		filters = append(filters, "status", Status)
	}

	filters = append(filters, "task_id", taskId)

	result, count := models.TaskLogGetList(page, pageSize, filters...)
	list := make([]map[string]interface{}, len(result))

	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["task_id"] = libs.JobKey(v.TaskId, v.ServerId)
		row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["process_time"] = float64(v.ProcessTime) / 1000

		row["server_id"] = v.ServerId
		row["server_name"] = v.ServerName + "#" + strconv.Itoa(v.ServerId)
		if v.Status == 0 {
			row["output_size"] = libs.SizeFormat(float64(len(v.Output)))
		} else {
			row["output_size"] = libs.SizeFormat(float64(len(v.Error)))
		}
		index := v.Status + 2
		if index > 2 {
			index = 2
		}
		row["status"] = TextStatus[index]

		list[k] = row
	}

	// self.ajaxList("成功", MSG_OK, count, list)
	ext := map[string]int{"count": int(count)}
	c.JSON(http.StatusOK, common.Success(c, list, ext))
}

func (self *TaskLogController) Detail(c *gin.Context) {

	//日志内容
	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	tasklog, err := models.TaskLogGetById(id)

	fmt.Println(tasklog)
	if err != nil {
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "日志不存在"))
		return
	}

	LogTextStatus := []string{
		"<font color='orange'><i class='fa fa-question-circle'></i>超时</font>",
		"<font color='red'><i class='fa fa-times-circle'></i> 错误</font>",
		"<font color='green'><i class='fa fa-check-square'></i> 正常</font>",
	}

	data := map[string]interface{}{}
	row := make(map[string]interface{})
	row["id"] = tasklog.Id
	row["task_id"] = tasklog.TaskId
	row["start_time"] = beego.Date(time.Unix(tasklog.CreateTime, 0), "Y-m-d H:i:s")
	row["process_time"] = float64(tasklog.ProcessTime) / 1000
	if tasklog.Status == 0 {
		row["output_size"] = libs.SizeFormat(float64(len(tasklog.Output)))
	} else {
		row["output_size"] = libs.SizeFormat(float64(len(tasklog.Error)))
	}

	row["server_name"] = tasklog.ServerName

	row["output"] = tasklog.Output
	row["error"] = tasklog.Error

	index := tasklog.Status + 2
	if index > 2 {
		index = 2
	}
	row["status"] = LogTextStatus[index]

	data["taskLog"] = row

	//任务详情
	task, err := models.TaskGetById(tasklog.TaskId)
	if err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	TextStatus := []string{
		"<font color='red'><i class='fa fa-minus-square'></i> 暂停</font>",
		"<font color='green'><i class='fa fa-check-square'></i> 运行中</font>",
		"<font color='orange'><i class='fa fa-check-square'></i> 待审核</font>",
		"<font color='blue'><i class='fa fa-times-circle'></i> 审核失败</font>",
	}

	data["TextStatus"] = TextStatus[task.Status]
	data["CreateTime"] = beego.Date(time.Unix(task.CreateTime, 0), "Y-m-d H:i:s")
	data["UpdateTime"] = beego.Date(time.Unix(task.UpdateTime, 0), "Y-m-d H:i:s")
	data["task"] = task
	// 分组列表
	tg, _ := service.TaskGroups(c.GetInt("uid"), c.GetString("role_id"))
	data["taskGroup"] = taskGroupLists(tg, c.GetInt("uid"))

	serverName := ""
	if task.ServerIds == "0" {
		serverName = "本地服务器"
	} else {
		serverIdSli := strings.Split(task.ServerIds, ",")
		for _, v := range serverIdSli {
			if v == "0" {
				serverName = "本地服务器  "
			}
		}
		servers, n := models.TaskServerGetByIds(task.ServerIds)
		if n > 0 {
			for _, server := range servers {
				if server.Status != 0 {
					serverName += server.ServerName + "【无效】 "
				} else {
					serverName += server.ServerName + " "
				}
			}
		} else {
			serverName = "服务器异常!!  "
		}
	}

	data["serverName"] = serverName

	//任务分组
	groupName := "默认分组"
	if task.GroupId > 0 {
		group, err := models.GroupGetById(task.GroupId)
		if err == nil {
			groupName = group.GroupName
		}
	}
	data["GroupName"] = groupName

	//创建人和修改人
	createName := "未知"
	updateName := "未知"
	if task.CreateId > 0 {
		admin, err := models.AdminGetById(task.CreateId)
		if err == nil {
			createName = admin.RealName
		}
	}

	if task.UpdateId > 0 {
		admin, err := models.AdminGetById(task.UpdateId)
		if err == nil {
			updateName = admin.RealName
		}
	}

	//是否出错通知
	data["adminInfo"] = []int{0}
	if task.NotifyUserIds != "0" && task.NotifyUserIds != "" {
		data["adminInfo"] = AllAdminInfo(task.NotifyUserIds)
	}

	data["CreateName"] = createName
	data["UpdateName"] = updateName
	data["pageTitle"] = "日志详细" + "(#" + strconv.Itoa(id) + ")"

	data["NotifyTplName"] = "未知"
	if task.IsNotify == 1 {
		notifyTpl, err := models.NotifyTplGetById(task.NotifyTplId)
		if err == nil {
			data["NotifyTplName"] = notifyTpl.TplName
		}
	}

	// self.display()
	c.HTML(http.StatusOK, "tasklog/detail.html", data)
}

// 批量操作日志
func (self *TaskLogController) AjaxDel(c *gin.Context) {
	ids := c.DefaultQuery("ids", "")
	idArr := strings.Split(ids, ",")

	if len(idArr) < 1 {
		// self.ajaxMsg("请选择要操作的项目", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "请选择要操作的项目"))
		return
	}

	for _, v := range idArr {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}
		models.TaskLogDelById(id)
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}
