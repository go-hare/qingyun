package main

import (
	"flag"
	"fmt"
	mconfig "github.com/realmicro/realmicro/config"
	"github.com/realmicro/realmicro/config/encoder/yaml"
	"github.com/realmicro/realmicro/config/reader"
	"github.com/realmicro/realmicro/config/reader/json"
	"github.com/realmicro/realmicro/config/source/file"
	"github.com/realmicro/realmicro/logger"
	mlogrus "github.com/realmicro/realmicro/logger/logrus"
	"github.com/sirupsen/logrus"
	"qingyun/services/wechat/internal/config"
	"qingyun/services/wechat/internal/logic"
	"qingyun/services/wechat/internal/middleware"
	"qingyun/services/wechat/internal/svc"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()
	c, _ := mconfig.NewConfig(
		mconfig.WithReader(
			json.NewReader(
				reader.WithEncoder(yaml.NewEncoder()),
			),
		),
	)
	var err error
	// load the config from a file source
	if err = c.Load(file.NewSource(
		file.WithPath(*configFile),
	)); err != nil {
		fmt.Println(err)
		return
	}
	var cfg config.Config
	if err = c.Scan(&cfg); err != nil {
		fmt.Println(err)
		return
	}
	logger.DefaultLogger = mlogrus.NewLogger(mlogrus.WithJSONFormatter(&logrus.JSONFormatter{}))
	ctx := svc.NewServiceContext(&cfg)
	ctx.HttpApp.Use(middleware.CheckAppKeyExistMiddleware(), middleware.CheckAppKeyIsLoggedInMiddleware())
	logic.InitLoginRoute(ctx.HttpApp)
	_ = ctx.HttpApp.Run(":8080")

}
