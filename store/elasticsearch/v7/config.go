//    Copyright 2021. Go-Ceres
//    Author https://github.com/go-ceres/go-ceres
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package elasticsearch

import (
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/logger"
	"github.com/olivere/elastic/v7"
	"net/http"
)

type Config struct {
	Address    []string     // es连接地址
	Scheme     string       // http协议
	Username   string       // es用户名
	Password   string       // es密码
	httpClient *http.Client // http客户端
	options    []elastic.ClientOptionFunc
	logger     *logger.Logger // 日志组件
}

func DefaultConfig() *Config {
	return &Config{
		Address: []string{"http://127.0.0.1:9200"},
		logger:  logger.FrameLogger.With(logger.FieldMod("store.elastic")),
	}
}

// RawConfig 根据配置key读取配置
func RawConfig(key string) *Config {
	conf := DefaultConfig()
	err := config.Get(key).Scan(conf)
	if err != nil {
		conf.logger.Panicd("parse config error", logger.FieldAny("key", key), logger.FieldValue(conf))
	}
	return conf
}

// ScanConfig 根据配置名称读取配置
func ScanConfig(name string) *Config {
	return RawConfig("ceres.store.elastic." + name)
}

// WithTransport 单独设置http客户端的transport
func (c *Config) WithTransport(transport *http.Transport) *Config {
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	c.httpClient.Transport = transport
	return c
}

// Build 构建客户端
func (c *Config) Build() *Client {
	var options []elastic.ClientOptionFunc
	options = append(options, elastic.SetURL(c.Address...))
	// 如果有协议
	if len(c.Scheme) > 0 {
		options = append(options, elastic.SetScheme(c.Scheme))
	}
	// 如果有账号密码
	if c.Username != "" && c.Password != "" {
		options = append(options, elastic.SetBasicAuth(c.Username, c.Password))
	}
	// 如果设置了transport
	if c.httpClient != nil {
		options = append(options, elastic.SetHttpClient(c.httpClient))
	}

	c.options = options
	return newClient(c)
}
