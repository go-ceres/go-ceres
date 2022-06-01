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

package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ceres/go-ceres/errors"
	"github.com/go-ceres/go-ceres/logger"
	"github.com/go-ceres/go-ceres/server"
	"net"
	"net/http"
	"reflect"
	"strconv"
)

type Server struct {
	*gin.Engine
	Server   *http.Server
	listener net.Listener
	Config   *Config
}

// ServiceDesc 服务描述
type ServiceDesc struct {
	ServiceName string       // 服务名称
	HandlerType interface{}  // 服务类型
	Routers     []RouterDesc // 路由描述
}

// methodHandler 服务回调方法
type methodHandler func(srv interface{}, ctx *Context, dec func(interface{}) error) (interface{}, error)

// RouterDesc 路由描述
type RouterDesc struct {
	Path    string
	Method  string
	Handler methodHandler
}

// newGinServer 创建gin服务
func newGinServer(config *Config) *Server {
	listener, err := net.Listen("tcp", config.Address())
	if err != nil {
		config.logger.Panicd("new gin server error", logger.FieldErr(err))
	}
	config.Port = listener.Addr().(*net.TCPAddr).Port
	gin.SetMode(config.Mode)
	return &Server{
		Engine:   gin.New(),
		Config:   config,
		listener: listener,
	}
}

// Upgrade 升级协议为websocket
func (s *Server) Upgrade(ws *WebSocket) gin.IRoutes {
	return s.GET(ws.Pattern, func(c *gin.Context) {
		ws.Upgrade(c.Writer, c.Request)
	})
}

// Start 启动服务
func (s *Server) Start() error {
	// 打印路由
	for _, route := range s.Engine.Routes() {
		s.Config.logger.Infod("add route", logger.FieldString("method", route.Method), logger.FieldString("path", route.Path))
	}
	// 配置服务信息
	s.Server = &http.Server{
		Addr:    s.Config.Address(),
		Handler: s,
	}
	// 启动服务
	err := s.Server.Serve(s.listener)
	if err == http.ErrServerClosed {
		s.Config.logger.Infod("gin server close", logger.FieldString("address", s.Config.Address()))
	}
	return err
}

// Stop 停止服务
func (s *Server) Stop() error {
	return s.Server.Close()
}

// Info 获取服务信息
func (s *Server) Info() *server.ServiceInfo {
	address := s.listener.Addr().String()
	if s.Config.PlainTextAddress != "" {
		address = s.Config.PlainTextAddress
	}
	address += ":" + strconv.Itoa(s.Config.Port)
	info := server.ApplyOptions(
		server.WithAddress(address),
		server.WithScheme("http"),
		server.WithMetadata("app_host", address),
	)
	return info
}

// RegisterService 注册服务
func (s *Server) RegisterService(sd *ServiceDesc, ss interface{}) {
	if ss != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			logger.Panicf("gin: Server.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
		}
	}
	s.register(sd, ss)
}

// register 根据描述文件注册路由
func (s *Server) register(sd *ServiceDesc, ss interface{}) {

	for _, router := range sd.Routers {
		s.Engine.Handle(router.Method, router.Path, s.buildHandler(ss, router.Handler))
	}
}

func (s *Server) buildHandler(ss interface{}, handler methodHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handler := handler
		df := func(v interface{}) error {
			return ctx.ShouldBind(v)
		}
		resp, err := handler(ss, ctx, df)
		if err != nil {
			if e, ok := err.(*errors.Error); ok {
				ctx.String(e.Code, e.Msg)
				return
			}
			ctx.String(500, err.Error())
			return
		}
		ctx.JSON(200, resp)
	}
}
