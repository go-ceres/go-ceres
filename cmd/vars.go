package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const ceresVersion = "0.1.0"

var (
	startTime string // 项目运行时间
	goVersion string // 项目运行时golang版本
	appRegion string // 数据中心
	appZone   string // 区域
)

var (
	appName     string // 应用名称
	usage       string // 应用名称后面的介绍
	description string // 应用简介
	version     string // 应用版本信息
	commit      string // git 提交信息
	branch      string // git 分支信息
	hostName    string // 服务器名称
	buildHost   string // 构建时的host
	buildTime   string // 应用构建时间
	buildUser   string // 构建的用户
	buildStatus string // 应用构建状态
)

func init() {
	if appName == "" {
		appName = filepath.Base(os.Args[0])
	}

	name, err := os.Hostname()
	if err != nil {
		name = "unknown"
	}
	hostName = name
	// 应用启动时间
	startTime = time.Now().Format("2006-01-02 15:04:05")
	// 设置构建时间
	SetBuildTime(buildTime)
	goVersion = runtime.Version()
}

// SetBuildTime 设置构建时间
func SetBuildTime(s string) {
	buildTime = strings.Replace(s, "--", " ", 1)
}

// GetCeresVersion 获取go-ceres版本
func GetCeresVersion() string {
	return ceresVersion
}

// SetAppName 设置应用名称
func SetAppName(name string) {
	appName = name
}

// GetAppName 获取应用名称
func GetAppName() string {
	return appName
}

// SetAppVersion 设置应用版本
func SetAppVersion(v string) {
	version = v
}

// GetAppVersion 获取应用版本
func GetAppVersion() string {
	return version
}

// GetRegion 获取数据中心
func GetRegion() string {
	return appRegion
}

// GetZone 获取区域
func GetZone() string {
	return appZone
}

// GetBuildTime 获取构建时间
func GetBuildTime() string {
	return buildTime
}

// GetBuildUser 获取构建用户
func GetBuildUser() string {
	return buildUser
}

// GetBuildHost 获取构建host
func GetBuildHost() string {
	return buildHost
}

// FrameVersion 获取框架版本
func FrameVersion() string {
	return ceresVersion
}

// GetStartTime 获取启动时间
func GetStartTime() string {
	return startTime
}

// GetHostname 获取hostname
func GetHostname() string {
	return hostName
}

// ShowVersion 输出版本信息
func ShowVersion() {
	fmt.Printf("%-9s]> %-30s => %s\n", "go-ceres", color.RedString("name"), color.BlueString(appName))
	fmt.Printf("%-9s]> %-30s => %s\n", "go-ceres", color.RedString("version"), color.BlueString(version))
	if len(branch) != 0 {
		fmt.Printf("%-9s]> %-30s => %s\n", "go-ceres", color.RedString("branch"), color.BlueString(branch))
	}

	if len(commit) != 0 {
		fmt.Printf("%-9s]> %-30s => %s\n", "go-ceres", color.RedString("commit"), color.BlueString(commit))
	}
	fmt.Printf("%-9s]> %-30s => %s\n", "go-ceres", color.RedString("ceresVersion"), color.BlueString(ceresVersion))
	fmt.Printf("%-9s]> %-30s => %s\n", "go-ceres", color.RedString("goVersion"), color.BlueString(goVersion))
	fmt.Printf("%-9s]> %-30s => %s\n", "go-ceres", color.RedString("buildUser"), color.BlueString(buildUser))
	fmt.Printf("%-9s]> %-30s => %s\n", "go-ceres", color.RedString("buildHost"), color.BlueString(buildHost))
	fmt.Printf("%-9s]> %-30s => %s\n", "go-ceres", color.RedString("buildTime"), color.BlueString(buildTime))
	fmt.Printf("%-9s]> %-30s => %s\n", "go-ceres", color.RedString("buildStatus"), color.BlueString(buildStatus))
}
