/**********************************************
** @Des: 用户
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
	"github.com/voioc/cjob/utils"
)

type UserController struct {
	BaseController
}

func (self *UserController) Edit(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")
	data["pageTitle"] = "资料修改"

	uid := c.GetInt("uid")
	// Admin, _ := model.AdminGetById(uid)
	Admin := model.Admin{}
	if err := model.DataByID(&Admin, uid); err != nil {
		fmt.Println(err.Error())
	}

	row := make(map[string]interface{})
	row["id"] = Admin.ID
	row["login_name"] = Admin.LoginName
	row["real_name"] = Admin.RealName
	row["phone"] = Admin.Phone
	row["email"] = Admin.Email
	row["dingtalk"] = Admin.Dingtalk
	row["wechat"] = Admin.Wechat
	data["admin"] = row
	// self.display()

	c.HTML(http.StatusOK, "user/edit.html", data)
}

func (self *UserController) AjaxSave(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	// admin, err := model.AdminGetById(id)
	admin := model.Admin{}
	if err := model.DataByID(&admin, c.GetInt("uid")); err != nil || admin.ID == 0 {
		msg := "用户ID错误"
		if err != nil {
			msg = err.Error()
		}

		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, msg))
		return
	}

	//修改
	admin.ID = id
	admin.UpdatedAt = time.Now().Unix()
	admin.UpdatedID = c.GetInt("uid")
	admin.LoginName = strings.TrimSpace(c.PostForm("login_name"))
	admin.RealName = strings.TrimSpace(c.PostForm("real_name"))
	admin.Phone = strings.TrimSpace(c.PostForm("phone"))
	admin.Email = strings.TrimSpace(c.PostForm("email"))
	admin.Dingtalk = strings.TrimSpace(c.PostForm("dingtalk"))
	admin.Wechat = strings.TrimSpace(c.PostForm("wechat"))

	resetPwd := strings.TrimSpace(c.PostForm("reset_pwd"))
	if resetPwd == "1" {
		pwdOld := strings.TrimSpace(c.PostForm("password_old"))
		pwdOldMd5 := utils.Md5([]byte(pwdOld + admin.Salt))
		if admin.Password != pwdOldMd5 {
			// self.ajaxMsg("旧密码错误", MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "旧密码错误"))
			return
		}

		pwdNew1 := strings.TrimSpace(c.PostForm("password_new1"))
		pwdNew2 := strings.TrimSpace(c.PostForm("password_new2"))

		if len(pwdNew1) < 6 {
			// self.ajaxMsg("密码长度需要六位以上", MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "密码长度需要六位以上"))
			return
		}
		if pwdNew1 != pwdNew2 {
			// self.ajaxMsg("两次密码不一致", MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "密码长度需要六位以上"))
			return
		}

		pwd, salt := utils.Password(4, pwdNew1)
		admin.Password = pwd
		admin.Salt = salt
	}
	admin.UpdatedAt = time.Now().Unix()
	admin.UpdatedID = c.GetInt("uid")
	admin.Status = 1

	if err := model.Update(admin.ID, &admin); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	// self.ajaxMsg("修改成功", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}
