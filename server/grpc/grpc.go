//   Copyright 2021 Go-Ceres
//   Author https://github.com/go-ceres/go-ceres
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package grpc

import (
	"github.com/go-ceres/go-ceres/logger"
	"github.com/go-ceres/go-ceres/server"
	"google.golang.org/grpc"
	"net"
	"time"
)

type grpcServer struct {
	Server   *grpc.Server
	listener net.Listener
	*Config
}

// NewServer 新建服务
func newServer(c *Config) *grpcServer {
	var streamInterceptors []grpc.StreamServerInterceptor
	if c.Debug {
		streamInterceptors = append(streamInterceptors, debugStreamServerInterceptor(c.logger, c.ServerSlowThreshold))
	}
	streamInterceptors = append(streamInterceptors, c.streamInterceptors...)

	var unaryInterceptors []grpc.UnaryServerInterceptor
	if c.Debug {
		unaryInterceptors = append(unaryInterceptors, debugUnaryServerInterceptor(c.logger, c.ServerSlowThreshold))
	}
	unaryInterceptors = append(unaryInterceptors, c.unaryInterceptors...)

	c.serverOptions = append(c.serverOptions,
		grpc.StreamInterceptor(StreamInterceptorChain(streamInterceptors...)),
		grpc.UnaryInterceptor(UnaryInterceptorChain(unaryInterceptors...)),
	)

	newServer := grpc.NewServer(c.serverOptions...)
	listener, err := net.Listen(c.Network, c.Address())
	if err != nil {
		c.logger.Panicd("new grpc server err", logger.FieldErr(err))
	}
	c.Port = listener.Addr().(*net.TCPAddr).Port

	return &grpcServer{
		Server:   newServer,
		listener: listener,
		Config:   c,
	}
}

// Start 启动grpc服务
func (s *grpcServer) Start() error {
	return s.Server.Serve(s.listener)
}

// Stop 停止服务
func (s *grpcServer) Stop() error {
	s.Server.Stop()
	return nil
}

// Info 服务信息
func (s *grpcServer) Info() *server.ServiceInfo {
	address := s.listener.Addr().String()
	ip, port, _ := net.SplitHostPort(address)
	internalIp := s.getInterfaceIp()
	if ip == "0.0.0.0" {
		if s.ping(internalIp + ":" + port) {
			ip = internalIp
		}
	}
	if s.Config.PlainTextAddress != "" {
		ip = s.Config.PlainTextAddress
	}
	address = ip + ":" + port
	info := server.ApplyOptions(
		server.WithAddress(address),
		server.WithScheme("grpc"),
		server.WithMetadata("app_host", address),
	)
	return info
}

// Ping 检查ip是否能通
func (s *grpcServer) ping(addr string) bool {
	// 3 秒超时
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return false
	} else {
		if conn != nil {
			_ = conn.Close()
			return true
		}
		return false
	}
}

// getInterfaceIp 获取内网ip
func (s *grpcServer) getInterfaceIp() string {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {

		if ipNet, isIpNet := addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}
