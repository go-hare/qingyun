package svc

import (
	"fmt"
	"github.com/realmicro/realmicro"
	"github.com/realmicro/realmicro/client"
	"github.com/realmicro/realmicro/server"
	"os"
	"qingyun/common/store/mysql"
	"qingyun/services/realmicro_web/internal/config"
	"qingyun/services/realmicro_web/internal/http"
	"qingyun/services/realmicro_web/models"
	"xorm.io/xorm"
)

type ServiceContext struct {
	Config *config.Config
	Server server.Server
	Client client.Client
}

func NewServiceContext(config *config.Config, service realmicro.Service) *ServiceContext {
	fmt.Println(config)
	mysql.Init(
		mysql.WithHost(config.DBInfo.Host),
		mysql.WithUser(config.DBInfo.User),
		mysql.WithPass(config.DBInfo.Pass),
		mysql.WithDBName(config.DBInfo.DBName),
		mysql.IfShowSql(config.DBInfo.IfShowSql),
		mysql.IfSyncDB(config.DBInfo.IfSyncDB),
		mysql.AfterInit(func(x *xorm.Engine) {
			if config.DBInfo.IfSyncDB {
				// sync tables
				if err := x.Sync2(
					new(models.App),
					new(models.AdminUser),
					new(models.Env),
					new(models.Project),
					new(models.Department),
				); err != nil {
					fmt.Printf("db sync error: %v\n", err)
					os.Exit(1)
				}
			}
		}),
	)
	http.InitEnvs()
	return &ServiceContext{
		Config: config,
		Server: service.Server(),
		Client: service.Client(),
	}
}
