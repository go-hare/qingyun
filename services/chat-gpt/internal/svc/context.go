package svc

import (
	"fmt"
	"github.com/realmicro/realmicro"
	"github.com/realmicro/realmicro/client"
	"github.com/realmicro/realmicro/server"
	"net/http"
	"net/url"
	"os"
	"qingyun/common/openai"
	"qingyun/common/store/mysql"
	"qingyun/services/chat-gpt/internal/config"
	"qingyun/services/chat-gpt/models"
	"time"
	"xorm.io/xorm"
)

type ServiceContext struct {
	Config       *config.Config
	Server       server.Server
	Client       client.Client
	OpenAiClient *openai.Client
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
					new(models.Category),
					new(models.Model),
					new(models.User),
				); err != nil {
					fmt.Printf("db sync error: %v\n", err)
					os.Exit(1)
				}
			}
		}),
	)
	conf := openai.DefaultConfig("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")
	hc := http.Client{Timeout: 30 * time.Second}
	hc.Transport = &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse("http://127.0.0.1:7890")
		},
	}
	conf.HTTPClient = &hc
	openAiClient := openai.NewClientWithConfig(conf)

	return &ServiceContext{
		Config:       config,
		Server:       service.Server(),
		Client:       service.Client(),
		OpenAiClient: openAiClient,
	}
}
