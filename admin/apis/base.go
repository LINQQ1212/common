package apis

import (
	"github.com/LINQQ1212/common/admin/models"
	"github.com/LINQQ1212/common/global"
	"github.com/LINQQ1212/common/middleware/jwt_server"
	"github.com/LINQQ1212/common/response"
	"github.com/gin-gonic/gin"
)

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
		}, c)
		return
	}
	response.FailWithMessage("用户名不存在或者密码错误", c)
}
