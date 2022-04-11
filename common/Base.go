package common

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var env string = "dev"

func init() {
	viper.SetDefault("envType", env)

	// 获取当前环境变量
	if realEnv := strings.ToLower(os.Getenv("envType")); realEnv != "" {
		env = realEnv
		viper.SetDefault("envType", env)
	}

	path, _ := filepath.Abs(filepath.Dir(""))        // 获取当前路径
	conf := path + "/config/config_" + env + ".toml" // 拼接配置文件
	// fmt.Println(conf)

	configFile := flag.String("c", conf, "配置文件路径") // 手动置顶配置文件
	flag.Parse()

	viper.SetConfigFile(*configFile) // 读取配置文件
	fmt.Println("Loading config file " + *configFile)

	if err := viper.ReadInConfig(); err != nil { //是否读取成功
		log.Fatalln("打开配置文件失败", err)
	}
}

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
