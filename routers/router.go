package routers

import (
	"html/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/voioc/cjob/app/handler"
	"github.com/voioc/cjob/app/model"
	"github.com/voioc/cjob/middleware"
	"github.com/voioc/cjob/utils"
)

func InitRouter(engine *gin.Engine) {
	model.Init(time.Now().Unix())

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

	engine.SetFuncMap(template.FuncMap{
		"str2html": func(str string) template.HTML {
			return template.HTML(str)
		},
	})

	// engine.StaticFile("/favicon.ico", "./static/admin/images/favicon.ico")
	engine.LoadHTMLGlob("views/**/*")

	// engine.LoadHTMLFiles(
	// 	"views/public/help.html",

	// 	"views/login/login.html",
	// 	"views/public/main.html",
	// 	"views/home/start.html",
	// 	"views/task/list.html",
	// 	"views/task/detail.html",
	// 	"views/task/add.html",
	// 	"views/task/auditlist.html",
	// 	"views/task/edit.html",
	// 	"views/task/copy.html",
	// 	"views/tasklog/list.html",
	// 	"views/tasklog/detail.html",
	// 	"views/group/list.html",
	// 	"views/group/add.html",
	// 	"views/group/edit.html",
	// 	"views/server/add.html",
	// 	"views/server/copy.html",
	// 	"views/server/edit.html",
	// 	"views/server/list.html",

	// 	"views/servergroup/add.html",
	// 	"views/servergroup/edit.html",
	// 	"views/servergroup/list.html",

	// 	"views/ban/add.html",
	// 	"views/ban/edit.html",
	// 	"views/ban/list.html",

	// 	"views/notify/add.html",
	// 	"views/notify/edit.html",
	// 	"views/notify/list.html",

	// 	"views/admin/add.html",
	// 	"views/admin/edit.html",
	// 	"views/admin/list.html",

	// 	"views/role/add.html",
	// 	"views/role/edit.html",
	// 	"views/role/list.html",

	// 	"views/auth/list.html",
	// 	"views/user/edit.html",
	// )

	loginC := handler.LoginController{}
	engine.GET(utils.URI("")+"/", loginC.Login)
	engine.GET(utils.URI("")+"/login", loginC.Login)
	engine.POST(utils.URI("")+"/login_in", loginC.LoginIn)
	engine.GET(utils.URI("")+"/login_out", loginC.LoginOut)

	r := engine.Group(utils.URI("")).Use(middleware.Auth())
	{
		homeC := handler.HomeController{}
		r.GET("/home", homeC.Index)
		r.GET("/home/start", homeC.Start)
		r.GET("/help", homeC.Help)

		taskC := handler.TaskController{}
		r.GET("/task/list", taskC.List)
		r.GET("/task/auditlist", taskC.AuditList)
		r.GET("/task/add", taskC.Add)
		r.GET("/task/edit", taskC.Edit)
		r.GET("/task/copy", taskC.Copy)
		r.POST("/task/save", taskC.Save)
		r.GET("/task/detail", taskC.Detail)

		r.GET("/task/table", taskC.Table)

		r.POST("/task/ajaxaudit", taskC.Audit)
		// r.POST("/task/ajaxnopass", taskC.AjaxBatchNoPass)
		r.POST("/task/ajaxstart", taskC.AjaxStart)
		r.POST("/task/ajaxpause", taskC.AjaxPause)
		r.POST("/task/ajaxrun", taskC.AjaxRun)
		r.POST("/task/ajaxdel", taskC.AjaxDel)
		r.POST("/task/notify/type", taskC.AjaxNotifyType)

		r.POST("/task/ajaxbatchstart", taskC.AjaxBatchStart)
		r.POST("/task/ajaxbatchpause", taskC.AjaxBatchPause)
		r.POST("/task/ajaxbatchdel", taskC.AjaxBatchDel)
		r.POST("/task/audit", taskC.AjaxBatchAudit)
		r.POST("/task/reject", taskC.Reject)

		r.POST("/task/apitask", taskC.ApiTask)
		r.GET("/task/apistart", taskC.ApiStart)
		r.GET("/task/apipause", taskC.ApiPause)

		taskLogC := handler.TaskLogController{}
		r.GET("/tasklog/list", taskLogC.List)
		r.GET("/tasklog/table", taskLogC.Table)
		r.GET("/tasklog/detail", taskLogC.Detail)
		r.GET("/tasklog/ajaxdel", taskLogC.AjaxDel)

		groupC := handler.GroupController{}
		r.GET("/group/list", groupC.List)
		r.GET("/group/table", groupC.Table)
		r.GET("/group/add", groupC.Add)
		r.GET("/group/edit", groupC.Edit)
		r.POST("/group/save", groupC.AjaxSave)
		r.POST("/group/del", groupC.AjaxDel)

		serverC := handler.ServerController{}
		r.GET("/server/list", serverC.List)
		r.GET("/server/table", serverC.Table)
		r.GET("/server/add", serverC.Add)
		r.GET("/server/edit", serverC.Edit)
		r.POST("/server/save", serverC.AjaxSave)
		r.POST("/server/del", serverC.AjaxDel)
		r.POST("/server/test", serverC.AjaxTestServer)

		serverGroupC := handler.ServerGroupController{}
		r.GET("/servergroup/list", serverGroupC.List)
		r.GET("/server/group/table", serverGroupC.Table)
		r.GET("/server/group/add", serverGroupC.Add)
		r.GET("/server/group/edit", serverGroupC.Edit)
		r.POST("/server/group/save", serverGroupC.AjaxSave)
		r.POST("/server/group/del", serverGroupC.AjaxDel)

		banC := handler.BanController{}
		r.GET("/ban/list", banC.List)
		r.GET("/ban/table", banC.Table)
		r.GET("/ban/add", banC.Add)
		r.GET("/ban/edit", banC.Edit)
		r.POST("/ban/save", banC.AjaxSave)
		r.POST("/ban/del", banC.AjaxDel)

		notifyC := handler.NotifyController{}
		r.GET("/notifytpl/list", notifyC.List)
		r.GET("/notify/table", notifyC.Table)
		r.GET("/notify/add", notifyC.Add)
		r.GET("/notify/edit", notifyC.Edit)
		r.POST("/notify/save", notifyC.AjaxSave)
		r.POST("/notify/del", notifyC.AjaxDel)

		adminC := handler.AdminController{}
		r.GET("/admin/list", adminC.List)
		r.GET("/admin/table", adminC.Table)
		r.GET("/admin/add", adminC.Add)
		r.GET("/admin/edit", adminC.Edit)
		r.POST("/admin/save", adminC.AjaxSave)
		r.POST("/admin/del", adminC.AjaxDel)

		roleC := handler.RoleController{}
		r.GET("/role/list", roleC.List)
		r.GET("/role/table", roleC.Table)
		r.GET("/role/add", roleC.Add)
		r.GET("/role/edit", roleC.Edit)
		r.POST("/role/save", roleC.AjaxSave)
		r.POST("/role/del", roleC.AjaxDel)

		authC := handler.AuthController{}
		// r.GET("/auth/index", authC.Index)
		r.GET("/auth/list", authC.List)
		r.GET("/auth/getnodes", authC.GetNodes)
		r.GET("/auth/getnode", authC.GetNode)
		r.POST("/auth/save", authC.AjaxSave)
		r.POST("/auth/del", authC.AjaxDel)

		userC := handler.UserController{}
		// r.GET("/auth/index", authC.Index)
		r.GET("/user/edit", userC.Edit)
		r.POST("/user/save", userC.AjaxSave)
	}

}
