package initialize

import (
	"github.com/LINQQ1212/common/global"
	"github.com/robfig/cron/v3"
	"time"
)

func NewCronin() *cron.Cron {
	c := cron.New()
	c.AddFunc("@daily", setDaily)
	c.AddFunc("@hourly", setHourly)
	c.Start()
	setDaily()
	setHourly()
	return c
}

func setDaily() {
	//fmt.Println("@daily")
	year, month, day := time.Now().Date()
	global.Day = time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func setHourly() {
	//fmt.Println("@daily")
	year, month, day := time.Now().Date()
	global.Hour = time.Date(year, month, day, time.Now().Hour(), 0, 0, 0, time.Local)
}
