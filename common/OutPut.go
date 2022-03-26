/*
 * @Author: Cedar
 * @Date: 2020-05-08 14:34:28
 * @LastEditors: Jianxuesong
 * @LastEditTime: 2021-06-15 15:00:01
 * @FilePath: /Melon/app/common/OutPutCommon.go
 */
package common

import (
	"github.com/gin-gonic/gin"
	"github.com/voioc/coco/public"
)

// Error 通用错误的错误返回
func Error(c *gin.Context, code int, message ...string) gin.H {
	msg := ""
	if mm, flag := INFO[code]; flag {
		msg = mm
	}

	if len(message) > 0 {
		msg = message[0]
	}

	return SetOutput(c, code, msg, nil, nil)
}

// Success 通用的成功返回
func Success(c *gin.Context, params ...interface{}) gin.H {
	var data, ext interface{}
	if len(params) > 0 {
		data = params[0]
	}

	if len(params) > 1 {
		ext = params[0]
	}

	return SetOutput(c, STATUS_OK, "success", data, ext)
}

// SetOutput 自定义返回结构，自己构造
func SetOutput(c *gin.Context, code int, msg string, data, ext interface{}, params ...map[string]interface{}) gin.H {
	// RealCode := http.StatusOK
	result := gin.H{"code": code, "msg": msg, "data": data}

	if data != nil {
		result["data"] = data
	}

	if ext != nil {
		result["ext"] = ext
	}

	// 兼容多余参数
	if len(params) > 0 {
		for key, value := range params[0] {
			result[key] = value
		}
	}

	base := Base{C: c}
	base.SetDebug(1, "[ALL] cost: %s)", public.TimeCost(c.GetTime("start")))

	if c.GetBool("_debug") {
		result["_debug"] = c.GetStringSlice("debug")
	}

	// output, _ := jsoniter.MarshalToString(result)
	// c.Set("output", output)

	return result
}
