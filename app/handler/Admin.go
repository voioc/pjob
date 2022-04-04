/**********************************************
** @Des: 管理员
** @Author: haodaquan
** @Date:   2017-09-16 14:17:37
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-17 11:14:07
***********************************************/
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
	"github.com/voioc/cjob/libs"
	"github.com/voioc/cjob/utils"
)

type AdminController struct {
	BaseController
}

func (self *AdminController) List(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "管理员管理"

	// self.display()
	c.HTML(http.StatusOK, "admin/list.html", data)
}

func (self *AdminController) Add(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "新增管理员"
	// 角色
	filters := make([]interface{}, 0)
	filters = append(filters, "status=", 1)
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
	c.HTML(http.StatusOK, "admin/add.html", data)
}

func (self *AdminController) Edit(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "编辑管理员"

	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	// Admin, _ := service.AdminS(c).AdminGetByID(id) // model.AdminGetById(id)
	admin := model.Admin{}
	if err := model.DataByID(&admin, id); err != nil {
		fmt.Println(id)
	}

	row := make(map[string]interface{})
	row["id"] = admin.ID
	row["login_name"] = admin.LoginName
	row["real_name"] = admin.RealName
	row["phone"] = admin.Phone
	row["email"] = admin.Email
	row["dingtalk"] = admin.Dingtalk
	row["wechat"] = admin.Wechat
	row["role_ids"] = admin.RoleIDs
	data["admin"] = row

	role_ids := strings.Split(admin.RoleIDs, ",")

	filters := make([]interface{}, 0)
	filters = append(filters, "status =", 1)
	// result, _ := service.RoleS(c).RoleList(1, 1000, filters...) // model.RoleGetList(1, 1000, filters...)
	result := make([]model.Role, 0)
	if err := model.List(&result, 1, 1000, filters...); err != nil {
		fmt.Println(err.Error())
	}

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["checked"] = 0
		for i := 0; i < len(role_ids); i++ {
			role_id, _ := strconv.Atoi(role_ids[i])
			if role_id == v.ID {
				row["checked"] = 1
			}
			// fmt.Println(role_ids[i])
		}
		row["id"] = v.ID
		row["role_name"] = v.RoleName
		list[k] = row
	}

	data["role"] = list
	// self.display()

	c.HTML(http.StatusOK, "admin/edit.html", data)
}

func (self *AdminController) AjaxSave(c *gin.Context) {
	uid := c.GetInt("uid")
	id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	if id == 0 {
		Admin := new(model.Admin)
		Admin.LoginName = strings.TrimSpace(c.DefaultPostForm("login_name", ""))
		Admin.RealName = strings.TrimSpace(c.DefaultPostForm("real_name", ""))
		Admin.Phone = strings.TrimSpace(c.DefaultPostForm("phone", ""))
		Admin.Email = strings.TrimSpace(c.DefaultPostForm("email", ""))
		Admin.Dingtalk = strings.TrimSpace(c.DefaultPostForm("dingtalk", ""))
		Admin.Wechat = strings.TrimSpace(c.DefaultPostForm("wechat", ""))
		Admin.RoleIDs = strings.TrimSpace(c.DefaultPostForm("roleids", ""))
		Admin.UpdatedAt = time.Now().Unix()
		Admin.UpdatedID = uid
		Admin.Status = 1

		// 检查登录名是否已经存在
		filters := make([]interface{}, 0)
		filters = append(filters, "login_name =", Admin.LoginName)
		filters = append(filters, "status =", 1)
		// result, _ := service.RoleS(c).RoleList(1, 1, filters...) // model.RoleGetList(1, 1000, filters...)
		result := make([]model.Role, 0)
		if err := model.List(&result, 1, 1000, filters...); err != nil {
			fmt.Println(err.Error())
		}

		// _, err := model.AdminGetByName(Admin.LoginName)
		if len(result) > 0 {
			// self.ajaxMsg("登录名已经存在", MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "登录名已经存在"))
			return
		}

		//新增
		pwd, salt := libs.Password(4, "")
		Admin.Password = pwd
		Admin.Salt = salt
		Admin.CreatedAt = time.Now().Unix()
		Admin.CreatedID = uid
		if err := model.Add(Admin); err != nil { // model.AdminAdd(Admin); err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}

		// self.ajaxMsg("", MSG_OK)
		c.JSON(http.StatusOK, common.Success(c))
		return
	}

	// Admin, _ := service.AdminS(c).AdminGetByID(id) // model.AdminGetById(id)
	Admin := model.Admin{}
	if err := model.DataByID(&Admin, id); err != nil {
		fmt.Println(err.Error())
	}
	//修改
	// Admin.Id = id
	Admin.UpdatedAt = time.Now().Unix()
	Admin.UpdatedID = uid
	Admin.LoginName = strings.TrimSpace(c.DefaultPostForm("login_name", ""))
	Admin.RealName = strings.TrimSpace(c.DefaultPostForm("real_name", ""))
	Admin.Phone = strings.TrimSpace(c.DefaultPostForm("phone", ""))
	Admin.Email = strings.TrimSpace(c.DefaultPostForm("email", ""))
	Admin.Dingtalk = strings.TrimSpace(c.DefaultPostForm("dingtalk", ""))
	Admin.Wechat = strings.TrimSpace(c.DefaultPostForm("wechat", ""))
	Admin.RoleIDs = strings.TrimSpace(c.DefaultPostForm("roleids", ""))
	Admin.UpdatedAt = time.Now().Unix()
	Admin.Status = 1

	resetPwd, _ := strconv.Atoi(c.DefaultPostForm("reset_pwd", "0"))
	if resetPwd == 1 {
		pwd, salt := libs.Password(4, "")
		Admin.Password = pwd
		Admin.Salt = salt
	}

	// 普通管理员不可修改超级管理员资料
	if uid != 1 && Admin.ID == 1 {
		// self.ajaxMsg("不可修改超级管理员资料", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "不可修改超级管理员资料"))
		return
	}
	if err := model.Update(Admin.ID, &Admin); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	// self.ajaxMsg(strconv.Itoa(resetPwd), MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *AdminController) AjaxDel(c *gin.Context) {

	id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	status := strings.TrimSpace(c.DefaultPostForm("status", "0"))
	if id == 1 {
		// self.ajaxMsg("超级管理员不允许操作", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "超级管理员不允许操作"))
		return
	}

	Admin_status := 0
	if status == "enable" {
		Admin_status = 1
	}

	// Admin, _ := service.AdminS(c).AdminGetByID(id) // model.AdminGetById(id)
	Admin := model.Admin{}
	if err := model.DataByID(&Admin, id); err != nil {
		fmt.Println(err.Error())
	}
	Admin.UpdatedAt = time.Now().Unix()
	Admin.Status = Admin_status
	Admin.ID = id

	if err := model.Update(Admin.ID, &Admin, true); err != nil { // Admin.Update(); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	// self.ajaxMsg("操作成功", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *AdminController) Table(c *gin.Context) {
	//列表
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pagesize", "20"))

	realName := strings.TrimSpace(c.DefaultQuery("realName", ""))

	StatusText := make(map[int]string)
	StatusText[0] = "<font color='red'>禁用</font>"
	StatusText[1] = "正常"

	//查询条件
	filters := make([]interface{}, 0)
	if realName != "" {
		filters = append(filters, "real_name LIKE '%"+realName+"%'", "")
	}

	// result, count, _ := service.AdminS(c).AdminList(page, pageSize, filters...) // model.AdminGetList(page, pageSize, filters...)
	result := make([]model.Admin, 0)
	if err := model.List(&result, page, pageSize, filters...); err != nil {
		fmt.Println(err.Error())
	}

	count, err := model.ListCount(&model.Admin{}, filters...)
	if err != nil {
		fmt.Println(err.Error())
	}

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["login_name"] = v.LoginName
		row["real_name"] = v.RealName
		row["phone"] = v.Phone
		row["email"] = v.Email
		row["dingtalk"] = v.Dingtalk
		row["wechat"] = v.Wechat
		row["role_ids"] = v.RoleIDs
		row["create_time"] = time.Unix(v.CreatedAt, 0).Format("2006-01-02 15:04:05")
		row["update_time"] = time.Unix(v.UpdatedAt, 0).Format("2006-01-02 15:04:05")
		row["status"] = v.Status
		row["status_text"] = StatusText[v.Status]
		list[k] = row
	}

	// self.ajaxList("成功", MSG_OK, count, list)
	ext := map[string]int{"total": int(count)}
	c.JSON(http.StatusOK, common.Success(c, list, ext))
}
