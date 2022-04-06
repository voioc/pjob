/************************************************************
** @Description: PPGo_Job2
** @Author: haodaquan
** @Date:   2018-06-05 22:24
** @Last Modified by:   haodaquan
** @Last Modified time: 2018-06-05 22:24
*************************************************************/
package main

import (
	"fmt"

	"github.com/spf13/viper"
	_ "github.com/voioc/cjob/init"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/service"

	"github.com/voioc/cjob/routers"
	"github.com/voioc/coco/logzap"
)

func main() {

	// db.Init()        // 初始化数据库连接
	// cache.Init()     // 初始化缓存连接
	logzap.InitZap() // 初始化日志组件
	// worker.InitJobs()

	service.TaskS(&gin.Context{}).Loading()

	r := gin.New()

	//加载静态资源文件路径
	r.Static("/static", "./static")

	// r.LoadHTMLFiles("templates/index.html")
	// r.GET("/detail", func(c *gin.Context) {
	// 	c.HTML(200, "index.html", "<a href='lizhouwen.com'>1232</a>")
	// })

	routers.InitRouter(r)

	fmt.Println("The service is running...")

	port := viper.GetString("server.port")
	endless.ListenAndServe(":"+port, r)

	// beego.Run()

}
