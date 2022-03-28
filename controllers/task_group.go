/************************************************************
** @Description: controllers
** @Author: haodaquan
** @Date:   2018-06-10 22:24
** @Last Modified by:   haodaquan
** @Last Modified time: 2018-06-10 22:24
*************************************************************/
package controllers

import (
	"net/http"
	"strings"
	"time"

	"fmt"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/models"
	"github.com/voioc/cjob/service"
	"github.com/voioc/cjob/utils"
)

type GroupController struct {
	BaseController
}

func (self *GroupController) List(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	data["pageTitle"] = "任务分组管理"

	// self.display()

	c.HTML(http.StatusOK, "group/list.html", data)
}

func (self *GroupController) Add(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	data["pageTitle"] = "新增分组"
	data["hideTop"] = true
	// self.display()

	c.HTML(http.StatusOK, "group/add.html", data)
}
func (self *GroupController) Edit(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "编辑分组"
	data["hideTop"] = true

	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	group, _ := models.GroupGetById(id)
	row := make(map[string]interface{})
	row["id"] = group.Id
	row["group_name"] = group.GroupName
	row["description"] = group.Description
	data["group"] = row

	// self.display()
	c.HTML(http.StatusOK, "group/edit.html", data)
}

func (self *GroupController) AjaxSave(c *gin.Context) {

	group := new(models.Group)
	group.GroupName = strings.TrimSpace(c.DefaultPostForm("group_name", ""))
	group.Description = strings.TrimSpace(c.DefaultPostForm("description", ""))
	group.Status = 1

	group_id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	fmt.Println(group_id)

	uid := c.GetInt("uid")
	if group_id == 0 {
		//新增
		group.CreateTime = time.Now().Unix()
		group.UpdateTime = time.Now().Unix()
		group.CreateId = uid
		group.UpdateId = uid
		if _, err := models.GroupAdd(group); err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}
		// self.ajaxMsg("", MSG_OK)
		c.JSON(http.StatusOK, common.Success(c))
		return
	}
	//修改
	group.Id = group_id
	group.UpdateTime = time.Now().Unix()
	group.UpdateId = self.userId
	if err := group.Update(); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *GroupController) AjaxDel(c *gin.Context) {

	group_id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	group, err := models.GroupGetById(group_id)
	if err != nil || group.Id == 0 {
		msg := "内部错误"
		if err != nil {
			msg = err.Error()
		}

		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, msg))
		return
	}

	group.Status = 0
	group.Id = group_id
	group.UpdateTime = time.Now().Unix()
	//TODO 如果分组下有任务 不处理
	//filters := make([]interface{}, 0)
	//filters = append(filters, "group_id", group_id)
	//filters = append(filters, "status", 0)
	//_, n := models.TaskServerGetList(1, 1, filters...)
	//if n > 0 {
	//	self.ajaxMsg("分组下有服务器资源，请先处理", MSG_ERR)
	//}

	if err := group.Update(); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}
	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *GroupController) Table(c *gin.Context) {
	uid := c.GetInt("uid")

	//列表
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pagesize", "20"))

	groupName := strings.TrimSpace(c.DefaultQuery("groupName", ""))
	// self.pageSize = limit

	//查询条件
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)

	if uid != 1 {
		tg, _ := service.TaskGroups(uid, c.GetString("role_id"))
		groups := strings.Split(tg, ",")

		groupsIds := make([]int, 0)
		for _, v := range groups {
			id, _ := strconv.Atoi(v)
			groupsIds = append(groupsIds, id)
		}
		filters = append(filters, "id__in", groupsIds)
	}
	if groupName != "" {
		filters = append(filters, "group_name__contains", groupName)
	}
	result, count := models.GroupGetList(page, pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["group_name"] = v.GroupName
		row["description"] = v.Description
		row["create_time"] = time.Unix(v.CreateTime, 0).Format("2006-01-02 15:04:05")
		row["update_time"] = time.Unix(v.UpdateTime, 0).Format("2006-01-02 15:04:05")
		list[k] = row
	}

	// self.ajaxList("成功", MSG_OK, count, list)
	ext := map[string]int{"count": int(count)}
	c.JSON(http.StatusOK, common.Success(c, list, ext))
}
