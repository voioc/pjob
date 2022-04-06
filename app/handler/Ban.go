/************************************************************
** @Description: controllers
** @Author: haodaquan
** @Date:   2018-06-10 19:50
** @Last Modified by:   haodaquan
** @Last Modified time: 2018-06-10 19:50
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
	// result, _ := service.RoleS(c).RoleList(1, 1000, filters) // model.RoleGetList(1, 1000, filters...)
	result := make([]model.Role, 0)
	if err := model.List(&result, 1, 1000, filters); err != nil {
		fmt.Println(err.Error())
	}

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
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
	// ban, _ := service.BanS(c).BanByID(id) // model.BanGetById(id)
	ban := &model.Ban{}
	if err := model.DataByID(ban, id); err != nil {
		fmt.Println(err.Error())
	}

	row := make(map[string]interface{})
	row["id"] = ban.ID
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
		ban.CreatedAt = time.Now().Unix()

		if err := model.Add(ban); err != nil { // model.BanAdd(ban); err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}

		c.JSON(http.StatusOK, common.Success(c))
		return
	}

	ban := &model.Ban{}
	if err := model.DataByID(ban, id); err != nil {
		fmt.Println(err.Error())
	}

	//修改
	// ban.Id = id
	ban.UpdatedAt = time.Now().Unix()
	ban.Code = strings.TrimSpace(c.DefaultPostForm("code", ""))

	if err := model.Update(ban.ID, ban); err != nil { // ban.Update(); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Success(c))
}

func (self *BanController) AjaxDel(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))

	// ban, _ := service.BanS(c).BanByID(id) // model.BanGetById(id)
	ban := &model.Ban{}
	if err := model.DataByID(ban, id); err != nil {
		fmt.Println(err.Error())
	}

	ban.UpdatedAt = time.Now().Unix()
	ban.Status = 2

	if err := model.Update(ban.ID, ban, true); err != nil { // ban.Update(); err != nil {
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
	filters = append(filters, "status =", 0)
	if code != "" {
		filters = append(filters, "code LIKE '%"+code+"%'", code)
	}

	// result, count, _ := service.BanS(c).BanList(page, pageSize, filters...) // model.BanGetList(page, pageSize, filters...)
	result := make([]model.Ban, 0)
	if err := model.List(&result, page, pageSize, filters...); err != nil {
		fmt.Println(err.Error())
	}

	count, _ := model.ListCount(&model.Ban{}, filters...)

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["code"] = v.Code
		row["create_time"] = time.Unix(v.CreatedAt, 0).Format("2006-01-02 15:04:05")
		list[k] = row
	}

	// self.ajaxList("成功", MSG_OK, count, list)
	ext := map[string]int{"total": int(count)}
	c.JSON(http.StatusOK, common.Success(c, list, ext))
}
