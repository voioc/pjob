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

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	_ "github.com/voioc/cjob/common"
	"github.com/voioc/cjob/routers"
	"github.com/voioc/coco/cache"
	"github.com/voioc/coco/logzap"
	// _ "github.com/voioc/cjob/routers"
)

func main() {

	// db.Init()        // 初始化数据库连接
	cache.Init()     // 初始化缓存连接
	logzap.InitZap() // 初始化日志组件

	r := gin.New()
	// pprof.Register(r)
	// r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
	// 	// 自定义日志格式
	// 	return fmt.Sprintf("[%s] - %s \"%s %s %s %d %s \"%s\" %s\"\n",
	// 		param.TimeStamp.Format(time.RFC3339),
	// 		param.ClientIP,
	// 		param.Method,
	// 		param.Path,
	// 		param.Request.Proto,
	// 		param.StatusCode,
	// 		param.Latency,
	// 		param.Request.UserAgent(),
	// 		param.ErrorMessage,
	// 	)
	// }))

	// if err := tool.SendSMS(); err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println("短信发送成功")
	// }

	// 中间件
	// r.Use(middleware.Trace(), gin.Recovery(), middleware.CORS(), middleware.ZapLogger())

	// 主页
	// r.GET("/hello", func(c *gin.Context) {
	// 	c.String(200, "Hello, Melon")
	// })

	//  beego.Router("/", &controllers.LoginController{}, "*:Login")
	// loginController := controllers.LoginController{}
	//加载相对路径
	// filepath.Abs(filepath.Dir(os.Args[0]))

	//加载静态资源文件路径
	r.Static("/static", "./static")

	//router.LoadHTMLFiles("templates/index.tmpl")
	//router.LoadHTMLFiles("templates/index.tmpl", "templates/goods.hmpl"
	// r.LoadHTMLGlob("views/**/*")
	//多层目录，多文件重名，在html文件中声明即可 {{defind "goods/list.html"}} {{end}}

	// r.SetFuncMap(template.FuncMap{
	// 	"abc": abc,
	// })

	// r.LoadHTMLFiles("templates/index.html")
	// r.GET("/detail", func(c *gin.Context) {
	// 	c.HTML(200, "index.html", "<a href='lizhouwen.com'>1232</a>")
	// })

	// r.LoadHTMLGlob("views/login/login.html")
	// r.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "login.html", template.FuncMap{
	// 		"urlfor": "dddd222dddd",
	// 	})
	// })

	routers.InitRouter(r)

	fmt.Println("The service is running...")
	endless.ListenAndServe(":8001", r)

	// beego.Run()

}
