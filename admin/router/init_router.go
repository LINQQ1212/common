package router

import (
	"github.com/LINQQ1212/common/admin/apis"
	"github.com/LINQQ1212/common/middleware"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func InitRouter(r *gin.Engine) {
	r.GET("/version.txt", apis.VersionPHP)
	r.GET("/version-jm.txt", apis.VersionPHP)
	r.GET("/version", apis.VersionPHP)
	r.GET("/versionjm", apis.VersionPHP)
	r.GET("/metrics", gin.BasicAuth(gin.Accounts{"fgadmin": "DHBOXlrZc71fOR8i"}), ginprom.PromHandler(promhttp.Handler()))

	a := r.Group("/h6hb7860q2")
	{
		a.GET("/", apis.Admin)
		a.GET("/new-version", apis.Admin)
		a.GET("/new-domain", apis.Admin)
		a.GET("/login", apis.Admin)
	}
	r.NoRoute(func(c *gin.Context) {
		c.Writer.WriteString("404")
		c.Status(http.StatusNotFound)
		c.Abort()
	})

	r.POST("/h6hb7860q2/api/login", apis.Login)
	admin := r.Group("/h6hb7860q2/api", middleware.JWTAuth())
	{
		admin.GET("menus", apis.Menus)
		admin.POST("/new/version", apis.NewVersion)
		admin.POST("/new/domain", apis.NewDomain)
	}
}
