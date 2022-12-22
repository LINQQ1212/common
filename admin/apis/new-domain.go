package apis

import (
	"encoding/json"
	"github.com/LINQQ1212/common/global"
	"github.com/LINQQ1212/common/models"
	"github.com/LINQQ1212/common/response"
	"github.com/gin-gonic/gin"
	"go.etcd.io/bbolt"
	"strconv"
	"strings"
)

func NewDomain(c *gin.Context) {
	var req struct {
		Version string `json:"version"`
		Domains string `json:"domains"`
	}
	if err := c.BindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	version, ok := global.Versions.Get(req.Version)
	if !ok {
		response.FailWithMessage("版本不存在", c)
		return
	}
	arr := strings.Split(req.Domains, "\n")
	datas := map[string]string{}
	for _, s := range arr {
		s = strings.TrimSpace(s)
		arr2 := strings.SplitN(s, "\t", 2)
		if len(arr2) < 2 {
			continue
		}
		datas[strings.TrimSpace(arr2[0])] = strings.TrimSpace(arr2[1])
	}
	count := 0
	for i, d := range version.Domains {
		if v, ok := datas[d]; ok {
			count++
			version.Domains[i] = v
		}
	}
	b, _ := json.Marshal(&version.Domains)
	err := version.DB.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(models.BDomain).Put(models.BDomain, b)
	})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("修改成功："+strconv.Itoa(count), c)
}
