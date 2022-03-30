/**********************************************
** @Des: This file ...
** @Author: haodaquan
** @Date:   2017-09-14 14:23:32
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-17 11:31:13
***********************************************/
package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/app/service"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/utils"
)

type RoleController struct {
	BaseController
}

func (self *RoleController) List(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "角色管理"

	// self.display()
	c.HTML(http.StatusOK, "role/list.html", data)
}

func (self *RoleController) Add(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	uid := c.GetInt("uid")
	tg, sg := service.TaskGroups(uid, c.GetString("role_id"))

	data["zTree"] = true //引入ztreecss
	data["taskGroup"] = taskGroupLists(tg, uid)
	data["serverGroup"] = serverLists(sg, uid)
	data["pageTitle"] = "新增角色"

	// self.display()
	c.HTML(http.StatusOK, "role/add.html", data)
}
func (self *RoleController) Edit(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	uid := c.GetInt("uid")
	tg, sg := service.TaskGroups(uid, c.GetString("role_id"))

	data["zTree"] = true //引入ztreecss
	data["pageTitle"] = "编辑角色"

	data["taskGroup"] = taskGroupLists(tg, uid)
	data["serverGroup"] = serverLists(sg, uid)

	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	role, _ := model.RoleGetById(id)
	row := make(map[string]interface{})
	row["id"] = role.Id
	row["role_name"] = role.RoleName
	row["detail"] = role.Detail
	row["task_group_ids"] = role.TaskGroupIds
	row["server_group_ids"] = role.ServerGroupIds
	data["role"] = row

	//获取选择的树节点
	roleAuth, _ := model.RoleAuthGetById(id)
	authId := make([]int, 0)
	for _, v := range roleAuth {
		authId = append(authId, v.AuthId)
	}

	taskGroupIdsArr := strings.Split(role.TaskGroupIds, ",")
	taskGroupIds := make([]int, 0)
	for _, v := range taskGroupIdsArr {
		id, _ := strconv.Atoi(v)
		taskGroupIds = append(taskGroupIds, id)
	}

	serverGroupIdsArr := strings.Split(role.ServerGroupIds, ",")
	serverGroupIds := make([]int, 0)
	for _, v := range serverGroupIdsArr {
		id, _ := strconv.Atoi(v)
		serverGroupIds = append(serverGroupIds, id)
	}

	data["server_group_ids"] = serverGroupIds
	data["task_group_ids"] = taskGroupIds

	data["auth"] = authId
	// self.display()

	c.HTML(http.StatusOK, "role/edit.html", data)
}

func (self *RoleController) AjaxSave(c *gin.Context) {

	uid := c.GetInt("uid")
	role := new(model.Role)
	role.RoleName = strings.TrimSpace(c.DefaultPostForm("role_name", ""))
	role.Detail = strings.TrimSpace(c.DefaultPostForm("detail", ""))
	role.ServerGroupIds = strings.TrimSpace(c.DefaultPostForm("server_group_ids", ""))
	role.TaskGroupIds = strings.TrimSpace(c.DefaultPostForm("task_group_ids", ""))
	role.CreateTime = time.Now().Unix()
	role.UpdateTime = time.Now().Unix()
	role.Status = 1

	auths := strings.TrimSpace(c.DefaultPostForm("nodes_data", ""))
	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	if id == 0 {
		//新增
		role.CreateTime = time.Now().Unix()
		role.UpdateTime = time.Now().Unix()
		role.CreateId = uid
		role.UpdateId = uid

		id, err := model.RoleAdd(role)
		if err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}

		ras := make([]model.RoleAuth, 0)
		authsSlice := strings.Split(auths, ",")
		for _, v := range authsSlice {
			//ra := new(model.RoleAuth)
			ra := model.RoleAuth{}
			aid, _ := strconv.Atoi(v)
			ra.AuthId = aid
			ra.RoleId = id
			ras = append(ras, ra)
		}

		if len(ras) > 0 {
			model.RoleAuthBatchAdd(&ras)
		}

		// self.ajaxMsg("", MSG_OK)
		c.JSON(http.StatusOK, common.Success(c))
		return
	}

	//修改
	role.Id = id
	role.UpdateTime = time.Now().Unix()
	role.UpdateId = self.userId
	if err := role.Update(); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}
	// 删除该角色权限
	model.RoleAuthDelete(id)

	ras := make([]model.RoleAuth, 0)
	authsSlice := strings.Split(auths, ",")
	for _, v := range authsSlice {
		//ra := new(model.RoleAuth)
		ra := model.RoleAuth{}
		aid, _ := strconv.Atoi(v)
		ra.AuthId = aid
		ra.RoleId = int64(id)
		ras = append(ras, ra)
	}
	if len(ras) > 0 {
		model.RoleAuthBatchAdd(&ras)
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *RoleController) AjaxDel(c *gin.Context) {

	id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	role, err := model.RoleGetById(id)
	if err != nil || role == nil {
		msg := "角色ID错误"
		if err != nil {
			msg = err.Error()
		}

		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, msg))
		return
	}

	role.Status = 0
	role.Id = id
	role.UpdateTime = time.Now().Unix()

	if err := role.Update(); err != nil {
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	// 删除该角色权限
	//model.RoleAuthDelete(role_id)
	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *RoleController) Table(c *gin.Context) {
	//列表
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pagesize", "20"))

	roleName := strings.TrimSpace(c.DefaultQuery("roleName", ""))

	//查询条件
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)
	if roleName != "" {
		filters = append(filters, "role_name__icontains", roleName)
	}
	result, count := model.RoleGetList(page, pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["role_name"] = v.RoleName
		row["detail"] = v.Detail
		row["create_time"] = time.Unix(v.CreateTime, 0).Format("2006-01-02 15:04:05")
		row["update_time"] = time.Unix(v.UpdateTime, 0).Format("2006-01-02 15:04:05")
		list[k] = row
	}

	ext := map[string]int{"total": int(count)}
	c.JSON(http.StatusOK, common.Success(c, list, ext))
}
