package config

import (
	"context"
	"go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

var (
	DefaultPrefix = "/real/config/"
	DirValue      = "etcdv3_dir_$2H#%gRe3*t"
)

type DBInfo struct {
	User      string
	Pass      string
	Host      string
	DBName    string
	IfShowSql bool
	IfSyncDB  bool
}

type Etcd struct {
	Address          []string
	RegisterTTL      int64
	RegisterInterval int64
}

type Config struct {
	Debug   bool
	Path    string
	Version string

	Address string

	Project      string
	ServerName   string
	IfProduction bool

	DBInfo DBInfo
	Etcd  Etcd
}

type EtcdConfig struct {
	timeout time.Duration
	prefix  string
	addr    []string
	client  *clientv3.Client
	opts    EtcdOpts
}

func NewEtcdConfig(opts ...Option) *EtcdConfig {
	options := NewOptions(opts...)

	endpoints := []string{"127.0.0.1:2379"}
	if options.address != nil {
		endpoints = options.address
	}
	defaultTimeout := 5 * time.Second
	if options.timeout != 0 {
		defaultTimeout = options.timeout
	}

	// create the client
	client, _ := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: defaultTimeout,
	})

	prefix := DefaultPrefix
	if options.prefix != "" {
		prefix = DefaultPrefix + prefix
	}

	return &EtcdConfig{
		timeout: defaultTimeout,
		addr:    endpoints,
		opts:    options,
		client:  client,
		prefix:  prefix,
	}
}

type Node struct {
	Key     string           `json:"key"`
	LongKey string           `json:"longKey"`
	Value   string           `json:"value,omitempty"`
	IsDir   bool             `json:"dir,omitempty"`
	Nodes   map[string]*Node `json:"nodes,omitempty"`
}

func (ec *EtcdConfig) Get(key string) (map[string]*Node, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ec.timeout)
	defer cancel()

	realKey := ec.prefix + key
	rsp, err := ec.client.Get(ctx, realKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	results := make(map[string]*Node)
	for i := 0; i < len(rsp.Kvs); i++ {
		//log.Logf("key: %s value: %s", string(rsp.Kvs[i].Key), string(rsp.Kvs[i].Value))
		nodeKey := string(rsp.Kvs[i].Key)
		nodeKey = strings.Replace(nodeKey, realKey, "", -1)
		keys := strings.Split(nodeKey, "/")
		//log.Logf("key: %s keys: %v %d", nodeKey, keys, len(keys))
		if len(keys) <= 1 {
			continue
		}
		kvals := results
		for j := 1; j < len(keys); j++ {
			kval, ok := kvals[keys[j]]
			if !ok {
				kval = &Node{Nodes: make(map[string]*Node)}
				kvals[keys[j]] = kval
			}
			if j == len(keys)-1 {
				//kval.LongKey = strings.Replace(nodeKey[1:], "/", "^", -1)
				kval.LongKey = nodeKey[1:]
				kval.Key = keys[j]
				kval.Value = string(rsp.Kvs[i].Value)
				if kval.Value == DirValue {
					kval.IsDir = true
				}
				break
			}
			kvals = kval.Nodes
		}
	}
	return results, nil
}

func (ec *EtcdConfig) Del(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ec.timeout)
	defer cancel()

	realKey := ec.prefix + key
	_, err := ec.client.Delete(ctx, realKey)
	if err != nil {
		return err
	}
	return nil
}

func (ec *EtcdConfig) Put(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ec.timeout)
	defer cancel()

	realKey := ec.prefix + key
	_, err := ec.client.Put(ctx, realKey, value)
	if err != nil {
		return err
	}
	return nil
}
