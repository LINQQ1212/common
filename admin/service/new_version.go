package service

import (
	"github.com/LINQQ1212/common/create_db"
	"github.com/LINQQ1212/common/global"
	"github.com/LINQQ1212/common/models"
)

func NewVersion(req models.NewVersionReq) error {
	cdb := create_db.New(req)
	err := cdb.Start()
	if err != nil {
		return err
	}
	v, err := models.NewVersion(global.VersionDir, cdb.GetDomain())
	if err != nil {
		return err
	}
	global.Versions.Set(cdb.GetDomain(), v)
	return global.VHost.NewVersion(cdb.GetDomain())
}
