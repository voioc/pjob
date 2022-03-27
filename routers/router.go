package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/voioc/pjob/controllers"
	"github.com/voioc/pjob/middleware"
)

func InitRouter(engine *gin.Engine) {
	// 默认登录
	// beego.Router("/", &controllers.LoginController{}, "*:Login")
	// beego.Router("/login_in", &controllers.LoginController{}, "*:LoginIn")
	// beego.Router("/login_out", &controllers.LoginController{}, "*:LoginOut")
	// beego.Router("/help", &controllers.HomeController{}, "*:Help")
	// beego.Router("/home", &controllers.HomeController{}, "*:Index")
	// beego.Router("/home/start", &controllers.HomeController{}, "*:Start")

	// beego.AutoRouter(&controllers.TaskController{})
	// beego.AutoRouter(&controllers.GroupController{})
	// beego.AutoRouter(&controllers.TaskLogController{})
	// //资源分组管理
	// beego.AutoRouter(&controllers.ServerGroupController{})
	// beego.AutoRouter(&controllers.ServerController{})
	// beego.AutoRouter(&controllers.BanController{})

	// //权限用户相关
	// beego.AutoRouter(&controllers.AuthController{})
	// beego.AutoRouter(&controllers.RoleController{})
	// beego.AutoRouter(&controllers.AdminController{})
	// beego.AutoRouter(&controllers.UserController{})

	// beego.AutoRouter(&controllers.NotifyTplController{})

	engine.LoadHTMLFiles(
		"views/login/login.html",
		"views/public/main.html",
		"views/home/start.html",
		"views/task/list.html",
		"views/task/detail.html",
		"views/task/add.html",
		"views/task/auditlist.html",
		"views/task/edit.html",
		"views/task/copy.html",
		"views/tasklog/list.html",
		"views/tasklog/detail.html",
	)

	loginC := controllers.LoginController{}
	engine.GET("/", loginC.Login)
	engine.POST("/login_in", loginC.LoginIn)

	r := engine.Group("").Use(middleware.Menu())
	{
		// self.Data["pageTitle"] = "系统首页"
		// r.GET("/home", func(c *gin.Context) {
		// 	fmt.Println(c.Get("menu"))
		// 	menu, _ := service.Menu(1)
		// 	c.HTML(http.StatusOK, "main.html", template.FuncMap{
		// 		"siteName":  "系统首页",
		// 		"SideMenu1": menu["SideMenu1"],
		// 		"SideMenu2": menu["SideMenu2"],
		// 	})
		// })

		homeC := controllers.HomeController{}
		r.GET("/home", homeC.Index)
		r.GET("/home/start", homeC.Start)

		taskC := controllers.TaskController{}
		r.GET("/task/list", taskC.List)
		r.GET("/task/auditlist", taskC.AuditList)
		r.GET("/task/add", taskC.Add)
		r.GET("/task/edit", taskC.Edit)
		r.GET("/task/copy", taskC.Copy)
		r.POST("/task/save", taskC.Save)
		r.GET("/task/detail", taskC.Detail)

		r.GET("/task/table", taskC.Table)

		r.GET("/task/ajaxaudit", taskC.AjaxAudit)
		r.GET("/task/ajaxnopass", taskC.AjaxBatchNoPass)
		r.GET("/task/ajaxstart", taskC.AjaxStart)
		r.GET("/task/ajaxpause", taskC.AjaxPause)
		r.GET("/task/ajaxrun", taskC.AjaxRun)
		r.GET("/task/ajaxdel", taskC.AjaxDel)
		r.POST("/task/notify/type", taskC.AjaxNotifyType)

		r.GET("/task/ajaxbatchstart", taskC.AjaxBatchStart)
		r.GET("/task/ajaxbatchpause", taskC.AjaxBatchPause)
		r.GET("/task/ajaxbatchdel", taskC.AjaxBatchDel)
		r.GET("/task/ajaxbatchaudit", taskC.AjaxAudit)
		r.GET("/task/ajaxbatchnopass", taskC.AjaxBatchNoPass)

		r.POST("/task/apitask", taskC.ApiTask)
		r.GET("/task/apistart", taskC.ApiStart)
		r.GET("/task/apipause", taskC.ApiPause)

		taskLogC := controllers.TaskLogController{}
		r.GET("/tasklog/list", taskLogC.List)
		r.GET("/tasklog/table", taskLogC.Table)
		r.GET("/tasklog/detail", taskLogC.Detail)
		r.GET("/tasklog/ajaxdel", taskLogC.AjaxDel)

	}

}
