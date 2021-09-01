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
	"github.com/go-ceres/go-ceres/logger"
	"github.com/olivere/elastic"
)

type Client struct {
	*elastic.Client // 客户端
}

func newClient(c *Config) *Client {
	client, err := elastic.NewClient(c.options...)
	if err != nil {
		c.logger.Panicd("init elasticsearch client error", logger.FieldErr(err), logger.FieldAny("config", c))
	}
	c.logger.Infod("init elasticsearch client success", logger.FieldAny("config", c))
	return &Client{
		Client: client,
	}
}
