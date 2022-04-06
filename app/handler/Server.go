/************************************************************
** @Description: controllers
** @Author: haodaquan
** @Date:   2018-06-09 16:11
** @Last Modified by:   Bee
** @Last Modified time: 2019-02-17 22:15:15
*************************************************************/
package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/app/service"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/libs"
	"github.com/voioc/cjob/utils"
)

type ServerController struct {
	BaseController
}

func (self *ServerController) List(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "执行资源管理"
	_, sg := service.AuthS(c).TaskGroups(c.GetInt("uid"), c.GetString("role_id"))
	data["serverGroup"], _ = service.ServerGroupS(c).GroupIDName(sg) // serverGroupLists(sg, c.GetInt("uid"))
	// self.display()

	c.HTML(http.StatusOK, "server/list.html", data)
}

func (self *ServerController) Add(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "新增执行资源"

	_, sg := service.AuthS(c).TaskGroups(c.GetInt("uid"), c.GetString("role_id"))
	data["serverGroup"], _ = service.ServerGroupS(c).GroupIDName(sg) // serverGroupLists(sg, c.GetInt("uid"))
	// self.display()

	c.HTML(http.StatusOK, "server/add.html", data)
}

func (self *ServerController) GetServerByGroupId(c *gin.Context) {
	gid, _ := strconv.Atoi(c.DefaultQuery("gid", "0"))
	if gid == 0 {
		// self.ajaxMsg("groupId is not exist", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "groupId is not exist"))
		return
	}

	//列表
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pagesize", "20"))

	//serverName := strings.TrimSpace(self.GetString("serverName"))
	StatusText := []string{
		"正常",
		"<font color='red'>禁用</font>",
	}

	loginType := [2]string{
		"密码",
		"密钥",
	}

	_, sg := service.AuthS(c).TaskGroups(c.GetInt("uid"), c.GetString("role_id"))
	serverGroup, _ := service.ServerGroupS(c).GroupIDName(sg) // serverGroupLists(sg, c.GetInt("uid"))

	//查询条件
	filters := make([]interface{}, 0)
	filters = append(filters, "status =", 0)
	filters = append(filters, "group_id =", gid)

	// result, count, _ := service.ServerS(c).ServerList(page, pageSize, filters...) // model.TaskServerGetList(page, pageSize, filters...)
	result := make([]model.TaskServer, 0)
	if err := model.List(&result, page, pageSize, filters...); err != nil {
		fmt.Println(err.Error())
	}

	count, err := model.ListCount(&model.TaskServer{}, filters...)
	if err != nil {
		fmt.Println(err.Error())
	}

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["connection_type"] = v.ConnectionType
		row["server_name"] = v.ServerName
		row["detail"] = v.Detail
		if serverGroup[v.GroupID] == "" {
			v.GroupID = 0
		}
		row["group_name"] = serverGroup[v.GroupID]
		row["type"] = loginType[v.Type]
		row["status"] = v.Status
		row["status_text"] = StatusText[v.Status]
		list[k] = row
	}

	// self.ajaxList("成功", MSG_OK, count, list)
	ext := map[string]int{"total": int(count)}
	c.JSON(http.StatusOK, common.Success(c, list, ext))
}

func (self *ServerController) Edit(c *gin.Context) {
	data := map[string]interface{}{}
	data["pageTitle"] = "编辑执行资源"

	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	// server, _ := service.ServerS(c).ServerByID(id) // model.TaskServerGetById(id)
	server := model.TaskServer{}
	if err := model.DataByID(&server, id); err != nil {
		fmt.Println(err.Error())
	}

	row := make(map[string]interface{})
	row["id"] = server.ID
	row["connection_type"] = server.ConnectionType
	row["server_name"] = server.ServerName
	row["group_id"] = server.GroupID
	row["server_ip"] = server.ServerIP
	row["server_account"] = server.ServerAccount
	row["server_outer_ip"] = server.ServerOuterIP
	row["port"] = server.Port
	row["type"] = server.Type
	row["password"] = server.Password
	row["public_key_src"] = server.PublicKeySrc
	row["private_key_src"] = server.PrivateKeySrc
	row["detail"] = server.Detail
	data["server"] = row

	_, sg := service.AuthS(c).TaskGroups(c.GetInt("uid"), c.GetString("role_id"))
	data["serverGroup"], _ = service.ServerGroupS(c).GroupIDName(sg) // serverGroupLists(sg, c.GetInt("uid"))
	// self.display()

	c.HTML(http.StatusOK, "server/edit.html", data)
}

func (self *ServerController) AjaxTestServer(c *gin.Context) {

	server := new(model.TaskServer)
	server.ConnectionType, _ = strconv.Atoi(c.DefaultPostForm("connection_type", "0"))
	server.ServerName = strings.TrimSpace(c.DefaultPostForm("server_name", ""))
	server.ServerAccount = strings.TrimSpace(c.DefaultPostForm("server_account", ""))
	server.ServerOuterIP = strings.TrimSpace(c.DefaultPostForm("server_outer_ip", ""))
	server.ServerIP = strings.TrimSpace(c.DefaultPostForm("server_ip", ""))
	server.PrivateKeySrc = strings.TrimSpace(c.DefaultPostForm("private_key_src", ""))
	server.PublicKeySrc = strings.TrimSpace(c.DefaultPostForm("public_key_src", ""))
	server.Password = strings.TrimSpace(c.DefaultPostForm("password", ""))
	server.Detail = strings.TrimSpace(c.DefaultPostForm("detail", ""))
	server.Type, _ = strconv.Atoi(c.DefaultPostForm("type", "0"))
	server.Port, _ = strconv.Atoi(c.DefaultPostForm("port", "0"))
	server.GroupID, _ = strconv.Atoi(c.DefaultPostForm("group_id", "0"))

	var err error

	if server.ConnectionType == 0 {
		if server.Type == 0 {
			// 密码登录
			if err = libs.RemoteCommandByPassword(server); err != nil {
				fmt.Println(err.Error())
			}
		}

		if server.Type == 1 {
			//密钥登录
			if err = libs.RemoteCommandByKey(server); err != nil {
				fmt.Println(err.Error())
			}
		}

		if err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}

		// self.ajaxMsg("Success", MSG_OK)
		c.JSON(http.StatusOK, common.Success(c))
		return
	} else if server.ConnectionType == 1 {
		if server.Type != 0 {
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "Telnet方式暂不支持密钥登陆！"))
			return
		}

		if err = libs.RemoteCommandByTelnetPassword(server); err != nil {
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}

		c.JSON(http.StatusOK, common.Success(c))
		return
	} else if server.ConnectionType == 2 {
		if err := libs.RemoteAgent(server); err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}

		c.JSON(http.StatusOK, common.Success(c))
		return
	}

	// self.ajaxMsg("未知连接方式", MSG_ERR)
	c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "未知连接方式"))
}

func (self *ServerController) Copy(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "复制服务器资源"

	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	// server, _ := service.ServerS(c).ServerByID(id) // model.TaskServerGetById(id)
	server := model.TaskServer{}
	if err := model.DataByID(&server, id); err != nil {
		fmt.Println(err.Error())
	}

	row := make(map[string]interface{})
	row["id"] = server.ID
	row["connection_type"] = server.ConnectionType
	row["server_name"] = server.ServerName
	row["group_id"] = server.GroupID
	row["server_ip"] = server.ServerIP
	row["server_account"] = server.ServerAccount
	row["server_outer_ip"] = server.ServerOuterIP
	row["port"] = server.Port
	row["type"] = server.Type
	row["password"] = server.Password
	row["public_key_src"] = server.PublicKeySrc
	row["private_key_src"] = server.PrivateKeySrc
	row["detail"] = server.Detail
	data["server"] = row

	_, sg := service.AuthS(c).TaskGroups(c.GetInt("uid"), c.GetString("role_id"))
	data["serverGroup"], _ = service.ServerGroupS(c).GroupIDName(sg) // serverGroupLists(sg, c.GetInt("uid"))

	// self.display()
	c.HTML(http.StatusOK, "group/copy.html", data)
}

func (self *ServerController) AjaxSave(c *gin.Context) {
	server_id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))

	if server_id == 0 {
		server := new(model.TaskServer)
		server.ConnectionType, _ = strconv.Atoi(c.DefaultPostForm("connection_type", "0"))
		server.ServerName = strings.TrimSpace(c.DefaultPostForm("server_name", ""))
		server.ServerAccount = strings.TrimSpace(c.DefaultPostForm("server_account", ""))
		server.ServerOuterIP = strings.TrimSpace(c.DefaultPostForm("server_outer_ip", ""))
		server.ServerIP = strings.TrimSpace(c.DefaultPostForm("server_ip", ""))
		server.PrivateKeySrc = strings.TrimSpace(c.DefaultPostForm("private_key_src", ""))
		server.PublicKeySrc = strings.TrimSpace(c.DefaultPostForm("public_key_src", ""))
		server.Password = strings.TrimSpace(c.DefaultPostForm("password", ""))

		server.Detail = strings.TrimSpace(c.DefaultPostForm("detail", ""))
		server.Type, _ = strconv.Atoi(c.DefaultPostForm("type", "0"))
		server.Port, _ = strconv.Atoi(c.DefaultPostForm("port", "0"))
		server.GroupID, _ = strconv.Atoi(c.DefaultPostForm("group_id", "0"))

		server.CreatedAt = time.Now().Unix()
		server.UpdatedAt = time.Now().Unix()
		server.Status = 0

		if err := model.Add(server); err != nil { // model.TaskServerAdd(server); err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}
		// self.ajaxMsg("", MSG_OK)
		c.JSON(http.StatusOK, common.Success(c))
		return
	}

	// server, _ := service.ServerS(c).ServerByID(server_id) // model.TaskServerGetById(server_id)
	server := model.TaskServer{}
	if err := model.DataByID(&server, server_id); err != nil {
		fmt.Println(err.Error())
	}

	//修改
	// server.Id = server_id
	server.UpdatedAt = time.Now().Unix()

	server.ConnectionType, _ = strconv.Atoi(c.DefaultPostForm("connection_type", "0"))
	server.ServerName = strings.TrimSpace(c.DefaultPostForm("server_name", ""))
	server.ServerAccount = strings.TrimSpace(c.DefaultPostForm("server_account", ""))
	server.ServerOuterIP = strings.TrimSpace(c.DefaultPostForm("server_outer_ip", ""))
	server.ServerIP = strings.TrimSpace(c.DefaultPostForm("server_ip", ""))
	server.PrivateKeySrc = strings.TrimSpace(c.DefaultPostForm("private_key_src", ""))
	server.PublicKeySrc = strings.TrimSpace(c.DefaultPostForm("public_key_src", ""))
	server.Detail = strings.TrimSpace(c.DefaultPostForm("detail", ""))
	server.Password = strings.TrimSpace(c.DefaultPostForm("password", ""))

	server.Type, _ = strconv.Atoi(c.DefaultPostForm("type", "0"))
	server.Port, _ = strconv.Atoi(c.DefaultPostForm("port", "0"))
	server.GroupID, _ = strconv.Atoi(c.DefaultPostForm("group_id", "0"))

	if err := model.Update(server.ID, &server); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *ServerController) AjaxDel(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	if id == 1 {
		// self.ajaxMsg("默认分组id=1，禁止删除", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "默认分组id=1，禁止删除"))
		return
	}

	// server, _ := service.ServerS(c).ServerByID(id) //  model.TaskServerGetById(id)
	server := model.TaskServer{}
	if err := model.DataByID(&server, id); err != nil {
		fmt.Println(err.Error())
	}
	if server.ID == 0 {
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "资源不存在"))
		return
	}

	server.UpdatedAt = time.Now().Unix()
	server.Status = 2
	server.ID = id

	// TODO 查询服务器是否用于定时任务
	if err := model.Update(server.ID, &server); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}
	// self.ajaxMsg("操作成功", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

func (self *ServerController) Table(c *gin.Context) {
	//列表
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pagesize", "20"))

	serverGroupId, _ := strconv.Atoi(c.DefaultQuery("serverGroupId", "0"))

	serverName := strings.TrimSpace(c.DefaultQuery("serverName", ""))
	StatusText := []string{
		"<i class='fa fa-refresh' style='color:#5FB878'></i>",
		"<i class='fa fa-bolt' style='color:#5FB878'></i>",
		"<i class='fa fa-ban' style='color:#FF5722'></i>",
	}
	//
	//loginType := [2]string{
	//	"密码",
	//	"密钥",
	//}

	connectionType := [4]string{
		"Unknown",
		"SSH",
		"Telnet",
		"Agent",
	}

	uid := c.GetInt("uid")
	_, sg := service.AuthS(c).TaskGroups(uid, c.GetString("role_id"))
	serverGroup, _ := service.ServerGroupS(c).GroupIDName(sg) // serverGroupLists(sg, uid)

	// self.pageSize = limit
	// 查询条件
	filters := make([]interface{}, 0)
	ids := []int{0, 1}
	filters = append(filters, "status", ids)

	groupsIds := make([]int, 0)
	if uid != 1 {
		groups := strings.Split(sg, ",")

		for _, v := range groups {
			id, _ := strconv.Atoi(v)
			if serverGroupId > 0 {
				if id == serverGroupId {
					groupsIds = append(groupsIds, id)
					break
				}
			} else {
				groupsIds = append(groupsIds, id)
			}
		}
		filters = append(filters, "group_id ", groupsIds)
	} else if serverGroupId > 0 {
		groupsIds = append(groupsIds, serverGroupId)
		filters = append(filters, "group_id ", groupsIds)
	}

	if serverName != "" {
		filters = append(filters, "server_name LIKE '%"+serverName+"%'", "")
	}

	// result, count, _ := service.ServerS(c).ServerList(page, pageSize, filters...) // model.TaskServerGetList(page, pageSize, filters...)
	result := make([]model.TaskServer, 0)
	if err := model.List(&result, page, pageSize, filters...); err != nil {
		fmt.Println(err.Error())
	}

	count, err := model.ListCount(&model.TaskServer{}, filters...)
	if err != nil {
		fmt.Println(err.Error())
	}

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["connection_type"] = connectionType[v.ConnectionType]
		row["server_name"] = StatusText[v.Status] + " " + v.ServerName
		row["detail"] = v.Detail
		if serverGroup[v.GroupID] == "" {
			v.GroupID = 0
		}
		row["ip_port"] = v.ServerIP + ":" + strconv.Itoa(v.Port)
		row["group_name"] = serverGroup[v.GroupID]
		//row["type"] = loginType[v.Type]
		row["status"] = v.Status
		list[k] = row
	}

	// self.ajaxList("成功", MSG_OK, count, list)
	ext := map[string]int{"total": int(count)}
	c.JSON(http.StatusOK, common.Success(c, list, ext))
}

// 以下函数为执行器接口
//注册
func (self *ServerController) ApiSave(c *gin.Context) {
	// 唯一确定值 ip+port
	serverIP := strings.TrimSpace(c.DefaultPostForm("server_ip", ""))
	port, _ := strconv.Atoi(c.DefaultPostForm("port", "0"))

	if serverIP == "" || port == 0 {
		// self.ajaxMsg("执行器和端口号必填", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "执行器和端口号必填"))
		return
	}

	defaultActName := "agent-" + serverIP + "-" + strconv.Itoa(port)

	filters := make([]interface{}, 0)
	filters = append(filters, "status", []int{0, 1})
	filters = append(filters, "server_ip", serverIP)
	filters = append(filters, "port", port)

	// server, _, _ := service.ServerS(c).ServerList(1, 1, filters...) // TaskServerGetList(1, 1, serverFilters...)
	server := make([]model.TaskServer, 0)
	if err := model.List(&server, 1, 1, filters...); err != nil {
		fmt.Println(err.Error())
	}

	// id := model.TaskServerForActuator(serverIp, port)
	if len(server) == 0 {
		//新增
		server := new(model.TaskServer)
		server.ConnectionType, _ = strconv.Atoi(c.DefaultPostForm("connection_type", "3"))
		server.ServerName = strings.TrimSpace(c.DefaultPostForm("server_name", defaultActName))
		server.ServerAccount = strings.TrimSpace(c.DefaultPostForm("server_account", "agent"))
		server.ServerOuterIP = strings.TrimSpace(c.DefaultPostForm("server_outer_ip", ""))
		server.ServerIP = strings.TrimSpace(c.DefaultPostForm("server_ip", ""))
		server.PrivateKeySrc = strings.TrimSpace(c.DefaultPostForm("private_key_src", ""))
		server.PublicKeySrc = strings.TrimSpace(c.DefaultPostForm("public_key_src", ""))
		server.Password = strings.TrimSpace(c.DefaultPostForm("password", "agent"))

		server.Detail = strings.TrimSpace(c.DefaultPostForm("detail", ""))
		server.Type, _ = strconv.Atoi(c.DefaultPostForm("type", "0"))
		server.Port, _ = strconv.Atoi(c.DefaultPostForm("port", "0"))
		server.GroupID, _ = strconv.Atoi(c.DefaultPostForm("group_id", "0"))
		server.Status = 0

		server.CreatedAt = time.Now().Unix()
		server.UpdatedAt = time.Now().Unix()
		server.Status = 0

		if err := model.Add(&server); err != nil { // model.TaskServerAdd(server)
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}
		// self.ajaxMsg(serverId, MSG_OK)

		data := map[string]interface{}{"server": server.ID}
		c.JSON(http.StatusOK, common.Success(c, data))
		return
	} else {
		//修改状态
		// server, _ := service.ServerS(c).ServerByID(server[0].ID) // model.TaskServerGetById(id)
		server[0].UpdatedAt = time.Now().Unix()
		server[0].Status, _ = strconv.Atoi(c.DefaultPostForm("status", "0"))
		if err := model.Update(server[0].ID, server[0], true); err != nil { // server.Update(); err != nil {
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}

		// self.ajaxMsg(id, MSG_OK)
		data := map[string]interface{}{"server": server[0].ID}
		c.JSON(http.StatusOK, common.Success(c, data))
		return
	}

}

//检测0-正常，1-异常，2-删除
func (self *ServerController) ApiStatus(c *gin.Context) {
	//唯一确定值 ip+port
	serverIP := strings.TrimSpace(c.DefaultPostForm("server_ip", ""))
	port, _ := strconv.Atoi(c.DefaultPostForm("port", "0"))
	status, _ := strconv.Atoi(c.DefaultPostForm("status", "0"))

	if serverIP == "" || port == 0 {
		// self.ajaxMsg("执行器和端口号必填", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "执行器和端口号必填"))
		return
	}

	filters := make([]interface{}, 0)
	filters = append(filters, "status", []int{0, 1})
	filters = append(filters, "server_ip =", serverIP)
	filters = append(filters, "port =", port)

	// server, _, _ := service.ServerS(c).ServerList(1, 1, filters...) // TaskServerGetList(1, 1, serverFilters...)
	server := make([]model.TaskServer, 0)
	if err := model.List(&server, 1, 1, filters...); err != nil {
		fmt.Println(err.Error())
	}

	if len(server) == 0 {
		// self.ajaxMsg("执行器不存在", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "执行器不存在"))
		return
	}

	if status != 0 && status != 1 {
		status = 0
	}

	// server, _ := service.ServerS(c).ServerByID(server[0].ID) // model.TaskServerGetById(id)
	server[0].UpdatedAt = time.Now().Unix()
	server[0].Status = status
	// server[0].ID = server[0].ID

	logs.Info(server)

	//TODO 查询执行器是否正在使用中
	if err := model.Update(server[0].ID, server[0]); err != nil { // server.Update(); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return

	}

	// self.ajaxMsg(id, MSG_OK)
	data := map[string]interface{}{"server": server[0].ID}
	c.JSON(http.StatusOK, common.Success(c, data))
}

//获取 不检测执行器状态
func (self *ServerController) ApiGet(c *gin.Context) {
	//唯一确定值 ip+port
	serverIP := strings.TrimSpace(c.DefaultPostForm("server_ip", ""))
	port, _ := strconv.Atoi(c.DefaultPostForm("port", "0"))

	if serverIP == "" || port == 0 {
		// self.ajaxMsg("执行器和端口号必填", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "执行器和端口号必填"))
		return
	}

	filters := make([]interface{}, 0)
	filters = append(filters, "server_ip =", serverIP)
	filters = append(filters, "port =", port)
	// server, _, _ := service.ServerS(c).ServerList(1, 1, filters...) // TaskServerGetList(1, 1, serverFilters...)
	server := make([]model.TaskServer, 0)
	if err := model.List(&server, 1, 1, filters...); err != nil {
		fmt.Println(err.Error())
	}

	if len(server) == 0 {
		// self.ajaxMsg("执行器不存在", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "执行器不存在"))
		return
	}

	// data := map[string]interface{}{"server": server[0]}
	c.JSON(http.StatusOK, common.Success(c, server[0]))
}
