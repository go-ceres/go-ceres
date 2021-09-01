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
	"context"
	"fmt"
	"github.com/go-ceres/go-ceres/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"runtime"
	"strings"
	"time"
)

// 日志拦截器
func debugUnaryServerInterceptor(log *logger.Logger, serverSlowThreshold int64) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 开始时间
		startTime := time.Now()
		// 事件类型
		event := "normal"
		fields := make([]logger.Field, 0)
		duration := int64(0)
		defer func() {
			if serverSlowThreshold > 0 {
				duration = int64(time.Since(startTime)) / 1e6
				if duration > serverSlowThreshold {
					event = "overtime"
				}
			}
			if rec := recover(); rec != nil {
				switch rec := rec.(type) {
				case error:
					err = rec
				default:
					err = fmt.Errorf("%v", rec)
				}
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				fields = append(fields, logger.FieldAny("stack", stack))
				event = "recover"
			}
			fields = append(fields,
				logger.FieldAny("grpc interceptor type", "unary"),
				logger.FieldAny("method", info.FullMethod),
				logger.FieldAny("duration", duration),
				logger.FieldAny("event", event),
			)
			for key, val := range getPeer(ctx) {
				fields = append(fields, logger.FieldAny(key, val))
			}

			if err != nil {
				fields = append(fields, logger.FieldErr(err))
				log.Errord("access", fields...)
			} else {
				log.Infod("access", fields...)
			}
		}()
		return handler(ctx, req)
	}
}

// 日志拦截器
func debugStreamServerInterceptor(log *logger.Logger, serverSlowThreshold int64) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) (err error) {
		// 开始时间
		startTime := time.Now()
		// 事件类型
		event := "normal"
		fields := make([]logger.Field, 0)
		duration := int64(0)
		defer func() {
			if serverSlowThreshold > 0 {
				duration = int64(time.Since(startTime)) / 1e6
				if duration > serverSlowThreshold {
					event = "overtime"
				}
			}
			if rec := recover(); rec != nil {
				switch rec := rec.(type) {
				case error:
					err = rec
				default:
					err = fmt.Errorf("%v", rec)
				}
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				fields = append(fields, logger.FieldAny("stack", stack))
				event = "recover"
			}
			fields = append(fields,
				logger.FieldAny("grpc interceptor type", "stream"),
				logger.FieldAny("method", info.FullMethod),
				logger.FieldAny("duration", duration),
				logger.FieldAny("event", event),
			)
			for key, val := range getPeer(ss.Context()) {
				fields = append(fields, logger.FieldAny(key, val))
			}

			if err != nil {
				fields = append(fields, logger.FieldErr(err))
				log.Errord("access", fields...)
			} else {
				log.Infod("access", fields...)
			}
		}()
		return handler(srv, ss)
	}
}

// getPeer 解析
func getPeer(ctx context.Context) map[string]string {
	var peerMeta = make(map[string]string)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// 获取应该id
		if val, ok := md["aid"]; ok {
			peerMeta["aid"] = strings.Join(val, ";")
		}
		var clientIP string
		if val, ok := md["client-ip"]; ok {
			clientIP = strings.Join(val, ";")
		} else {
			client, ok := peer.FromContext(ctx)
			if ok {
				clientIP = client.Addr.String()
			}
		}
		peerMeta["clientIP"] = clientIP
		// 获取客户端主机名
		if val, ok := md["client-host"]; ok {
			peerMeta["host"] = strings.Join(val, ";")
		}
	}
	return peerMeta

}
