/************************************************************
** @Description: controllers
** @Author: haodaquan
** @Date:   2018-06-08 21:57
** @Last Modified by:   haodaquan
** @Last Modified time: 2018-06-08 21:57
*************************************************************/
package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/app/service"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/utils"
)

type ServerGroupController struct {
	BaseController
}

func (self *ServerGroupController) List(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "资源分组管理"
	// self.display()
	c.HTML(http.StatusOK, "servergroup/list.html", data)
}

func (self *ServerGroupController) Add(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "新增分组"
	data["hideTop"] = true

	// self.display()
	c.HTML(http.StatusOK, "servergroup/add.html", data)
}
func (self *ServerGroupController) Edit(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "编辑分组"
	data["hideTop"] = true

	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	// group, _ := service.TaskGroupS(c).GroupByID(id) // model.TaskGroupGetById(id)
	group := model.ServerGroup{}
	if err := model.DataByID(&group, id); err != nil {
		fmt.Println(err.Error())
	}

	row := make(map[string]interface{})
	row["id"] = group.ID
	row["group_name"] = group.GroupName
	row["description"] = group.Description
	data["group"] = row

	// self.display()
	c.HTML(http.StatusOK, "servergroup/edit.html", data)
}

func (self *ServerGroupController) AjaxSave(c *gin.Context) {
	servergroup := new(model.ServerGroup)
	servergroup.GroupName = strings.TrimSpace(c.DefaultPostForm("group_name", ""))
	servergroup.Description = strings.TrimSpace(c.DefaultPostForm("description", ""))
	servergroup.Status = 1

	servergroup_id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))

	uid := c.GetInt("uid")
	if servergroup_id == 0 {
		//新增
		servergroup.CreatedAt = time.Now().Unix()
		servergroup.UpdatedAt = time.Now().Unix()
		servergroup.CreatedID = uid
		servergroup.UpdatedID = uid
		if err := model.Add(servergroup); err != nil { //  model.ServerGroupAdd(servergroup); err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}

		c.JSON(http.StatusOK, common.Success(c))
		return
	}

	//修改
	servergroup.ID = servergroup_id
	servergroup.UpdatedAt = time.Now().Unix()
	servergroup.UpdatedID = uid
	if err := model.Update(servergroup.ID, servergroup); err != nil { // servergroup.Update(); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Success(c))
}

func (self *ServerGroupController) AjaxDel(c *gin.Context) {

	group_id, _ := strconv.Atoi(c.PostForm("id"))
	// group, _ := service.TaskGroupS(c).GroupByID(group_id) // model.TaskGroupGetById(group_id)
	group := model.ServerGroup{}
	if err := model.DataByID(&group, group_id); err != nil {
		fmt.Println(err.Error())
	}

	group.Status = 2
	group.ID = group_id
	group.UpdatedAt = time.Now().Unix()

	// TODO 如果分组下有服务器 需要处理
	filters := make([]interface{}, 0)
	filters = append(filters, "group_id =", group_id)
	filters = append(filters, "status =", 1)
	// _, n, _ := service.ServerS(c).ServerList(1, 1, filters...) // model.TaskServerGetList(1, 1, filters...)
	count, err := model.ListCount(&model.ServerGroup{}, 1, 1, filters)
	if err != nil {
		fmt.Println(err.Error())
	}

	if count > 0 {
		// self.ajaxMsg("分组下有服务器资源，请先处理", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "分组下有服务器资源，请先处理"))
		return
	}

	if err := model.Update(group.ID, group, true); err != nil { // group.Update(); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Success(c))
}

func (self *ServerGroupController) Table(c *gin.Context) {
	//列表
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pagesize", "20"))

	groupName := strings.TrimSpace(c.DefaultQuery("groupName", ""))
	//查询条件
	filters := make([]interface{}, 0)
	filters = append(filters, "status =", 1)

	uid := c.GetInt("uid")
	if uid != 1 {
		_, sg := service.AuthS(c).TaskGroups(uid, c.GetString("role_id"))
		groups := strings.Split(sg, ",")

		groupsIds := make([]int, 0)
		for _, v := range groups {
			id, _ := strconv.Atoi(v)
			groupsIds = append(groupsIds, id)
		}
		filters = append(filters, "id", groupsIds)
	}

	if groupName != "" {
		filters = append(filters, "group_name LIKE '%"+groupName+"'%", "")
	}

	// result, count, _ := service.ServerGroupS(c).List(page, pageSize, filters...) // model.ServerGroupGetList(page, pageSize, filters...)
	result := make([]model.ServerGroup, 0)
	if err := model.List(&result, page, pageSize, filters...); err != nil {
		fmt.Println(err.Error())
	}

	count, err := model.ListCount(&model.ServerGroup{}, filters...)
	if err != nil {
		fmt.Println(err.Error())
	}

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["group_name"] = v.GroupName
		row["description"] = v.Description
		row["create_time"] = time.Unix(v.CreatedAt, 0).Format("2006-01-02 15:04:05")
		row["update_time"] = time.Unix(v.UpdatedAt, 0).Format("2006-01-02 15:04:05")
		list[k] = row
	}

	// self.ajaxList("成功", MSG_OK, count, list)
	ext := map[string]int{"total": int(count)}
	c.JSON(http.StatusOK, common.Success(c, list, ext))
}
