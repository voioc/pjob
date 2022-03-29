/************************************************************
** @Description: controllers
** @Author: haodaquan
** @Date:   2018-06-09 16:11
** @Last Modified by:   Bee
** @Last Modified time: 2019-02-17 22:15:15
*************************************************************/
package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/libs"
	"github.com/voioc/cjob/models"
	"github.com/voioc/cjob/service"
	"github.com/voioc/cjob/utils"
)

type ServerController struct {
	BaseController
}

func (self *ServerController) List(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "执行资源管理"
	_, sg := service.TaskGroups(c.GetInt("uid"), c.GetString("role_id"))
	data["serverGroup"] = serverGroupLists(sg, c.GetInt("uid"))
	// self.display()

	c.HTML(http.StatusOK, "server/list.html", data)
}

func (self *ServerController) Add(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "新增执行资源"

	_, sg := service.TaskGroups(c.GetInt("uid"), c.GetString("role_id"))
	data["serverGroup"] = serverGroupLists(sg, c.GetInt("uid"))
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

	_, sg := service.TaskGroups(c.GetInt("uid"), c.GetString("role_id"))
	serverGroup := serverGroupLists(sg, c.GetInt("uid"))

	//查询条件
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 0)
	filters = append(filters, "group_id", gid)

	result, count := models.TaskServerGetList(page, pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["connection_type"] = v.ConnectionType
		row["server_name"] = v.ServerName
		row["detail"] = v.Detail
		if serverGroup[v.GroupId] == "" {
			v.GroupId = 0
		}
		row["group_name"] = serverGroup[v.GroupId]
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
	server, _ := models.TaskServerGetById(id)
	row := make(map[string]interface{})
	row["id"] = server.Id
	row["connection_type"] = server.ConnectionType
	row["server_name"] = server.ServerName
	row["group_id"] = server.GroupId
	row["server_ip"] = server.ServerIp
	row["server_account"] = server.ServerAccount
	row["server_outer_ip"] = server.ServerOuterIp
	row["port"] = server.Port
	row["type"] = server.Type
	row["password"] = server.Password
	row["public_key_src"] = server.PublicKeySrc
	row["private_key_src"] = server.PrivateKeySrc
	row["detail"] = server.Detail
	data["server"] = row

	_, sg := service.TaskGroups(c.GetInt("uid"), c.GetString("role_id"))
	data["serverGroup"] = serverGroupLists(sg, c.GetInt("uid"))
	// self.display()

	c.HTML(http.StatusOK, "server/edit.html", data)
}

func (self *ServerController) AjaxTestServer(c *gin.Context) {

	server := new(models.TaskServer)
	server.ConnectionType, _ = strconv.Atoi(c.DefaultPostForm("connection_type", "0"))
	server.ServerName = strings.TrimSpace(c.DefaultPostForm("server_name", ""))
	server.ServerAccount = strings.TrimSpace(c.DefaultPostForm("server_account", ""))
	server.ServerOuterIp = strings.TrimSpace(c.DefaultPostForm("server_outer_ip", ""))
	server.ServerIp = strings.TrimSpace(c.DefaultPostForm("server_ip", ""))
	server.PrivateKeySrc = strings.TrimSpace(c.DefaultPostForm("private_key_src", ""))
	server.PublicKeySrc = strings.TrimSpace(c.DefaultPostForm("public_key_src", ""))
	server.Password = strings.TrimSpace(c.DefaultPostForm("password", ""))
	server.Detail = strings.TrimSpace(c.DefaultPostForm("detail", ""))
	server.Type, _ = strconv.Atoi(c.DefaultPostForm("type", "0"))
	server.Port, _ = strconv.Atoi(c.DefaultPostForm("port", "0"))
	server.GroupId, _ = strconv.Atoi(c.DefaultPostForm("group_id", "0"))

	var err error

	if server.ConnectionType == 0 {
		if server.Type == 0 {
			//密码登录
			err = libs.RemoteCommandByPassword(server)
		}

		if server.Type == 1 {
			//密钥登录
			err = libs.RemoteCommandByKey(server)
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
	server, _ := models.TaskServerGetById(id)
	row := make(map[string]interface{})
	row["id"] = server.Id
	row["connection_type"] = server.ConnectionType
	row["server_name"] = server.ServerName
	row["group_id"] = server.GroupId
	row["server_ip"] = server.ServerIp
	row["server_account"] = server.ServerAccount
	row["server_outer_ip"] = server.ServerOuterIp
	row["port"] = server.Port
	row["type"] = server.Type
	row["password"] = server.Password
	row["public_key_src"] = server.PublicKeySrc
	row["private_key_src"] = server.PrivateKeySrc
	row["detail"] = server.Detail
	data["server"] = row

	_, sg := service.TaskGroups(c.GetInt("uid"), c.GetString("role_id"))
	data["serverGroup"] = serverGroupLists(sg, c.GetInt("uid"))

	// self.display()
	c.HTML(http.StatusOK, "group/copy.html", data)
}

func (self *ServerController) AjaxSave(c *gin.Context) {
	server_id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))

	if server_id == 0 {
		server := new(models.TaskServer)
		server.ConnectionType, _ = strconv.Atoi(c.DefaultPostForm("connection_type", "0"))
		server.ServerName = strings.TrimSpace(c.DefaultPostForm("server_name", ""))
		server.ServerAccount = strings.TrimSpace(c.DefaultPostForm("server_account", ""))
		server.ServerOuterIp = strings.TrimSpace(c.DefaultPostForm("server_outer_ip", ""))
		server.ServerIp = strings.TrimSpace(c.DefaultPostForm("server_ip", ""))
		server.PrivateKeySrc = strings.TrimSpace(c.DefaultPostForm("private_key_src", ""))
		server.PublicKeySrc = strings.TrimSpace(c.DefaultPostForm("public_key_src", ""))
		server.Password = strings.TrimSpace(c.DefaultPostForm("password", ""))

		server.Detail = strings.TrimSpace(c.DefaultPostForm("detail", ""))
		server.Type, _ = strconv.Atoi(c.DefaultPostForm("type", "0"))
		server.Port, _ = strconv.Atoi(c.DefaultPostForm("port", "0"))
		server.GroupId, _ = strconv.Atoi(c.DefaultPostForm("group_id", "0"))

		server.CreateTime = time.Now().Unix()
		server.UpdateTime = time.Now().Unix()
		server.Status = 0

		if _, err := models.TaskServerAdd(server); err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}
		// self.ajaxMsg("", MSG_OK)
		c.JSON(http.StatusOK, common.Success(c))
		return
	}

	server, _ := models.TaskServerGetById(server_id)

	//修改
	// server.Id = server_id
	server.UpdateTime = time.Now().Unix()

	server.ConnectionType, _ = strconv.Atoi(c.DefaultPostForm("connection_type", "0"))
	server.ServerName = strings.TrimSpace(c.DefaultPostForm("server_name", ""))
	server.ServerAccount = strings.TrimSpace(c.DefaultPostForm("server_account", ""))
	server.ServerOuterIp = strings.TrimSpace(c.DefaultPostForm("server_outer_ip", ""))
	server.ServerIp = strings.TrimSpace(c.DefaultPostForm("server_ip", ""))
	server.PrivateKeySrc = strings.TrimSpace(c.DefaultPostForm("private_key_src", ""))
	server.PublicKeySrc = strings.TrimSpace(c.DefaultPostForm("public_key_src", ""))
	server.Detail = strings.TrimSpace(c.DefaultPostForm("detail", ""))
	server.Password = strings.TrimSpace(c.DefaultPostForm("password", ""))

	server.Type, _ = strconv.Atoi(c.DefaultPostForm("type", "0"))
	server.Port, _ = strconv.Atoi(c.DefaultPostForm("port", "0"))
	server.GroupId, _ = strconv.Atoi(c.DefaultPostForm("group_id", "0"))

	if err := server.Update(); err != nil {
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

	server, _ := models.TaskServerGetById(id)
	if server == nil {
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "资源不存在"))
		return
	}

	server.UpdateTime = time.Now().Unix()
	server.Status = 2
	server.Id = id

	// TODO 查询服务器是否用于定时任务
	if err := server.Update(); err != nil {
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
		"<i class='fa fa-ban' style='color:#FF5722'></i>",
	}
	//
	//loginType := [2]string{
	//	"密码",
	//	"密钥",
	//}

	connectionType := [3]string{
		"SSH",
		"Telnet",
		"Agent",
	}

	uid := c.GetInt("uid")
	_, sg := service.TaskGroups(uid, c.GetString("role_id"))
	serverGroup := serverGroupLists(sg, uid)

	// self.pageSize = limit
	//查询条件
	filters := make([]interface{}, 0)
	ids := []int{0, 1}
	filters = append(filters, "status__in", ids)

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
		filters = append(filters, "group_id__in", groupsIds)
	} else if serverGroupId > 0 {
		groupsIds = append(groupsIds, serverGroupId)
		filters = append(filters, "group_id__in", groupsIds)
	}

	if serverName != "" {
		filters = append(filters, "server_name__icontains", serverName)
	}

	result, count := models.TaskServerGetList(page, pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["connection_type"] = connectionType[v.ConnectionType]
		row["server_name"] = StatusText[v.Status] + " " + v.ServerName
		row["detail"] = v.Detail
		if serverGroup[v.GroupId] == "" {
			v.GroupId = 0
		}
		row["ip_port"] = v.ServerIp + ":" + strconv.Itoa(v.Port)
		row["group_name"] = serverGroup[v.GroupId]
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
	serverIp := strings.TrimSpace(c.DefaultPostForm("server_ip", ""))
	port, _ := strconv.Atoi(c.DefaultPostForm("port", "0"))

	if serverIp == "" || port == 0 {
		// self.ajaxMsg("执行器和端口号必填", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "执行器和端口号必填"))
		return
	}

	defaultActName := "agent-" + serverIp + "-" + strconv.Itoa(port)

	id := models.TaskServerForActuator(serverIp, port)
	if id == 0 {
		//新增
		server := new(models.TaskServer)
		server.ConnectionType, _ = strconv.Atoi(c.DefaultPostForm("connection_type", "3"))
		server.ServerName = strings.TrimSpace(c.DefaultPostForm("server_name", defaultActName))
		server.ServerAccount = strings.TrimSpace(c.DefaultPostForm("server_account", "agent"))
		server.ServerOuterIp = strings.TrimSpace(c.DefaultPostForm("server_outer_ip", ""))
		server.ServerIp = strings.TrimSpace(c.DefaultPostForm("server_ip", ""))
		server.PrivateKeySrc = strings.TrimSpace(c.DefaultPostForm("private_key_src", ""))
		server.PublicKeySrc = strings.TrimSpace(c.DefaultPostForm("public_key_src", ""))
		server.Password = strings.TrimSpace(c.DefaultPostForm("password", "agent"))

		server.Detail = strings.TrimSpace(c.DefaultPostForm("detail", ""))
		server.Type, _ = strconv.Atoi(c.DefaultPostForm("type", "0"))
		server.Port, _ = strconv.Atoi(c.DefaultPostForm("port", "0"))
		server.GroupId, _ = strconv.Atoi(c.DefaultPostForm("group_id", "0"))
		server.Status = 0

		server.CreateTime = time.Now().Unix()
		server.UpdateTime = time.Now().Unix()
		server.Status = 0
		serverId, err := models.TaskServerAdd(server)
		if err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}
		// self.ajaxMsg(serverId, MSG_OK)
		data := map[string]interface{}{"server": serverId}
		c.JSON(http.StatusOK, common.Success(c, data))
		return
	} else {
		//修改状态
		server, _ := models.TaskServerGetById(id)
		server.UpdateTime = time.Now().Unix()
		server.Status, _ = strconv.Atoi(c.DefaultPostForm("status", "0"))
		if err := server.Update(); err != nil {
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}

		// self.ajaxMsg(id, MSG_OK)
		data := map[string]interface{}{"server": id}
		c.JSON(http.StatusOK, common.Success(c, data))
		return
	}

}

//检测0-正常，1-异常，2-删除
func (self *ServerController) ApiStatus(c *gin.Context) {
	//唯一确定值 ip+port
	serverId := strings.TrimSpace(c.DefaultPostForm("server_ip", ""))
	port, _ := strconv.Atoi(c.DefaultPostForm("port", "0"))
	status, _ := strconv.Atoi(c.DefaultPostForm("status", "0"))

	if serverId == "" || port == 0 {
		// self.ajaxMsg("执行器和端口号必填", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "执行器和端口号必填"))
		return
	}

	id := models.TaskServerForActuator(serverId, port)
	if id == 0 {
		// self.ajaxMsg("执行器不存在", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "执行器不存在"))
		return
	}

	if status != 0 && status != 1 {
		status = 0
	}

	server, _ := models.TaskServerGetById(id)
	server.UpdateTime = time.Now().Unix()
	server.Status = status
	server.Id = id

	logs.Info(server)

	//TODO 查询执行器是否正在使用中
	if err := server.Update(); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return

	}

	// self.ajaxMsg(id, MSG_OK)
	data := map[string]interface{}{"server": id}
	c.JSON(http.StatusOK, common.Success(c, data))
}

//获取 不检测执行器状态
func (self *ServerController) ApiGet(c *gin.Context) {
	//唯一确定值 ip+port
	serverId := strings.TrimSpace(c.DefaultPostForm("server_ip", ""))
	port, _ := strconv.Atoi(c.DefaultPostForm("port", "0"))

	if serverId == "" || port == 0 {
		// self.ajaxMsg("执行器和端口号必填", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "执行器和端口号必填"))
		return
	}

	id := models.TaskServerForActuator(serverId, port)
	if id == 0 {
		// self.ajaxMsg("执行器不存在", MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, "执行器不存在"))
		return
	}

	server, err := models.TaskServerGetById(id)

	if err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	} else {
		// self.ajaxMsg(server, MSG_OK)
	}

	data := map[string]interface{}{"server": server}
	c.JSON(http.StatusOK, common.Success(c, data))
}
