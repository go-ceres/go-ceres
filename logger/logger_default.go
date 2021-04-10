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
	"github.com/go-ceres/go-ceres/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

type (
	Func   func(string, ...zap.Field)
	Field  = zap.Field
	Level  = zapcore.Level
	logger struct {
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
func newLogger(c *Config) *logger {
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
	if c.Rotate {
		ws["rotate"] = zapcore.AddSync(c.RotateConfig.Build())
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
			if key == "stdout" {
				encoderConfig := defaultEncoderConfig()
				encoderConfig.EncodeLevel = debugEncodeLevel
				encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
				cores = append(cores, zapcore.NewCore(
					zapcore.NewConsoleEncoder(*encoderConfig),
					w,
					lv,
				))
			} else {
				encoderConfig := defaultEncoderConfig()
				cores = append(cores, zapcore.NewCore(
					zapcore.NewJSONEncoder(*encoderConfig),
					w,
					lv,
				))
			}
		}
		core = zapcore.NewTee(cores...)
	}
	zapLogger := zap.New(
		core,
		zapOptions...,
	)
	zapLogger.WithOptions(zap.Hooks())
	return &logger{
		logger:        zapLogger,
		lv:            &lv,
		config:        *c,
		sugaredLogger: zapLogger.Sugar(),
	}
}

func (l *logger) AutoLevel(key string) {
	config.OnChange(func(v config.Values) {
		lvText := strings.ToLower(v.Get(key).String(""))
		if lvText != "" {
			l.Info("update level", String("level", lvText))
			if err := l.lv.UnmarshalText([]byte(lvText)); err != nil {
				l.Error("UnmarshalText error: " + err.Error())
			}
		}
	})
	l.sugaredLogger.Debug()
}

// 主动设置等级
func (l *logger) SetLevel(level Level) {
	l.lv.SetLevel(level)
}
func (l *logger) With(fields ...Field) Logger {
	logger := l.logger.With(fields...)
	clone := l.clone()
	clone.logger = logger
	clone.sugaredLogger = clone.logger.Sugar()
	return clone
}

// 跨栈调用，增加栈
func (l *logger) AddCallerSkip(i int) Logger {
	clone := l.clone()
	clone.logger = clone.logger.WithOptions(zap.AddCallerSkip(i))
	clone.sugaredLogger = clone.logger.Sugar()
	return clone
}

// clone  克隆
func (l *logger) clone() *logger {
	clone := *l
	return &clone
}

// Sugar
func (l *logger) Debug(args ...interface{}) {
	l.sugaredLogger.Debug(args...)
}
func (l *logger) Info(args ...interface{}) {
	l.sugaredLogger.Info(args...)
}
func (l *logger) Warn(args ...interface{}) {
	l.sugaredLogger.Warn(args...)
}
func (l *logger) Error(args ...interface{}) {
	l.sugaredLogger.Error(args...)
}
func (l *logger) DPanic(args ...interface{}) {
	l.sugaredLogger.DPanic(args...)
}
func (l *logger) Panic(args ...interface{}) {
	l.sugaredLogger.Panic(args...)
}

// Sugar f
func (l *logger) Debugf(template string, args ...interface{}) {
	l.sugaredLogger.Debugf(template, args...)
}
func (l *logger) Infof(template string, args ...interface{}) {
	l.sugaredLogger.Infof(template, args...)
}
func (l *logger) Warnf(template string, args ...interface{}) {
	l.sugaredLogger.Warnf(template, args...)
}
func (l *logger) Errorf(template string, args ...interface{}) {
	l.sugaredLogger.Errorf(template, args...)
}
func (l *logger) DPanicf(template string, args ...interface{}) {
	l.sugaredLogger.DPanicf(template, args...)
}
func (l *logger) Panicf(template string, args ...interface{}) {
	l.sugaredLogger.Panicf(template, args...)
}

// Sugar w
func (l *logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Debugw(msg, keysAndValues...)
}
func (l *logger) Infow(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Infow(msg, keysAndValues...)
}
func (l *logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Warnw(msg, keysAndValues...)
}
func (l *logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Errorw(msg, keysAndValues...)
}
func (l *logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.DPanicw(msg, keysAndValues...)
}
func (l *logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Panicw(msg, keysAndValues...)
}

// DeSugar
func (l *logger) Debugd(msg string, fields ...Field) {
	l.logger.Debug(msg, fields...)
}
func (l *logger) Infod(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}
func (l *logger) Warnd(msg string, fields ...Field) {
	l.logger.Warn(msg, fields...)
}
func (l *logger) Errord(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}
func (l *logger) DPanicd(msg string, fields ...Field) {
	l.logger.DPanic(msg, fields...)
}
func (l *logger) Panicd(msg string, fields ...Field) {
	l.logger.Panic(msg, fields...)
}
