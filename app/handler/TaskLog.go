/************************************************************
** @Description: controllers
** @Author: george hao
** @Date:   2018-07-05 16:43
** @Last Modified by:  george hao
** @Last Modified time: 2018-07-05 16:43
*************************************************************/
package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/app/service"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/utils"

	"strings"
	"time"
)

type TaskLogController struct {
	BaseController
}

func (self *TaskLogController) List(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.DefaultQuery("task_id", "0"))
	// if err != nil {
	// 	taskId = 1
	// }

	// task, err := service.TaskS(c).TaskByID(taskID) // model.TaskGetById(taskId)
	task := model.Task{}
	if err := model.DataByID(&task, taskID); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	data["pageTitle"] = "日志管理 - " + task.TaskName + "(#" + strconv.Itoa(task.ID) + ")"
	data["task_id"] = task.ID

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
	taskID, err := strconv.Atoi(c.DefaultQuery("task_id", "0"))
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
		filters = append(filters, "status =", Status)
	}

	filters = append(filters, "task_id =", taskID)
	filters = append(filters, "order", "id desc")

	// result, count, _ := service.TaskLogS(c).LogList(page, pageSize, filters...) // model.TaskLogGetList(page, pageSize, filters...)
	result := make([]model.TaskLog, 0)
	if err := model.List(&result, page, pageSize, filters...); err != nil {
		fmt.Println(err.Error())
	}

	count, err := model.ListCount(&model.TaskLog{}, filters...)
	if err != nil {
		fmt.Println(err.Error())
	}

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["task_id"] = utils.JobKey(v.TaskID, v.ServerID)
		row["start_time"] = time.Unix(v.CreatedAt, 0).Format("2006-01-02 15:04:05")
		row["process_time"] = float64(v.ProcessTime) / 1000

		row["server_id"] = v.ServerID
		row["server_name"] = v.ServerName + "#" + strconv.Itoa(v.ServerID)
		if v.Status == 0 {
			row["output_size"] = utils.SizeFormat(float64(len(v.Output)))
		} else {
			row["output_size"] = utils.SizeFormat(float64(len(v.Error)))
		}
		index := v.Status + 2
		if index > 2 {
			index = 2
		}
		row["status"] = TextStatus[index]

		list[k] = row
	}

	// self.ajaxList("成功", MSG_OK, count, list)
	ext := map[string]int{"total": int(count)}
	c.JSON(http.StatusOK, common.Success(c, list, ext))
}

func (self *TaskLogController) Detail(c *gin.Context) {
	//日志内容
	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	// tasklogs, err := service.TaskLogS(c).LogByID([]int{id}) // model.TaskLogGetById(id)
	tasklog := model.TaskLog{}
	if err := model.DataByID(&tasklog, id); err != nil {
		fmt.Println(err.Error())
	}

	// fmt.Println(tasklog)
	// if err != nil || len(tasklogs) < 1 {
	// 	c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "日志不存在"))
	// 	return
	// }

	// tasklog := tasklogs[id]

	LogTextStatus := []string{
		"<font color='orange'><i class='fa fa-question-circle'></i>超时</font>",
		"<font color='red'><i class='fa fa-times-circle'></i> 错误</font>",
		"<font color='green'><i class='fa fa-check-square'></i> 正常</font>",
	}

	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	row := make(map[string]interface{})
	row["id"] = tasklog.ID
	row["task_id"] = tasklog.TaskID
	row["start_time"] = time.Unix(tasklog.CreatedAt, 0).Format("2006-01-02 15:04:05")
	row["process_time"] = float64(tasklog.ProcessTime) / 1000
	if tasklog.Status == 0 {
		row["output_size"] = utils.SizeFormat(float64(len(tasklog.Output)))
	} else {
		row["output_size"] = utils.SizeFormat(float64(len(tasklog.Error)))
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

	// 任务详情
	// task, err := service.TaskS(c).TaskByID(tasklog.TaskID) // model.TaskGetById(tasklog.TaskID)
	task := model.Task{}
	if err := model.DataByID(&task, tasklog.TaskID); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	data["TextStatus"] = common.TextStatus[task.Status]
	data["CreateTime"] = time.Unix(task.CreatedAt, 0).Format("2006-01-02 15:04:05")
	data["UpdateTime"] = time.Unix(task.UpdatedAt, 0).Format("2006-01-02 15:04:05")
	data["task"] = task

	// 分组列表
	tg, _ := service.AuthS(c).TaskGroups(c.GetInt("uid"), c.GetString("role_id"))
	data["taskGroup"], _ = service.TaskGroupS(c).GroupIDName(tg) // taskGroupLists(tg, c.GetInt("uid"))

	serverName := ""
	if task.ServerIDs == "0" {
		serverName = "本地服务器"
	} else {
		serverIdSli := strings.Split(task.ServerIDs, ",")
		for _, v := range serverIdSli {
			if v == "0" {
				serverName = "本地服务器"
			}
		}

		// servers, err := service.ServerS(c).ServersListID(strings.Split(task.ServerIDs, ",")) // model.TaskServerGetByIds(task.ServerIDs)
		servers := make([]*model.TaskServer, 0)
		if err := model.DataByIDs(&servers, strings.Split(task.ServerIDs, ",")); err != nil || len(servers) == 0 {
			if err != nil {
				fmt.Println(err.Error())
			}
			serverName = "服务器异常!!  "
		}

		for _, server := range servers {
			if server.Status != 0 {
				serverName += server.ServerName + "【无效】 "
			} else {
				serverName += server.ServerName + " "
			}
		}
	}

	data["serverName"] = serverName

	//任务分组
	groupName := "默认分组"
	if task.GroupID > 0 {
		// group, err := service.TaskGroupS(c).GroupByID(task.GroupID) // model.GroupGetById(task.GroupID)
		group := model.TaskGroup{}
		if err := model.DataByID(&group, task.GroupID); err == nil {
			groupName = group.GroupName
		}
	}
	data["GroupName"] = groupName

	//创建人和修改人
	createName := "未知"
	updateName := "未知"
	if task.CreatedID > 0 {
		// admin, err := service.AdminS(c).AdminGetByID(task.CreatedID) // model.AdminGetById(task.CreatedID)
		admin := model.Admin{}
		if err := model.DataByID(&admin, task.GroupID); err == nil {
			createName = admin.RealName
		}
	}

	if task.UpdatedID > 0 {
		// admin, err := service.AdminS(c).AdminGetByID(task.UpdatedID) // model.AdminGetById(task.UpdatedID)
		admin := model.Admin{}
		if err := model.DataByID(&admin, task.GroupID); err == nil {
			createName = admin.RealName
		}
	}

	//是否出错通知
	data["adminInfo"] = []int{0}
	if task.NotifyUserIDs != "0" && task.NotifyUserIDs != "" {
		data["adminInfo"], _ = service.AdminS(c).AdminInfo(nil) // AllAdminInfo(task.NotifyUserIds)
	}

	data["CreateName"] = createName
	data["UpdateName"] = updateName
	data["pageTitle"] = "日志详细" + "(#" + strconv.Itoa(id) + ")"

	data["NotifyTplName"] = "未知"
	if task.IsNotify == 1 {
		notifyTpl, err := service.NotifyS(c).NotifyListIDs([]int{task.NotifyTplID}) // model.NotifyTplGetById(task.NotifyTplID)
		if err == nil && len(notifyTpl) > 0 {
			data["NotifyTplName"] = notifyTpl[0].TplName
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

	// for _, v := range idArr {
	// 	id, _ := strconv.Atoi(v)
	// 	if id < 1 {
	// 		continue
	// 	}

	// 	// model.TaskLogDelById(id)
	// }
	if err := model.Del(&model.TaskLog{}, idArr); err != nil {
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}
