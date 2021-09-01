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

package etcdv3

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/go-ceres/go-ceres/errors"
	"github.com/go-ceres/go-ceres/logger"
	"github.com/jinzhu/copier"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"google.golang.org/grpc"
	"io/ioutil"
)

type Client struct {
	*clientv3.Client
	config *Config
}

// newClient 创建Etcd客户端
func newClient(c *Config) *Client {
	conf := clientv3.Config{}
	err := copier.Copy(&conf, c)
	if err != nil {
		c.logger.Panicd("etcd client init error", logger.FieldErr(err))
	}

	// 使用安全连接
	if c.Secure {
		conf.DialOptions = append(conf.DialOptions, grpc.WithInsecure())
	}

	// tls证书
	tlsEnabled := false
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}
	if c.CaCert != "" {
		certBytes, err := ioutil.ReadFile(c.CaCert)
		if err != nil {
			c.logger.Panicd("parse CaCert failed", logger.FieldErr(err))
		}

		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(certBytes)

		if ok {
			tlsConfig.RootCAs = caCertPool
		}
		tlsEnabled = true
	}
	if c.CertFile != "" && c.KeyFile != "" {
		tlsCert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
		if err != nil {
			c.logger.Panicd("load CertFile or KeyFile failed", logger.FieldValue(c), logger.FieldErr(err))
		}
		tlsConfig.Certificates = []tls.Certificate{tlsCert}
		tlsEnabled = true
	}
	if tlsEnabled {
		conf.TLS = tlsConfig
	}
	client, err := clientv3.New(conf)
	if err != nil {
		c.logger.Panicd("client etcd start panic", logger.FieldMod(errors.ModClientEtcd), logger.FieldErr(err), logger.FieldValue(c))
	}
	if c.logger != nil {
		client = client.WithLogger(c.logger.AddCallerSkip(-1).ZapLogger())
	}
	cli := &Client{
		Client: client,
		config: c,
	}
	c.logger.Infod("etcd init success", logger.FieldAny("Endpoints", c.Endpoints))
	return cli
}

// GetLeaseSession 创建一个session会话
func (c *Client) GetLeaseSession(opts ...concurrency.SessionOption) (leaseSession *concurrency.Session, err error) {
	return concurrency.NewSession(c.Client, opts...)
}
