/************************************************************
** @Description: controllers
** @Author: Bee
** @Date:   2019-02-15 20:21
** @Last Modified by:   Bee
** @Last Modified time: 2019-02-15 20:21
*************************************************************/
package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/models"
	"github.com/voioc/cjob/service"
	"github.com/voioc/cjob/utils"
)

type NotifyController struct {
	BaseController
}

func (self *NotifyController) List(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "通知模板"
	// self.display()
	c.HTML(http.StatusOK, "notify/list.html", data)
}

func (self *NotifyController) Add(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	data["pageTitle"] = "新增通知模板"

	uid := c.GetInt("uid")
	_, sg := service.TaskGroups(uid, c.GetString("role_id"))
	data["serverGroup"] = serverGroupLists(sg, uid)
	// self.display()

	c.HTML(http.StatusOK, "notify/add.html", data)
}

func (self *NotifyController) Edit(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	data["pageTitle"] = "编辑通知模板"

	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	notifyTpl, _ := models.NotifyTplGetById(id)
	row := make(map[string]interface{})
	row["id"] = notifyTpl.Id
	row["tpl_name"] = notifyTpl.TplName
	row["tpl_type"] = notifyTpl.TplType
	row["title"] = notifyTpl.Title
	row["content"] = notifyTpl.Content
	row["status"] = notifyTpl.Status
	data["notifyTpl"] = row

	// self.display()
	c.HTML(http.StatusOK, "notify/edit.html", data)
}

func (self *NotifyController) AjaxSave(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	uid := c.GetInt("uid")
	id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	if id == 0 {
		notifyTpl := new(models.NotifyTpl)
		notifyTpl.TplName = strings.TrimSpace(c.DefaultPostForm("tpl_name", ""))
		notifyTpl.TplType, _ = strconv.Atoi(c.DefaultPostForm("tpl_type", "0"))
		notifyTpl.Title = strings.TrimSpace(c.DefaultPostForm("title", ""))
		notifyTpl.Content = strings.TrimSpace(c.DefaultPostForm("content", ""))
		notifyTpl.CreateId = uid
		notifyTpl.CreateTime = time.Now().Unix()
		notifyTpl.Type = models.NotifyTplTypeDefault
		notifyTpl.Status, _ = strconv.Atoi(c.DefaultPostForm("status", "0"))

		if notifyTpl.TplType == 1 || notifyTpl.TplType == 2 || notifyTpl.TplType == 3 {
			m := make(map[string]string)
			err := json.Unmarshal([]byte(notifyTpl.Content), &m)
			if err != nil {
				// self.ajaxMsg("模板内容格式错误,"+err.Error(), MSG_ERR)
				c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "模板内容格式错误,"+err.Error()))
				return
			}
		}

		if _, err := models.NotifyTplAdd(notifyTpl); err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}

		c.JSON(http.StatusOK, common.Success(c))
		return
	}

	notifyTpl, _ := models.NotifyTplGetById(id)
	//修改
	// notifyTpl.Id = id
	notifyTpl.UpdateId = uid
	notifyTpl.UpdateTime = time.Now().Unix()

	notifyTpl.TplName = strings.TrimSpace(c.DefaultPostForm("tpl_name", ""))
	notifyTpl.TplType, _ = strconv.Atoi(c.DefaultPostForm("tpl_type", "0"))
	notifyTpl.Title = strings.TrimSpace(c.DefaultPostForm("title", ""))
	notifyTpl.Content = strings.TrimSpace(c.DefaultPostForm("content", ""))
	notifyTpl.Status, _ = strconv.Atoi(c.DefaultPostForm("status", "0"))

	if notifyTpl.TplType == 1 || notifyTpl.TplType == 2 || notifyTpl.TplType == 3 {
		m := make(map[string]string)
		err := json.Unmarshal([]byte(notifyTpl.Content), &m)
		if err != nil {
			// self.ajaxMsg("模板内容格式错误,"+err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "模板内容格式错误,"+err.Error()))
			return
		}
	}

	if notifyTpl.Type == models.NotifyTplTypeSystem {
		// self.ajaxMsg("系统模板禁止更新", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "系统模板禁止更新"))
		return
	}

	if err := notifyTpl.Update(); err != nil {
		// self.ajaxMsg("更新失败,"+err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "更新失败,"+err.Error()))
		return
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *NotifyController) AjaxDel(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	notifyTpl, _ := models.NotifyTplGetById(id)

	if notifyTpl.Type == models.NotifyTplTypeSystem {
		// self.ajaxMsg("系统模板禁止删除", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "系统模板禁止删除"))
		return
	}

	if err := models.NotifyTplDelById(id); err != nil {
		// self.ajaxMsg("删除失败,"+err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "删除失败: "+err.Error()))
		return
	}

	// self.ajaxMsg("操作成功", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *NotifyController) Table(c *gin.Context) {

	//列表
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pagesize", "20"))
	tplName := strings.TrimSpace(c.DefaultQuery("tplName", ""))

	StatusText := []string{
		"<font color='red'>禁用</font>",
		"正常",
	}

	TplTypeText := []string{
		"邮件",
		"信息",
		"钉钉",
		"微信",
	}

	//查询条件
	filters := make([]interface{}, 0)

	if tplName != "" {
		filters = append(filters, "tpl_name__icontains", tplName)
	}
	result, count := models.NotifyTplGetList(page, pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["type"] = v.Type
		row["tpl_name"] = v.TplName
		row["tpl_type"] = v.TplType
		row["tpl_type_text"] = TplTypeText[v.TplType]
		row["status"] = v.Status
		row["status_text"] = StatusText[v.Status]
		row["create_time"] = time.Unix(v.CreateTime, 0).Format("2006-01-02 15:04:05")
		row["update_time"] = time.Unix(v.UpdateTime, 0).Format("2006-01-02 15:04:05")
		list[k] = row
	}

	// self.ajaxList("成功", MSG_OK, count, list)
	ext := map[string]int{"total": int(count)}
	c.JSON(http.StatusOK, common.Success(c, list, ext))
}
