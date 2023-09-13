package main

import (
	"bbs-go/controllers"
	"bbs-go/model"
	"bbs-go/pkg/common"
	"bbs-go/pkg/config"
	"bbs-go/pkg/iplocator"
	"bbs-go/scheduler"
	_ "bbs-go/services/eventhandler"
	"flag"
	"io"
	"os"
	"time"

	"github.com/mlogclub/simple/sqls"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var configFile = flag.String("config", "./bbs-go.yaml", "配置文件路径")

func init() {
	flag.Parse()

	// 初始化配置
	conf := config.Init(*configFile)

	// 初始化日志
	if file, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
		logrus.SetOutput(io.MultiWriter(os.Stdout, file))
	} else {
		logrus.SetOutput(os.Stdout)
		logrus.Error(err)
	}

	// 连接数据库
	gormConf := &gorm.Config{
		Logger: logger.New(logrus.StandardLogger(), logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
		}),
	}
	if err := sqls.Open(conf.DB, gormConf, model.Models...); err != nil {
		logrus.Error(err)
	}

	// ip定位
	iplocator.InitIpLocator(conf.IpDataPath)
}

func main() {
	if common.IsProd() {
		// 开启定时任务
		scheduler.Start()
	}
	controllers.Router()
}
