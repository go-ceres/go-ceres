// Copyright 2020 Go-Ceres
// Author https://github.com/go-ceres/go-ceres
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"fmt"
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/logger/encoder"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

type (
	Func   func(string, ...zap.Field)
	Field  = zap.Field
	Level  = zapcore.Level
	Logger struct {
		logger        *zap.Logger
		lv            *zap.AtomicLevel
		core          zapcore.Core
		config        Config
		sugaredLogger *zap.SugaredLogger
		encoderConfig *zapcore.EncoderConfig
	}
)

var (
	// String ...
	String = zap.String
	// Any ...
	Any = zap.Any
	// Int64 ...
	Int64 = zap.Int64
	// Int ...
	Int = zap.Int
	// Int32 ...
	Int32 = zap.Int32
	// Uint ...
	Uint = zap.Uint
	// Duration ...
	Duration = zap.Duration
	// Durationp ...
	Durationp = zap.Durationp
	// Object ...
	Object = zap.Object
	// Namespace ...
	Namespace = zap.Namespace
	// Reflect ...
	Reflect = zap.Reflect
	// Skip ...
	Skip = zap.Skip()
	// ByteString ...
	ByteString = zap.ByteString
)

// 根据配置信息创建logger
func newLogger(c *Config) *Logger {
	zapOptions := make([]zap.Option, 0)
	zapOptions = append(zapOptions, zap.AddStacktrace(zap.DPanicLevel))
	if c.AddCaller {
		zapOptions = append(zapOptions, zap.AddCaller(), zap.AddCallerSkip(c.CallerSkip))
	}
	if len(c.Fields) > 0 {
		zapOptions = append(zapOptions, zap.Fields(c.Fields...))
	}
	var ws = make(map[string]zapcore.WriteSyncer)
	if c.Stdout {
		ws["stdout"] = os.Stdout
	}
	if len(c.writer) > 0 {
		for s, writer := range c.writer {
			ws[s] = zapcore.AddSync(writer)
		}
	}
	var lv zap.AtomicLevel
	// 如果日志等级不为空
	if c.Level != "" {
		if err := lv.UnmarshalText([]byte(c.Level)); err != nil {
			panic(err)
		}
	}
	// 如果开启了debug模式
	if c.Debug {
		lv = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	} else {
		lv = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	core := c.Core
	if core == nil {
		cores := make([]zapcore.Core, 0)
		for key, w := range ws {
			encoderConfig := *c.EncoderConfig
			cores = append(cores, zapcore.NewCore(
				func() zapcore.Encoder {
					if key == "stdout" {
						encoderConfig.EncodeLevel = debugEncodeLevel
						encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
						encoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, arrayEncoder zapcore.PrimitiveArrayEncoder) {
							arrayEncoder.AppendString(caller.FullPath())
						}
						return encoder.NewConsoleEncoder(encoderConfig)
					}
					return zapcore.NewJSONEncoder(encoderConfig)
				}(),
				w,
				lv,
			))
		}
		core = zapcore.NewTee(cores...)
	}
	zapLogger := zap.New(
		core,
		zapOptions...,
	)
	zapLogger.WithOptions(zap.Hooks())
	return &Logger{
		logger:        zapLogger,
		lv:            &lv,
		config:        *c,
		sugaredLogger: zapLogger.Sugar(),
	}
}

func (l *Logger) AutoLevel(key string) {
	config.OnChange(func(v config.Values) {
		lvText := strings.ToLower(v.Get(key).String(""))
		if lvText != "" {
			l.Info("update level", String("level", lvText))
			if err := l.lv.UnmarshalText([]byte(lvText)); err != nil {
				l.Error("UnmarshalText error: " + err.Error())
			}
		}
	})
}

// SetLevel 主动设置等级
func (l *Logger) SetLevel(level Level) *Logger {
	clone := l.clone()
	clone.lv.SetLevel(level)
	return clone
}
func (l *Logger) With(fields ...Field) *Logger {
	logger := l.logger.With(fields...)
	clone := l.clone()
	clone.logger = logger
	clone.sugaredLogger = clone.logger.Sugar()
	return clone
}

// AddCallerSkip 跨栈调用，增加栈
func (l *Logger) AddCallerSkip(i int) *Logger {
	clone := l.clone()
	clone.logger = clone.logger.WithOptions(zap.AddCallerSkip(i))
	clone.sugaredLogger = clone.logger.Sugar()
	return clone
}

// clone  克隆
func (l *Logger) clone() *Logger {
	clone := *l
	return &clone
}

// Sugar
func (l *Logger) Debug(args ...interface{}) {
	l.sugaredLogger.Debug(args...)
}
func (l *Logger) Info(args ...interface{}) {
	l.sugaredLogger.Info(args...)
}
func (l *Logger) Warn(args ...interface{}) {
	l.sugaredLogger.Warn(args...)
}
func (l *Logger) Error(args ...interface{}) {
	l.sugaredLogger.Error(args...)
}
func (l *Logger) DPanic(args ...interface{}) {
	l.sugaredLogger.DPanic(args...)
}
func (l *Logger) Panic(args ...interface{}) {
	l.sugaredLogger.Panic(args...)
}

// Sugar f
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugaredLogger.Debugf(template, args...)
}
func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugaredLogger.Infof(template, args...)
}
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugaredLogger.Warnf(template, args...)
}
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugaredLogger.Errorf(template, args...)
}
func (l *Logger) DPanicf(template string, args ...interface{}) {
	l.sugaredLogger.DPanicf(template, args...)
}
func (l *Logger) Panicf(template string, args ...interface{}) {
	l.sugaredLogger.Panicf(template, args...)
}

// Sugar w
func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Debugw(msg, keysAndValues...)
}
func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Infow(msg, keysAndValues...)
}
func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Warnw(msg, keysAndValues...)
}
func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Errorw(msg, keysAndValues...)
}
func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.DPanicw(msg, keysAndValues...)
}
func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Panicw(msg, keysAndValues...)
}

// DeSugar
func (l *Logger) Debugd(msg string, fields ...Field) {
	l.logger.Debug(msg, fields...)
}
func normalizeMessage(msg string) string {
	return fmt.Sprintf("%-32s", msg)
}
func (l *Logger) Infod(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}
func (l *Logger) Warnd(msg string, fields ...Field) {
	l.logger.Warn(msg, fields...)
}
func (l *Logger) Errord(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}
func (l *Logger) DPanicd(msg string, fields ...Field) {
	l.logger.DPanic(msg, fields...)
}
func (l *Logger) Panicd(msg string, fields ...Field) {
	l.logger.Panic(msg, fields...)
}

// ZapLogger 获取zapLogger
func (l *Logger) ZapLogger() *zap.Logger {
	clone := l.clone()

	return clone.logger
}
