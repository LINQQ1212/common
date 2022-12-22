package initialize

import (
	"github.com/LINQQ1212/common/global"
	"github.com/LINQQ1212/common/models"
	"github.com/LINQQ1212/common/utils"
	"os"
	"path"
	"strings"
)

func InitVersions() {
	global.VersionDir = path.Join(global.CONFIG.System.MainDir, "version_data")
	if b, _ := utils.PathExists(global.VersionDir); !b {
		os.RemoveAll(global.VersionDir)
		err := os.Mkdir(global.VersionDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
		return
	}

	fs, err := os.ReadDir(global.VersionDir)
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		if strings.HasPrefix(f.Name(), "~") || f.IsDir() || !strings.HasSuffix(f.Name(), ".db") {
			continue
		}
		name := strings.TrimSuffix(f.Name(), ".db")
		v, err := models.NewVersion(global.VersionDir, name)
		if err != nil {
			panic(err)
		}
		global.Versions.Set(name, v)
	}

}
