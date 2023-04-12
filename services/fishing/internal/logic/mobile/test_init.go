package mobile

import (
	"fmt"
	"github.com/realmicro/realmicro"
	mconfig "github.com/realmicro/realmicro/config"
	"github.com/realmicro/realmicro/config/encoder/yaml"
	"github.com/realmicro/realmicro/config/reader"
	"github.com/realmicro/realmicro/config/reader/json"
	"github.com/realmicro/realmicro/config/source/file"
	"github.com/realmicro/realmicro/registry"
	"github.com/realmicro/realmicro/registry/etcd"
	"github.com/realmicro/realmicro/wrapper/select/dc"
	"qingyun/common/wrapper"
	"qingyun/services/fishing/internal/config"
	"qingyun/services/fishing/internal/svc"
)

var conf config.Config
var mobileFishingService *MobileFishingService
var configTestFile = "../../../etc/config.yaml"

func init() {
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
		file.WithPath(configTestFile),
	)); err != nil {
		fmt.Println(err)
		return
	}
	if err = c.Scan(&conf); err != nil {
		fmt.Println(err)
		return
	}
	service := realmicro.NewService(
		realmicro.Name(conf.ServiceName),
		realmicro.Version(conf.Version),
		realmicro.Metadata(map[string]string{
			"env":     conf.Env,
			"project": conf.Project,
		}),
		realmicro.Registry(etcd.NewRegistry(registry.Addrs(conf.Hosts.Etcd.Address...))),
		realmicro.WrapHandler(wrapper.LogHandler()),
		realmicro.WrapClient(dc.NewDCWrapper, wrapper.LogCall),
	)

	service.Init()
	conf.Hosts.Mysql.Host = "101.43.58.200:3306"
	conf.Hosts.Mysql.User = "fishing"
	conf.Hosts.Mysql.Pass = "Emt3rfwCdtGack5D"
	ctx := svc.NewServiceContext(&conf, service)
	mobileFishingService = &MobileFishingService{
		SvcContext: ctx,
	}
}
