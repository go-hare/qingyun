package config

import (
	"time"
)

type Option func(o *EtcdOpts)

type EtcdOpts struct {
	address []string
	prefix  string
	timeout time.Duration
}

func NewOptions(opts ...Option) EtcdOpts {
	options := EtcdOpts{
		address: nil,
		timeout: 0,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

func WithAddress(a []string) Option {
	return func(o *EtcdOpts) {
		o.address = a
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *EtcdOpts) {
		o.timeout = timeout
	}
}

func WithPrefix(p string) Option {
	return func(o *EtcdOpts) {
		o.prefix = p
	}
}
