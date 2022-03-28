/**********************************************
** @Des: 权限因子
** @Author: haodaquan
** @Date:   2017-09-09 16:14:31
** @Last Modified by:   haodaquan
** @Last Modified time: 2017-09-17 11:23:40
***********************************************/

package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/models"
	"github.com/voioc/cjob/utils"
)

type AuthController struct {
	BaseController
}

func (self *AuthController) Index(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["pageTitle"] = "权限因子"

	// self.display()
	c.HTML(http.StatusOK, "auth/list.html", data)
}

func (self *AuthController) List(c *gin.Context) {
	data := map[string]interface{}{}
	data["uri"] = utils.URI("")

	data["zTree"] = true // 引入ztreecss
	data["pageTitle"] = "权限因子"
	// self.display()

	c.HTML(http.StatusOK, "auth/list.html", data)
}

//获取全部节点
func (self *AuthController) GetNodes(c *gin.Context) {
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)
	result, count := models.AuthGetList(1, 1000, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["pId"] = v.Pid
		row["name"] = v.AuthName
		row["open"] = true
		list[k] = row
	}

	// self.ajaxList("成功", MSG_OK, count, list)
	ext := map[string]int{"count": int(count)}
	c.JSON(http.StatusOK, common.Success(c, list, ext))
}

//获取一个节点
func (self *AuthController) GetNode(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	result, _ := models.AuthGetById(id)
	// if err == nil {
	// 	self.ajaxMsg(err.Error(), MSG_ERR)
	// }
	row := make(map[string]interface{})
	row["id"] = result.Id
	row["pid"] = result.Pid
	row["auth_name"] = result.AuthName
	row["auth_url"] = result.AuthUrl
	row["sort"] = result.Sort
	row["is_show"] = result.IsShow
	row["icon"] = result.Icon

	fmt.Println(row)

	// self.ajaxList("成功", MSG_OK, 0, row)
	ext := map[string]int{"count": 0}
	c.JSON(http.StatusOK, common.Success(c, row, ext))
}

//新增或修改
func (self *AuthController) AjaxSave(c *gin.Context) {

	uid := c.GetInt("uid")
	auth := new(models.Auth)
	auth.UserId = uid
	auth.Pid, _ = strconv.Atoi(c.DefaultPostForm("pid", "0"))
	auth.AuthName = strings.TrimSpace(c.DefaultPostForm("auth_name", ""))
	auth.AuthUrl = strings.TrimSpace(c.DefaultPostForm("auth_url", ""))
	auth.Sort, _ = strconv.Atoi(c.DefaultPostForm("sort", "0"))
	auth.IsShow, _ = strconv.Atoi(c.DefaultPostForm("is_show", "0"))
	auth.Icon = strings.TrimSpace(c.DefaultPostForm("icon", ""))
	auth.UpdateTime = time.Now().Unix()

	auth.Status = 1

	id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	if id == 0 {
		//新增
		auth.CreateTime = time.Now().Unix()
		auth.CreateId = uid
		auth.UpdateId = uid
		if _, err := models.AuthAdd(auth); err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}
	} else {
		auth.Id = id
		auth.UpdateId = self.userId
		if err := auth.Update(); err != nil {
			// self.ajaxMsg(err.Error(), MSG_ERR)
			c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
			return
		}
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}

//删除
func (self *AuthController) AjaxDel(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	auth, err := models.AuthGetById(id)
	if err != nil || auth == nil {
		msg := "角色ID错误"
		if err != nil {
			msg = err.Error()
		}

		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, msg))
		return
	}

	auth.Id = id
	auth.Status = 0
	if err := auth.Update(); err != nil {
		// self.ajaxMsg(err.Error(), MSG_ERR)
		c.JSON(http.StatusOK, common.Error(c, MSG_ERR, err.Error()))
		return
	}

	// self.ajaxMsg("", MSG_OK)
	c.JSON(http.StatusOK, common.Success(c))
}
