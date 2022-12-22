package apis

import (
	"github.com/LINQQ1212/common/global"
	"github.com/LINQQ1212/common/models"
	"github.com/LINQQ1212/common/response"
	"github.com/gin-gonic/gin"
)

type menu struct {
	Id    int    `json:"id"`
	Key   string `json:"key"`
	Title string `json:"title"`
	Path  string `json:"path"`
	Order int    `json:"order"`
}

func Menus(c *gin.Context) {
	l := []menu{
		{
			Id:    1,
			Key:   "new-version",
			Title: "新版本",
			Path:  "/new-version",
			Order: 1,
		},
		{
			Id:    2,
			Key:   "new-domain",
			Title: "修改域名",
			Path:  "/new-domain",
			Order: 800,
		},
	}
	i := 3
	global.Versions.Range(func(k string, v *models.Version) bool {
		l = append(l, menu{
			Id:    i,
			Key:   v.Info.Name,
			Title: v.Info.Name,
			Path:  "/",
			Order: i + 3,
		})
		i++
		return true
	})

	response.OkWithData(l, c)
}