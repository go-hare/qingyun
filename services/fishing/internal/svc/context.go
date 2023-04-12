package svc

import (
	"fmt"
	"github.com/realmicro/realmicro"
	"github.com/realmicro/realmicro/client"
	"github.com/realmicro/realmicro/server"
	"os"
	"qingyun/common/store/mysql"
	"qingyun/services/fishing/internal/config"
	"qingyun/services/fishing/models"
	"xorm.io/xorm"
)

type ServiceContext struct {
	Config *config.Config
	Server server.Server
	Client client.Client
}

func NewServiceContext(config *config.Config, service realmicro.Service) *ServiceContext {
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
					new(models.User),
					new(models.UserPlatform),
					new(models.UserWechat),
					new(models.Banner),
					new(models.UserFollow),
					new(models.Infomation),
					new(models.UserLike),
					new(models.UserCollection),
					new(models.Buy),
					new(models.Article),
					new(models.ArticleTag),
					new(models.ArticleCategory),
					new(models.ArticleTheme),
					new(models.ArticleTagDetail),
					new(models.Comment),
					new(models.UserCommentLike),
					new(models.Message),
					new(models.UserExtend),
					new(models.InfomationOrder),
					new(models.FishingOrder),
					new(models.UserCoinPayBill),
				); err != nil {
					fmt.Printf("db sync error: %v\n", err)
					os.Exit(1)
				}
			}
		}),
	)
	return &ServiceContext{
		Config: config,
		Server: service.Server(),
		Client: service.Client(),
	}
}
