package svc

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"qingyun/common/store/mysql"
	"qingyun/services/wechat/internal/config"
	"qingyun/services/wechat/internal/wechat"
	"qingyun/services/wechat/models"
	"xorm.io/xorm"
)

type ServiceContext struct {
	Config  *config.Config
	HttpApp *gin.Engine
}

func NewServiceContext(config *config.Config) *ServiceContext {
	mysql.Init(
		mysql.WithHost(config.Hosts.Mysql.Host),
		mysql.WithUser(config.Hosts.Mysql.User),
		mysql.WithPass(config.Hosts.Mysql.Pass),
		mysql.WithDBName(config.Hosts.Mysql.DBName),
		mysql.IfShowSql(config.Hosts.Mysql.IfShowSql),
		mysql.IfSyncDB(config.Hosts.Mysql.IfSyncDB),
		mysql.AfterInit(func(x *xorm.Engine) {
			if config.Hosts.Mysql.IfSyncDB {
				// sync tables
				if err := x.Sync2(
					new(models.WechatBot),
				); err != nil {
					fmt.Printf("db sync error: %v\n", err)
					os.Exit(1)
				}
			}
		}),
	)
	wechat.InitWechatBotsMap()
	wechat.InitBotWithStart()
	return &ServiceContext{
		Config:  config,
		HttpApp: gin.Default(),
	}
}
