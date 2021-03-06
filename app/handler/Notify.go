/************************************************************
** @Description: controllers
** @Author: Bee
** @Date:   2019-02-15 20:21
** @Last Modified by:   Bee
** @Last Modified time: 2019-02-15 20:21
*************************************************************/
package handler

import (
	"encoding/json"
	"fmt"
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
	_, sg := service.AuthS(c).TaskGroups(uid, c.GetString("role_id"))
	data["serverGroup"], _ = service.ServerGroupS(c).GroupIDName(sg) // serverGroupLists(sg, uid)
	// self.display()

	c.HTML(http.StatusOK, "notify/add.html", data)
}

func (self *NotifyController) Edit(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	data["pageTitle"] = "编辑通知模板"

	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	// notifyTpl, _ := service.NotifyS(c).NotifyByID(id) // model.NotifyTplGetById(id)
	notifyTpl := model.NotifyTpl{}
	if err := model.DataByID(&notifyTpl, id); err != nil {
		fmt.Println(err.Error())
	}

	row := make(map[string]interface{})
	row["id"] = notifyTpl.ID
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
		notifyTpl := new(model.NotifyTpl)
		notifyTpl.TplName = strings.TrimSpace(c.DefaultPostForm("tpl_name", ""))
		notifyTpl.TplType, _ = strconv.Atoi(c.DefaultPostForm("tpl_type", "0"))
		notifyTpl.Title = strings.TrimSpace(c.DefaultPostForm("title", ""))
		notifyTpl.Content = strings.TrimSpace(c.DefaultPostForm("content", ""))
		notifyTpl.CreatedID = uid
		notifyTpl.CreatedAt = time.Now().Unix()
		notifyTpl.Type = model.NotifyTplTypeDefault
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

		if err := model.Add(&notifyTpl); err != nil { // model.NotifyTplAdd(notifyTpl); err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}

		c.JSON(http.StatusOK, common.Success(c))
		return
	}

	// notifyTpl, _ := service.NotifyS(c).NotifyByID(id) // model.NotifyTplGetById(id)
	notifyTpl := model.NotifyTpl{}
	if err := model.DataByID(&notifyTpl, id); err != nil {
		fmt.Println(err.Error())
	}

	//修改
	// notifyTpl.Id = id
	notifyTpl.UpdatedID = uid
	notifyTpl.UpdatedAt = time.Now().Unix()

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

	if notifyTpl.Type == model.NotifyTplTypeSystem {
		// self.ajaxMsg("系统模板禁止更新", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "系统模板禁止更新"))
		return
	}

	if err := model.Update(notifyTpl.ID, notifyTpl); err != nil { // notifyTpl.Update(); err != nil {
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
	// notifyTpl, _ := model.NotifyTplGetById(id)
	notifyTpl := model.NotifyTpl{}
	if err := model.DataByID(&notifyTpl, id); err != nil {
		fmt.Println(err.Error())
	}

	if notifyTpl.Type == model.NotifyTplTypeSystem {
		// self.ajaxMsg("系统模板禁止删除", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "系统模板禁止删除"))
		return
	}

	if err := model.Del(model.NotifyTpl{}, []int{id}); err != nil { // model.NotifyTplDelById(id); err != nil {
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
		"<font color='green'>正常</font>",
	}

	TplTypeText := []string{
		"未知",
		"邮件",
		"信息",
		"钉钉",
		"微信",
	}

	//查询条件
	filters := make([]interface{}, 0)

	if tplName != "" {
		filters = append(filters, "tpl_name LIKE '%"+tplName+"%'", "")
	}

	result, count, _ := service.NotifyS(c).NotifyList(page, pageSize, filters...) // model.NotifyTplGetList(page, pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["type"] = v.Type
		row["tpl_name"] = v.TplName
		row["tpl_type"] = v.TplType
		row["tpl_type_text"] = TplTypeText[v.TplType]
		row["status"] = v.Status
		row["status_text"] = StatusText[v.Status]
		row["create_time"] = time.Unix(v.CreatedAt, 0).Format("2006-01-02 15:04:05")
		row["update_time"] = time.Unix(v.UpdatedAt, 0).Format("2006-01-02 15:04:05")
		list[k] = row
	}

	// self.ajaxList("成功", MSG_OK, count, list)
	ext := map[string]int{"total": int(count)}
	c.JSON(http.StatusOK, common.Success(c, list, ext))
}
