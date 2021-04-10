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

var (
	DefaultLogger = Config{
		Stdout:       true,
		Debug:        true,
		AddCaller:    true,
		CallerSkip:   2,
		Rotate:       false,
		RotateConfig: NewDefaultRotateConfig(),
	}.Build()
	FrameLogger = Config{
		Debug:      true,
		Stdout:     true,
		AddCaller:  true,
		CallerSkip: 2,
	}.Build()
)

// Sugar
func Debug(args ...interface{}) {
	DefaultLogger.Debug(args...)
}
func Info(args ...interface{}) {
	DefaultLogger.Info(args...)
}
func Warn(args ...interface{}) {
	DefaultLogger.Warn(args...)
}
func Error(args ...interface{}) {
	DefaultLogger.Error(args...)
}
func DPanic(args ...interface{}) {
	DefaultLogger.DPanic(args...)
}
func Panic(args ...interface{}) {
	DefaultLogger.Panic(args...)
}

// Sugar f
func Debugf(template string, args ...interface{}) {
	DefaultLogger.Debugf(template, args...)
}
func Infof(template string, args ...interface{}) {
	DefaultLogger.Infof(template, args...)
}
func Warnf(template string, args ...interface{}) {
	DefaultLogger.Warnf(template, args...)
}
func Errorf(template string, args ...interface{}) {
	DefaultLogger.Errorf(template, args...)
}
func DPanicf(template string, args ...interface{}) {
	DefaultLogger.DPanicf(template, args...)
}
func Panicf(template string, args ...interface{}) {
	DefaultLogger.Panicf(template, args...)
}

// Sugar w
func Debugw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Debugw(msg, keysAndValues...)
}
func Infow(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Infow(msg, keysAndValues...)
}
func Warnw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Warnw(msg, keysAndValues...)
}
func Errorw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Errorw(msg, keysAndValues...)
}
func DPanicw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.DPanicw(msg, keysAndValues...)
}
func Panicw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Panicw(msg, keysAndValues...)
}

// DeSugar
func Debugd(msg string, fields ...Field) {
	DefaultLogger.Debugd(msg, fields...)
}
func Infod(msg string, fields ...Field) {
	DefaultLogger.Infod(msg, fields...)
}
func Warnd(msg string, fields ...Field) {
	DefaultLogger.Warnd(msg, fields...)
}
func Errord(msg string, fields ...Field) {
	DefaultLogger.Errord(msg, fields...)
}
func DPanicd(msg string, fields ...Field) {
	DefaultLogger.DPanicd(msg, fields...)
}
func Panicd(msg string, fields ...Field) {
	DefaultLogger.Panicd(msg, fields...)
}
