package apis

import (
	"github.com/LINQQ1212/common/admin/service"
	"github.com/LINQQ1212/common/global"
	"github.com/LINQQ1212/common/models"
	"github.com/LINQQ1212/common/response"
	"github.com/LINQQ1212/common/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
	"path"
	"strings"
)

func NewVersion(c *gin.Context) {
	req := models.NewVersionReq{}
	if err := c.BindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 下载 产品数据
	if strings.HasPrefix(req.ProductTarLink, "http") {
		fileP := path.Join(os.TempDir(), path.Base(req.ProductTarLink))
		err := utils.DownFile(req.ProductTarLink, fileP)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		req.ProductTarLink = fileP
	}
	// 下载google 图片
	if strings.HasPrefix(req.GoogleImg, "http") {
		fileP := path.Join(os.TempDir(), path.Base(req.GoogleImg))
		err := utils.DownFile(req.GoogleImg, fileP)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		req.GoogleImg = fileP
	}
	go func() {
		err := service.NewVersion(req)
		if err != nil {
			global.LOG.Error("NewVersion", zap.String("version", req.Domain), zap.Error(err))
		}
	}()
	response.OkWithMessage("后台执行中", c)
}
