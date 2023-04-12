package mysql

import "xorm.io/xorm"

type Options struct {
	Host      string
	User      string
	Pass      string
	DBName    string
	IfShowSql bool
	IfSyncDB  bool
	AfterInit func(x *xorm.Engine)
}

type Option func(o *Options)

func WithHost(host string) Option {
	return func(o *Options) {
		o.Host = host
	}
}

func WithUser(user string) Option {
	return func(o *Options) {
		o.User = user
	}
}

func WithPass(pass string) Option {
	return func(o *Options) {
		o.Pass = pass
	}
}

func WithDBName(name string) Option {
	return func(o *Options) {
		o.DBName = name
	}
}

func IfShowSql(b bool) Option {
	return func(o *Options) {
		o.IfShowSql = b
	}
}

func IfSyncDB(b bool) Option {
	return func(o *Options) {
		o.IfSyncDB = b
	}
}

func AfterInit(after func(x *xorm.Engine)) Option {
	return func(o *Options) {
		o.AfterInit = after
	}
}
