package common

/*
 * @Description: Do not edit
 * @Author: Jianxuesong
 * @Date: 2021-06-11 14:31:23
 * @LastEditors: Jianxuesong
 * @LastEditTime: 2021-06-11 14:31:54
 * @FilePath: /Melon/app/common/Base.go
 */

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Base struct {
	C *gin.Context
}

// Flush Flush
func (bm *Base) Flush() bool {
	return bm.C.GetBool("_flush")
	// return false
}

// Debug debug
func (bm *Base) Debug() bool {
	return bm.C.GetBool("_debug")
}

// SetDebug 设置debug
func (bm *Base) SetDebug(depth int, message string, a ...interface{}) {
	if !bm.C.GetBool("_debug") {
		return
	}

	if depth == 0 {
		depth = 1
	}

	_, file, line, _ := runtime.Caller(depth)
	path := strings.LastIndexByte(file, '/')

	info := fmt.Sprintf(message, a...)
	tmp := string([]byte(file)[path+1:]) + "(line " + strconv.Itoa(line) + "): " + info

	debug := bm.C.GetStringSlice("debug")
	bm.C.Set("debug", append(debug, tmp))
}

// TimeCost 计算花费时间
func (bm *Base) TimeCost(key string) string {
	cost := ""
	if start := bm.C.GetTime(key); start.IsZero() {
		bm.C.Set(key, time.Now())
	} else {
		tc := time.Since(start)
		cost = fmt.Sprintf("%v", tc)
	}

	return cost
}
