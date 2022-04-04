package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/app/service"
	"github.com/voioc/cjob/common"
)

func Menu() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, _ := c.Cookie("auth")
		fmt.Println("cookie: ", cookie)
		arr := strings.Split(cookie, "|")
		fmt.Println("arr:", arr)
		// uid, _ := strconv.Atoi(arr[0])
		uid := 1
		// user, err := service.AdminS(c).AdminGetByID(uid)
		user := &model.Admin{}
		if err := model.DataByID(user, uid); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusOK, common.Error(c, common.ERROR_AUTH))
			c.Abort()
			return
		}

		if user.ID == 0 {
			c.JSON(http.StatusOK, common.Error(c, common.ERROR_AUTH))
			c.Abort()
			return
		}

		tg, sg := service.RoleS(c).TaskGroups(uid, user.RoleIDs)

		// c.Set("menu", data)
		c.Set("uid", uid)
		c.Set("role_ids", user.RoleIDs)
		c.Set("tg", tg) // taskgroups
		c.Set("sg", sg) // taskgroups

		c.Next()
	}
}
