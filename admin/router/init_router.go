package router

import (
	"github.com/LINQQ1212/common/admin/apis"
	"github.com/LINQQ1212/common/middleware"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitRouter(r *gin.Engine) {
	r.GET("/metrics", gin.BasicAuth(gin.Accounts{"fgadmin": "DHBOXlrZc71fOR8i"}), ginprom.PromHandler(promhttp.Handler()))

	a := r.Group("/h6hb7860q2")
	{
		a.Static("/static", "admin_static/static")
		a.StaticFile("/", "admin_static/index.html")
		//a.Static("/", "admin")
	}
	r.NoRoute(func(c *gin.Context) {
		c.File("admin_static/index.html")
	})

	r.POST("/h6hb7860q2/api/login", apis.Login)
	admin := r.Group("/h6hb7860q2/api", middleware.JWTAuth())
	{

		admin.GET("menus", apis.Menus)
		admin.POST("/new/version", apis.NewVersion)
		admin.POST("/new/domain", apis.NewDomain)

	}
}
