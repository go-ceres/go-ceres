package etcd

import (
	"context"
	"github.com/go-ceres/go-ceres/config"
	"go.etcd.io/etcd/client/v3"
	"time"
)

type Config struct {
	*clientv3.Config
	Prefix     string // etcd配置路径
	TrimPrefix string // 删除掉的头部字符串
	Encoding   string // 加解密
	Ctx        context.Context
}

// 默认的配置信息
func DefaultConfig() *Config {
	conf := &Config{
		Config: &clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: 5 * time.Second,
		},
		Ctx:      context.Background(),
		Prefix:   DefaultPrefix,
		Encoding: "json",
	}
	return conf
}

// RawConfig 返回配置
func RawConfig(key string) *Config {
	conf := DefaultConfig()
	if err := config.Get(key).Scan(conf); err != nil {
		panic(err)
	}
	return conf
}

// ScanConfig 返回配置
func ScanConfig(name string) *Config {
	return RawConfig("ceres.config.source.etcd." + name)
}

// Addr 连接地址
func (c *Config) WithEndpoints(addrs ...string) *Config {
	c.Endpoints = addrs
	return c
}

// Prefix 前缀
func (c *Config) WithPrefix(prefix string) *Config {
	c.Prefix = prefix
	return c
}

// TrimPrefix 去掉前缀前缀
func (c *Config) WithStripPrefix(trimPrefix string) *Config {
	c.TrimPrefix = trimPrefix
	return c
}

// Build 根据配置信息
func (c *Config) Build() config.Source {
	return NewSource(c)
}
