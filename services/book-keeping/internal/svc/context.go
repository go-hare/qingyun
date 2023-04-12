package svc

import (
	"github.com/realmicro/realmicro"
	"github.com/realmicro/realmicro/client"
	"github.com/realmicro/realmicro/server"
	"qingyun/services/book-keeping/internal/config"
)

type ServiceContext struct {
	Config *config.Config
	Server server.Server
	Client client.Client
}

func NewServiceContext(config *config.Config, service realmicro.Service) *ServiceContext {
	return &ServiceContext{
		Config: config,
		Server: service.Server(),
		Client: service.Client(),
	}
}
