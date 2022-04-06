package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/app/service"
	"github.com/voioc/cjob/common"
	"github.com/voioc/cjob/utils"
	"github.com/voioc/coco/logzap"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tmpCookie, _ := c.Cookie("job_cookie")
		cookie := strings.Split(tmpCookie, "|")

		key := "123456781234567812345678" // 16,24,32 位的密钥
		ciphertext := cookie[0]
		nonce := cookie[1]
		// ciphertext, nonce := utils.AESGCMEncrypt(plaintext, key)
		// uid, _ := strconv.Atoi(arr[0])

		// fmt.Println("解密结果: ", utils.AESGCMDecrypt(ciphertext, key, nonce))
		uidString := utils.AESGCMDecrypt(ciphertext, key, nonce)
		uid, err := strconv.Atoi(uidString)
		if uid == 0 || err != nil {
			logzap.Ex(context.Background(), "AUTH", "auth failed: %s", uidString)
			c.JSON(http.StatusOK, common.Error(c, common.ERROR_AUTH))
			c.Abort()
			return
		}

		user := model.Admin{}
		if err := model.DataByID(&user, uid); err != nil {
			logzap.Ex(context.Background(), "AUTH", "record not found: %s", err.Error())
			c.JSON(http.StatusOK, common.Error(c, common.ERROR_AUTH))
			c.Abort()
			return
		}

		if user.ID == 0 {
			c.JSON(http.StatusOK, common.Error(c, common.ERROR_AUTH))
			c.Abort()
			return
		}

		tg, sg := service.RoleS(c).Resources(uid, user.RoleIDs)

		// c.Set("menu", data)
		c.Set("uid", uid)
		c.Set("role_ids", user.RoleIDs)
		c.Set("tg", tg) // taskgroups
		c.Set("sg", sg) // taskgroups

		c.Next()
	}
}
