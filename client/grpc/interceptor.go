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

package grpc

import (
	"context"
	"fmt"
	"github.com/go-ceres/go-ceres/errors"
	"github.com/go-ceres/go-ceres/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"time"
)

// debugUnaryClientInterceptor 调试拦截器
func debugUnaryClientInterceptor(log *logger.Logger, addr string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var p peer.Peer
		addr := fmt.Sprintf("[%s]", addr)
		if remote, ok := peer.FromContext(ctx); ok && remote.Addr != nil {
			addr = addr + "(" + remote.Addr.String() + ")"
		}
		log.Infod("before Call", logger.FieldAny("addr", addr), logger.FieldAny("method", method), logger.FieldAny("req", req))
		err := invoker(ctx, method, req, reply, cc, append(opts, grpc.Peer(&p))...)
		if err != nil {
			log.Errord("after Call", logger.FieldAny("addr", addr), logger.FieldErr(err))
		} else {
			log.Infod("after Call", logger.FieldAny("addr", addr), logger.FieldAny("method", method), logger.FieldAny("reply", reply))
		}

		return err
	}
}

// timeoutUnaryClientInterceptor gRPC客户端超时拦截器
func timeoutUnaryClientInterceptor(log *logger.Logger, timeout time.Duration, slowThreshold time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		now := time.Now()
		// 若无自定义超时设置，默认设置超时
		_, ok := ctx.Deadline()
		if !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		err := invoker(ctx, method, req, reply, cc, opts...)
		du := time.Since(now)
		addr := fmt.Sprintf("[%s]", cc.Target())
		if remote, ok := peer.FromContext(ctx); ok && remote.Addr != nil {
			addr = addr + "(" + remote.Addr.String() + ")"
		}

		if du > slowThreshold {
			log.Error("slow",
				logger.FieldErr(errors.New(errors.CodeCallGrpcUnaryTimeout, errors.MsgCallGrpcUnaryTimeout)),
				logger.FieldAny("method", method),
				logger.FieldAny("cost", du),
				logger.FieldAny("addr", addr),
			)
		}
		return err
	}
}
