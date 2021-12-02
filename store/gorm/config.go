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

package gorm

import (
	"fmt"
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/errors"
	"github.com/go-ceres/go-ceres/logger"
	"gorm.io/gorm"
	log "gorm.io/gorm/logger"
	"time"
)

// Config 配置信息
type Config struct {
	Drive           string        `json:"drive"`             // 驱动
	DNS             string        `json:"dns"`               // 连接字符串
	Debug           bool          `json:"debug"`             // 是否开启调试
	MaxIdleConns    int           `json:"max_idle_conns"`    // 最大空闲连接数
	MaxOpenConns    int           `json:"max_open_conns"`    // 最大活动连接数
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"` // 连接的最大存活时间

	*GormConfig                // gorm初始化配置
	*LogConfig                 // 日志配置
	dialect     Dialector      // 驱动连接器
	logger      *logger.Logger // 日志库
}

// LogConfig 日志配置
type LogConfig struct {
	SlowThreshold time.Duration // 日志时间阈值
	Colorful      bool          // 是否开启日志颜色区别
	LogLevel      string        // 日志等级
}

// DefaultLogConfig 默认的日志配置
func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		SlowThreshold: time.Second,
		Colorful:      false,
		LogLevel:      "",
	}
}

type GormConfig gorm.Config

// DefaultConfig 默认gorm配置
func DefaultConfig() *Config {
	return &Config{
		Drive:           "mysql",
		DNS:             "",
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
		GormConfig:      &GormConfig{},
		logger:          logger.FrameLogger.With(logger.FieldMod(errors.ModStoreGorm)),
		LogConfig:       DefaultLogConfig(),
	}
}

// RawConfig 根据key读取配置信息
func RawConfig(key string) *Config {
	conf := DefaultConfig()
	err := config.Get(key).Scan(conf)
	if err != nil {
		conf.logger.Panicd("parse config error", logger.FieldErr(err), logger.FieldAny("config", conf))
	}
	return conf
}

// ScanConfig 根据名称读取配置信息
func ScanConfig(name string) *Config {
	return RawConfig("ceres.store.gorm." + name)
}

// initLogger 初始化日志
func (c *Config) initLogger() {
	// 默认日志配置
	logConf := log.Config{
		SlowThreshold: time.Second, // 慢 SQL 阈值
		LogLevel:      log.Silent,  // Log level
		Colorful:      false,       // 禁用彩色打印
	}
	// 转换等级
	if c.LogLevel != "" {
		logConf.LogLevel = ConvertLevel(c.LogLevel)
	}
	logConf.Colorful = c.Colorful
	logConf.SlowThreshold = c.SlowThreshold
	dbLog := newLog(c.logger, logConf)
	if c.Debug {
		dbLog = dbLog.LogMode(log.Info)
	}
	// gorm的配置信息
	c.GormConfig.Logger = dbLog
}

// WithDialector 单独设置Dialector
func (c *Config) WithDialector(dialect Dialector) *Config {
	c.dialect = dialect
	return c
}

// Build 构建gorm数据库链接
func (c *Config) Build() *DB {
	// 创建驱动
	if driver, ok := drivers[c.Drive]; !ok {
		c.logger.Panicd(fmt.Sprintf("%s driver is not set", c.Drive))
	} else {
		c.dialect = driver(c.DNS)
	}
	// 数据库
	db, err := Open(c.dialect, c)
	if err != nil {
		c.logger.Panicd("open gorm", logger.FieldErr(err), logger.FieldAny("value", c))
	}
	return db
}
