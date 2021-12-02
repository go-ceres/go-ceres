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
	"github.com/go-ceres/go-ceres/logger"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// loggerMiddleware 日志中间件
func loggerMiddleware(log *logger.Logger, slowQueryThresholdInMilli int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		var now = time.Now()
		var fields = make([]logger.Field, 0, 0)
		var errBlack bool
		defer func() {
			// 花费时长
			fields = append(fields, logger.FieldString("cost", time.Since(now).String()))
			// 超时记录
			if slowQueryThresholdInMilli > 0 {
				if cost := int64(time.Since(now)) / 1e6; cost > slowQueryThresholdInMilli {
					fields = append(fields, zap.Int64("slow", cost))
				}
			}
			// 错误日志
			if rec := recover(); rec != nil {
				if ne, ok := rec.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							errBlack = true
						}
					}
				}
				err := rec.(error)
				fields = append(fields, logger.FieldErr(err))
				logger.Errord("gin request error", fields...)
				if errBlack {
					c.Error(err)
					c.Abort()
					return
				}
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			// 正常日志
			fields = append(fields,
				logger.FieldString("method", c.Request.Method),
				logger.FieldString("path", c.Request.URL.Path),
				logger.FieldString("host", c.Request.Host),
			)
			logger.Infod("gin request success", fields...)
		}()
		c.Next()
	}
}
