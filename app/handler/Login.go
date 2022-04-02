/**********************************************
** @Des: login
** @Author: haodaquan
** @Date:   2017-09-07 16:30:10
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-17 11:55:21
***********************************************/
package handler

import (
	"net/http"

	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/utils"
)

type LoginController struct {
	BaseController
}

func (self *LoginController) Login(c *gin.Context) {
	//if self.userId > 0 {
	//	self.redirect(beego.URLFor("HomeController.Index"))
	//}
	// c.HTML("index", "login.html")

	// data := map[string]interface{}{"url"}
	fm := template.FuncMap{
		"url": utils.URI("login_in"),
	}
	c.HTML(http.StatusOK, "login/login.html", fm)
}

//登录 TODO:XSRF过滤
func (self *LoginController) LoginIn(c *gin.Context) {

	// self.ajaxMsg("登录成功", MSG_OK)
	// if self.userId > 0 {
	// 	self.ajaxMsg("登录成功", MSG_OK)
	// }

	// if self.isPost() {
	// 	username := strings.TrimSpace(self.GetString("username"))
	// 	password := strings.TrimSpace(self.GetString("password"))
	// 	if username != "" && password != "" {
	// 		user, err := model.AdminGetByName(username)
	// 		if err != nil || user.Password != libs.Md5([]byte(password+user.Salt)) {
	// 			self.ajaxMsg("帐号或密码错误", MSG_ERR)
	// 		} else if user.Status == -1 {
	// 			self.ajaxMsg("该帐号已禁用", MSG_ERR)
	// 		} else {
	// 			user.LastIp = self.getClientIp()
	// 			user.LastLogin = time.Now().Unix()
	// 			user.Update()
	// 			authkey := libs.Md5([]byte(self.getClientIp() + "|" + user.Password + user.Salt))
	// self.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey, 7*86400)
	c.SetCookie("auth", "1|authkey", 7*86400, "/", "10.74.6.23", false, true)

	// 			self.ajaxMsg("登录成功", MSG_OK)
	// 		}
	// 	}
	// }
	// self.ajaxMsg("请求方式错误", MSG_ERR)
	c.JSON(http.StatusOK, common.Success(c))
}

//登出
func (self *LoginController) LoginOut(c *gin.Context) {
	c.SetCookie("auth", "1|authkey", -1, "/", "10.74.6.23", false, true)
	c.Redirect(302, "/home/start")
	// self.Ctx.SetCookie("auth", "")
	// self.redirect(beego.URLFor("LoginController.Login"))
}

func (self *LoginController) NoAuth() {
	self.Ctx.WriteString("没有权限")
}
