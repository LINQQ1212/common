package service

import (
	"github.com/LINQQ1212/common/create_db"
	"github.com/LINQQ1212/common/global"
	"github.com/LINQQ1212/common/models"
	"go.uber.org/zap"
)

func NewVersion(req models.NewVersionReq) error {
	cdb := create_db.New(req)
	err := cdb.Start()
	if err != nil {
		global.LOG.Error("Start", zap.Error(err))
		return err
	}
	v, err := models.NewVersion(global.VersionDir, cdb.GetDomain())
	if err != nil {
		global.LOG.Error("NewVersion", zap.Error(err))
		return err
	}
	global.Versions.Set(cdb.GetDomain(), v)
	return global.VHost.NewVersion(cdb.GetDomain())
}
