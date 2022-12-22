package global

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/LINQQ1212/common/config"
	"github.com/LINQQ1212/common/models"
	"github.com/LINQQ1212/common/vhost"
	"github.com/LINQQ1212/common/words"
	"github.com/cornelk/hashmap"
	"github.com/robfig/cron/v3"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"github.com/tdewolff/minify/v2"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"time"
)

var (
	// CONFIG 配置
	CONFIG config.Server
	// VP 配置处理
	VP *viper.Viper
	// LOG 日志
	LOG *zap.Logger

	// BlackCache 本地缓存
	BlackCache local_cache.Cache

	// ConcurrencyControl 一个重复的函数调用抑制机制
	ConcurrencyControl = &singleflight.Group{}

	Cron *cron.Cron

	// IMGDB 图片存储
	//IMGDB *badger.DB

	VersionDir string
	//Versions *hashmap.List[string, *models.Version]
	Versions = hashmap.New[string, *models.Version]()

	// VHost  站点信息
	VHost *vhost.V

	// View 模板引擎
	View *jet.Set

	// Words 促销词
	Words *words.Words
	// EnWords 英语单词
	EnWords *words.Words
	// Names  人名
	Names *words.Words
	//Review *review.Review

	Minify *minify.M

	Day  time.Time // 日期
	Hour time.Time // 日期

)
