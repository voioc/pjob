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
	"time"

	_ "github.com/voioc/cjob/init"
	"github.com/voioc/cjob/utils"

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

	// gin.SetMode(gin.DebugMode)
	// gin.DefaultWriter = io.MultiWriter(os.Stdout)

	r := gin.New()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 自定义日志格式
		return fmt.Sprintf("[%s] - %s \"%s %s %s %d %s \"%s\" %s\"\n",
			param.TimeStamp.Format(time.RFC3339),
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	//加载静态资源文件路径
	r.Static(utils.URI("")+"/static", "./static")

	// r.LoadHTMLFiles("templates/index.html")
	// r.GET("/detail", func(c *gin.Context) {
	// 	c.HTML(200, "index.html", "<a href='lizhouwen.com'>1232</a>")
	// })

	routers.InitRouter(r)

	fmt.Println("The service is running...")

	// port := viper.GetString("server.port")
	endless.ListenAndServe(":8001", r)

	// beego.Run()

}
