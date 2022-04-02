/**********************************************
** @Des: base controller
** @Author: haodaquan
** @Date:   2017-09-07 16:54:40
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-18 10:28:01
***********************************************/
package handler

import (
	"github.com/astaxie/beego"
	"github.com/voioc/cjob/app/model"
)

const (
	MSG_OK  = 0
	MSG_ERR = -1
)

type BaseController struct {
	beego.Controller
	controllerName string
	actionName     string
	user           *model.Admin
	userId         int
	userName       string
	loginName      string
	pageSize       int
	allowUrl       string
	serverGroups   string
	taskGroups     string
}

// //前期准备
// func (self *BaseController) Prepare(c *gin.Context) {
// 	self.pageSize = 20
// 	controllerName, actionName := self.GetControllerAndAction()
// 	self.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
// 	self.actionName = strings.ToLower(actionName)
// 	self.Data["version"] = beego.AppConfig.String("version")
// 	self.Data["siteName"] = beego.AppConfig.String("site.name")
// 	self.Data["curRoute"] = self.controllerName + "." + self.actionName
// 	self.Data["curController"] = self.controllerName
// 	self.Data["curAction"] = self.actionName
// 	// noAuth := "ajaxsave/ajaxdel/table/loginin/loginout/getnodes/start"
// 	// isNoAuth := strings.Contains(noAuth, self.actionName)
// 	//fmt.Println(self.controllerName)
// 	//if (strings.Compare(self.controllerName, "apidoc")) != 0 {
// 	//
// 	//}

// 	self.Auth(c)
// 	self.Data["loginUserId"] = self.userId
// 	self.Data["loginUserName"] = self.userName
// }

// // 登录权限验证
// func (self *BaseController) Auth(c *gin.Context) {
// 	// arr := strings.Split(self.Ctx.GetCookie("auth"), "|")
// 	cookie, _ := c.Cookie("auth")
// 	arr := strings.Split(cookie, "|")

// 	self.userId = 0
// 	if len(arr) == 2 {
// 		idstr, password := arr[0], arr[1]
// 		userId, _ := strconv.Atoi(idstr)
// 		if userId > 0 {
// 			user, err := model.AdminGetById(userId)

// 			if err == nil && password == libs.Md5([]byte(self.getClientIp()+"|"+user.Password+user.Salt)) {
// 				self.userId = user.ID
// 				self.loginName = user.LoginName
// 				self.userName = user.RealName
// 				self.user = user
// 				self.AdminAuth()
// 				// self.dataAuth(user)
// 			}

// 			isHasAuth := strings.Contains(self.allowUrl, self.controllerName+"/"+self.actionName)
// 			noAuth := "ajaxsave/table/loginin/loginout/getnodes/start/apitask/apistart/apipause"
// 			isNoAuth := strings.Contains(noAuth, self.actionName)

// 			if isHasAuth == false && isNoAuth == false {
// 				if strings.Contains(self.actionName, "ajax") {
// 					self.ajaxMsg("没有权限", MSG_ERR)
// 					return
// 				}

// 				flash := beego.NewFlash()
// 				flash.Error("没有权限")
// 				flash.Store(&self.Controller)
// 				return
// 			}
// 		}
// 	}

// 	if self.userId == 0 &&
// 		(self.controllerName != "login" &&
// 			self.actionName != "loginin" &&
// 			self.actionName != "apistart" &&
// 			self.actionName != "apitask" &&
// 			self.actionName != "apipause" &&
// 			self.actionName != "apisave" &&
// 			self.actionName != "apistatus" &&
// 			self.actionName != "apiget") {
// 		self.redirect(beego.URLFor("LoginController.Login"))
// 	}
// }

// // func (self *BaseController) dataAuth(user *model.Admin) {
// // 	if user.RoleIDs == "0" || user.ID == 1 {
// // 		return
// // 	}

// // 	Filters := make([]interface{}, 0)
// // 	Filters = append(Filters, "status", 1)

// // 	RoleIdsArr := strings.Split(user.RoleIDs, ",")

// // 	RoleIds := make([]int, 0)
// // 	for _, v := range RoleIdsArr {
// // 		id, _ := strconv.Atoi(v)
// // 		RoleIds = append(RoleIds, id)
// // 	}

// // 	Filters = append(Filters, "id", RoleIds)

// // 	Result, _, _ := service.RoleS(c).RoleList(1, 1000, Filters...)
// // 	serverGroups := ""
// // 	taskGroups := ""
// // 	for _, v := range Result {
// // 		serverGroups += v.ServerGroupIDs + ","
// // 		taskGroups += v.TaskGroupIDs + ","
// // 	}

// // 	self.serverGroups = strings.Trim(serverGroups, ",")
// // 	self.taskGroups = strings.Trim(taskGroups, ",")
// // }

// func (self *BaseController) AdminAuth() {
// 	// 左侧导航栏
// 	filters := make([]interface{}, 0)
// 	filters = append(filters, "status", 1)
// 	if self.userId != 1 {
// 		//普通管理员
// 		adminAuthIds, _ := model.RoleAuthGetByIds(self.user.RoleIDs)
// 		adminAuthIdArr := strings.Split(adminAuthIds, ",")
// 		filters = append(filters, "id__in", adminAuthIdArr)
// 	}

// 	c := &gin.Context{}
// 	result, _, _ := service.AuthS(c).AuthList(1, 1000, filters...)
// 	list := make([]map[string]interface{}, len(result))
// 	list2 := make([]map[string]interface{}, len(result))
// 	allow_url := ""
// 	i, j := 0, 0
// 	for _, v := range result {
// 		if v.AuthUrl != " " || v.AuthUrl != "/" {
// 			allow_url += v.AuthUrl
// 		}
// 		row := make(map[string]interface{})
// 		if v.PID == 1 && v.IsShow == 1 {
// 			row["Id"] = int(v.ID)
// 			row["Sort"] = v.Sort
// 			row["AuthName"] = v.AuthName
// 			row["AuthUrl"] = v.AuthUrl
// 			row["Icon"] = v.Icon
// 			row["Pid"] = int(v.PID)
// 			list[i] = row
// 			i++
// 		}
// 		if v.PID != 1 && v.IsShow == 1 {
// 			row["Id"] = int(v.ID)
// 			row["Sort"] = v.Sort
// 			row["AuthName"] = v.AuthName
// 			row["AuthUrl"] = v.AuthUrl
// 			row["Icon"] = v.Icon
// 			row["Pid"] = int(v.PID)
// 			list2[j] = row
// 			j++
// 		}
// 	}

// 	self.Data["SideMenu1"] = list[:i]  //一级菜单
// 	self.Data["SideMenu2"] = list2[:j] //二级菜单

// 	self.allowUrl = allow_url + "/home/index"
// }

// // 是否POST提交
// func (self *BaseController) isPost() bool {
// 	return self.Ctx.Request.Method == "POST"
// }

// //获取用户IP地址
// func (self *BaseController) getClientIp() string {
// 	s := strings.Split(self.Ctx.Request.RemoteAddr, ":")
// 	return s[0]
// }

// // 重定向
// func (self *BaseController) redirect(url string) {
// 	self.Redirect(url, 302)
// 	self.StopRun()
// }

// //加载模板
// func (self *BaseController) display(tpl ...string) {
// 	var tplname string
// 	if len(tpl) > 0 {
// 		tplname = strings.Join([]string{tpl[0], "html"}, ".")
// 	} else {
// 		tplname = self.controllerName + "/" + self.actionName + ".html"
// 	}
// 	self.Layout = "public/layout.html"
// 	self.TplName = tplname
// }

// //ajax返回
// func (self *BaseController) ajaxMsg(msg interface{}, msgno int) {
// 	out := make(map[string]interface{})
// 	out["status"] = msgno
// 	out["message"] = msg
// 	self.Data["json"] = out
// 	self.ServeJSON()
// 	self.StopRun()
// }

// //ajax返回 列表
// func (self *BaseController) ajaxList(msg interface{}, msgno int, count int64, data interface{}) {
// 	out := make(map[string]interface{})
// 	out["code"] = msgno
// 	out["msg"] = msg
// 	out["count"] = count
// 	out["data"] = data
// 	self.Data["json"] = out
// 	self.ServeJSON()
// 	self.StopRun()
// }

// //资源分组信息
// func serverGroupLists(authStr string, adminId int) (sgl map[int]string) {
// 	Filters := make([]interface{}, 0)
// 	Filters = append(Filters, "status", 1)
// 	if authStr != "0" && adminId != 1 {
// 		serverGroupIdsArr := strings.Split(authStr, ",")
// 		serverGroupIds := make([]int, 0)
// 		for _, v := range serverGroupIdsArr {
// 			id, _ := strconv.Atoi(v)
// 			serverGroupIds = append(serverGroupIds, id)
// 		}
// 		Filters = append(Filters, "id__in", serverGroupIds)
// 	}

// 	groupResult, n := model.ServerGroupGetList(1, 1000, Filters...)
// 	sgl = make(map[int]string, n)
// 	for _, gv := range groupResult {
// 		sgl[gv.ID] = gv.GroupName
// 	}
// 	//sgl[0] = "默认分组"
// 	return sgl
// }

// func taskGroupLists(authStr string, adminId int) (gl map[int]string) {
// 	groupFilters := make([]interface{}, 0)
// 	groupFilters = append(groupFilters, "status", 1)
// 	if authStr != "0" && adminId != 1 {
// 		taskGroupIdsArr := strings.Split(authStr, ",")
// 		taskGroupIds := make([]int, 0)
// 		for _, v := range taskGroupIdsArr {
// 			id, _ := strconv.Atoi(v)
// 			taskGroupIds = append(taskGroupIds, id)
// 		}
// 		groupFilters = append(groupFilters, "id__in", taskGroupIds)
// 	}
// 	groupResult, n := model.GroupGetList(1, 1000, groupFilters...)
// 	gl = make(map[int]string, n)
// 	for _, gv := range groupResult {
// 		gl[gv.ID] = gv.GroupName
// 	}
// 	return gl
// }

// func serverListByGroupId(groupId int) []string {
// 	Filters := make([]interface{}, 0)
// 	Filters = append(Filters, "status", 1)
// 	Filters = append(Filters, "group_id", groupId)
// 	Result, _ := model.TaskServerGetList(1, 1000, Filters...)

// 	servers := make([]string, 0)
// 	for _, v := range Result {
// 		servers = append(servers, strconv.Itoa(v.ID), v.ServerName)
// 	}

// 	return servers
// }

// type AdminInfo struct {
// 	Id       int
// 	Email    string
// 	Phone    string
// 	RealName string
// }

// func AllAdminInfo(adminIds string) []*AdminInfo {
// 	Filters := make([]interface{}, 0)
// 	Filters = append(Filters, "status", 1)
// 	//Filters = append(Filters, "id__gt", 1)
// 	var notifyUserIds []int
// 	if adminIds != "0" && adminIds != "" {
// 		notifyUserIdsStr := strings.Split(adminIds, ",")
// 		for _, v := range notifyUserIdsStr {
// 			i, _ := strconv.Atoi(v)
// 			notifyUserIds = append(notifyUserIds, i)
// 		}
// 		Filters = append(Filters, "id__in", notifyUserIds)
// 	}
// 	Result, _ := model.AdminGetList(1, 1000, Filters...)

// 	adminInfos := make([]*AdminInfo, 0)
// 	for _, v := range Result {
// 		ai := AdminInfo{
// 			Id:       v.ID,
// 			Email:    v.Email,
// 			Phone:    v.Phone,
// 			RealName: v.RealName,
// 		}
// 		adminInfos = append(adminInfos, &ai)
// 	}

// 	return adminInfos
// }

// type serverList struct {
// 	GroupId   int
// 	GroupName string
// 	Servers   map[int]string
// }

// func serverLists(authStr string, adminId int) (sls []serverList) {
// 	serverGroup := serverGroupLists(authStr, adminId)
// 	Filters := make([]interface{}, 0)
// 	Filters = append(Filters, "status__in", []int{0, 1})

// 	Result, _ := model.TaskServerGetList(1, 1000, Filters...)
// 	for k, v := range serverGroup {
// 		sl := serverList{}
// 		sl.GroupId = k
// 		sl.GroupName = v
// 		servers := make(map[int]string)
// 		for _, sv := range Result {
// 			if sv.GroupID == k {
// 				servers[sv.ID] = sv.ServerName
// 			}
// 		}
// 		sl.Servers = servers
// 		sls = append(sls, sl)
// 	}
// 	return sls
// }
