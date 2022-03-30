package middleware

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/app/service"
)

func Menu() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, _ := c.Cookie("auth")
		fmt.Println("cookie: ", cookie)
		arr := strings.Split(cookie, "|")
		fmt.Println("arr:", arr)
		uid := 1
		var user *model.Admin
		if len(arr) == 2 {
			// idstr, password := arr[0], arr[1]
			uid, _ := strconv.Atoi(arr[0])
			if uid > 0 {
				user, _ = service.AdminGetById(uid)
			}
		}

		// c.Set("menu", data)
		c.Set("uid", uid)
		c.Set("role_ids", user.RoleIds)
		c.Next()
	}

}
