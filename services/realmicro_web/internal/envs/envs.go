package envs

import (
	"fmt"
	"github.com/realmicro/realmicro/client"
	"github.com/realmicro/realmicro/registry"
	"github.com/realmicro/realmicro/registry/etcd"
	log "github.com/sirupsen/logrus"
	"qingyun/services/realmicro_web/internal/config"
	"qingyun/services/realmicro_web/models"
	"sort"
	"strings"
	"sync"
)

type envList []*Env

func (e envList) Len() int {
	return len(e)
}
func (e envList) Less(i, j int) bool {
	return e[i].Env < e[j].Env || (e[i].Env == e[j].Env && e[i].Name < e[j].Name)
}
func (e envList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

type EnvDetail struct {
	Env  string
	Envs []*Env
}

type Env struct {
	Env        string
	Name       string
	Cluster    string
	IfDefault  bool
	IfSelected bool
	Addresses  []string
	register   registry.Registry
	config     *config.EtcdConfig
	client     client.Client
}

func (e *Env) Build() string {
	return e.Env + "/" + e.Name
}

func BuildEnv(env, name string) string {
	return env + "/" + name
}

type PaaSEnvs struct {
	envMutex sync.RWMutex
	envs     map[string]*Env
}

func NewPaaSEnvs() *PaaSEnvs {
	pe := &PaaSEnvs{
		envs: make(map[string]*Env),
	}
	//pe.AddEnv(&Env{
	//	Env:       "TEST",
	//	Name:      "Default",
	//	Cluster:   "单机",
	//	IfDefault: true,
	//	Addresses: []string{"127.0.0.1:2379"},
	//})
	//pe.AddEnv(&Env{
	//	Env:       "PROD",
	//	Name:      "Default",
	//	Cluster:   "集群",
	//	IfDefault: false,
	//	Addresses: []string{"127.0.0.1:2379"},
	//})
	el, err := models.GetEnvList(nil)
	if err != nil {
		log.Info("get env list error: %v", err)
	} else {
		for i := 0; i < len(el); i++ {
			pe.AddEnv(&Env{
				Env:       el[i].Env,
				Name:      el[i].Name,
				Cluster:   el[i].Cluster,
				IfDefault: el[i].IfDefault == 1,
				Addresses: strings.Split(el[i].Addresses, ","),
			})
		}
	}
	return pe
}

func (pe *PaaSEnvs) GetEnvs() ([]*EnvDetail, *Env) {
	pe.envMutex.RLock()
	defer pe.envMutex.RUnlock()

	var eList []*Env
	for _, v := range pe.envs {
		eList = append(eList, v)
	}
	sort.Sort(envList(eList))

	var nowEnv string
	var list []*EnvDetail
	defaultEnv := new(Env)
	for i := 0; i < len(eList); i++ {
		if eList[i].Env != nowEnv {
			list = append(list, &EnvDetail{
				Env: eList[i].Env,
				Envs: []*Env{&Env{
					Env:  eList[i].Env,
					Name: eList[i].Name,
				}},
			})
		} else {
			list[len(list)-1].Envs = append(list[len(list)-1].Envs, &Env{
				Env:  eList[i].Env,
				Name: eList[i].Name,
			})
		}
		if eList[i].IfDefault {
			defaultEnv = eList[i]
		}
	}

	return list, defaultEnv
}

func (pe *PaaSEnvs) AddEnv(env *Env) {
	pe.envMutex.Lock()
	env.register = etcd.NewRegistry(registry.Addrs(env.Addresses...))
	env.config = config.NewEtcdConfig(config.WithAddress(env.Addresses))
	env.client = client.NewClient(client.Registry(env.register))
	pe.envs[env.Build()] = env
	pe.envMutex.Unlock()
}

func (pe *PaaSEnvs) GetClient(env string) client.Client {
	pe.envMutex.RLock()
	defer pe.envMutex.RUnlock()
	e, ok := pe.envs[env]
	if !ok {
		return nil
	}
	return e.client
}

func (pe *PaaSEnvs) ListServices(env string) ([]*registry.Service, error) {
	pe.envMutex.RLock()
	e, ok := pe.envs[env]
	pe.envMutex.RUnlock()
	if !ok {
		return nil, fmt.Errorf("cannot found this env: %s", env)
	}
	return e.register.ListServices()
}

func (pe *PaaSEnvs) GetService(env, svc string) ([]*registry.Service, error) {
	pe.envMutex.RLock()
	e, ok := pe.envs[env]
	pe.envMutex.RUnlock()
	if !ok {
		return nil, fmt.Errorf("cannot found this env: %s", env)
	}
	return e.register.GetService(svc)
}

func (pe *PaaSEnvs) GetSvcConfig(env, svc string) (map[string]*config.Node, error) {
	pe.envMutex.RLock()
	e, ok := pe.envs[env]
	pe.envMutex.RUnlock()
	if !ok {
		return nil, fmt.Errorf("cannot found this env: %s", env)
	}
	return e.config.Get(svc)
}

func (pe *PaaSEnvs) ConfigDel(env, key string) error {
	pe.envMutex.RLock()
	e, ok := pe.envs[env]
	pe.envMutex.RUnlock()
	if !ok {
		return fmt.Errorf("cannot found this env: %s", env)
	}
	return e.config.Del(key)
}

func (pe *PaaSEnvs) ConfigPut(env, key, value string) error {
	pe.envMutex.RLock()
	e, ok := pe.envs[env]
	pe.envMutex.RUnlock()
	if !ok {
		return fmt.Errorf("cannot found this env: %s", env)
	}
	return e.config.Put(key, value)
}
