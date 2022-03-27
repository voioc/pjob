package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/voioc/pjob/models"
	"github.com/voioc/pjob/service"
)

func Menu() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, _ := c.Cookie("auth")
		arr := strings.Split(cookie, "|")

		uid := 1
		var user *models.Admin
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
