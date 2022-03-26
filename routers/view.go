package routers

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/voioc/pjob/controllers"
)

func ViewRouter(engine *gin.Engine) {
	// Simple group: v1
	// engine.LoadHTMLGlob("views/login/login.html")
	engine.LoadHTMLFiles(
		"views/login/login.html",
		"views/public/main.html",
		"views/home/start.html",
		"views/task/list.html",
	)
	view := engine.Group("")
	{
		view.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", template.FuncMap{
				"url": "login_in",
			})
		})

		// self.Data["pageTitle"] = "系统首页"
		view.GET("/home", func(c *gin.Context) {
			c.HTML(http.StatusOK, "main.html", template.FuncMap{
				"siteName": "系统首页",
			})
		})

		homeC := controllers.HomeController{}
		view.GET("/home/start", homeC.Start)

		view.GET("/home/list", func(c *gin.Context) {
			c.HTML(http.StatusOK, "list.html", template.FuncMap{
				"siteName": "系统首页",
			})
		})
	}
}
