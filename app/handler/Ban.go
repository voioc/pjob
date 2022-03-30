/************************************************************
** @Description: controllers
** @Author: haodaquan
** @Date:   2018-06-10 19:50
** @Last Modified by:   haodaquan
** @Last Modified time: 2018-06-10 19:50
*************************************************************/
package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/utils"
)

type BanController struct {
	BaseController
}

func (self *BanController) List(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "禁用命令管理"
	// self.display()
	c.HTML(http.StatusOK, "ban/list.html", data)
}

func (self *BanController) Add(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	data["pageTitle"] = "新增禁用命令"

	// 角色
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)
	result, _ := model.RoleGetList(1, 1000, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["role_name"] = v.RoleName
		list[k] = row
	}

	data["role"] = list

	// self.display()
	c.HTML(http.StatusOK, "ban/add.html", data)
}

func (self *BanController) Edit(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	data["pageTitle"] = "编辑禁用命令"

	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	ban, _ := model.BanGetById(id)
	row := make(map[string]interface{})
	row["id"] = ban.Id
	row["code"] = ban.Code
	data["ban"] = row
	// self.display()
	c.HTML(http.StatusOK, "ban/edit.html", data)
}

func (self *BanController) AjaxSave(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	if id == 0 {
		ban := new(model.Ban)
		ban.Code = strings.TrimSpace(c.DefaultPostForm("code", ""))
		ban.CreateTime = time.Now().Unix()

		if _, err := model.BanAdd(ban); err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}

		c.JSON(http.StatusOK, common.Success(c))
		return
	}

	ban, _ := model.BanGetById(id)
	//修改
	// ban.Id = id
	ban.UpdateTime = time.Now().Unix()
	ban.Code = strings.TrimSpace(c.DefaultPostForm("code", ""))

	if err := ban.Update(); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Success(c))
}

func (self *BanController) AjaxDel(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	ban, _ := model.BanGetById(id)
	ban.UpdateTime = time.Now().Unix()
	ban.Status = 1

	if err := ban.Update(); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	// self.ajaxMsg("操作成功", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *BanController) Table(c *gin.Context) {
	//列表
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pagesize", "20"))

	code := strings.TrimSpace(c.DefaultQuery("code", ""))

	//查询条件
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 0)
	if code != "" {
		filters = append(filters, "code__icontains", code)
	}
	result, count := model.BanGetList(page, pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["code"] = v.Code
		row["create_time"] = time.Unix(v.CreateTime, 0).Format("2006-01-02 15:04:05")
		list[k] = row
	}

	// self.ajaxList("成功", MSG_OK, count, list)
	ext := map[string]int{"total": int(count)}
	c.JSON(http.StatusOK, common.Success(c, list, ext))
}
