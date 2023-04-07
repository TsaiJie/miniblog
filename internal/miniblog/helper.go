// Copyright 2022 Innkeeper tsai <mengjietsai@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/tsai/miniblog.

package miniblog

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

const (
	// recommendedHomeDir 定义放置miniblog服务配置的默认目录
	recommendedHomeDir = ".miniblog"
	// defaultConfigName 指定miniblog服务的默认配置文件名
	defaultConfigName = "miniblog.yaml"
)

// initConfig 设置需要读取的配置文件名、环境变量、并读取配置文件到viper中
func initConfig() {
	if cfgFile != "" {
		//	从命令行 --config指定的配置文件中读取
		viper.SetConfigFile(cfgFile)
	} else {
		//查找主目录
		home, err := os.UserHomeDir()
		//	如果获取用户主目录失败，打印`Error：xxx`错误，并退出程序
		cobra.CheckErr(err)

		// 用$HOME/<recommendedHomeDir>目录加入到配置文件的搜索路径中
		viper.AddConfigPath(filepath.Join(home, recommendedHomeDir))

		// 把当前目录加入到配置文件的搜索路径中
		viper.AddConfigPath(".")

		//	设置配置文件格式为 YAML
		viper.SetConfigType("yaml")

		//	配置文件名称（没有文件扩展名）
		viper.SetConfigName(defaultConfigName)
	}
	//	 读取匹配的环境变量
	viper.AutomaticEnv()

	// 读取环境变量的前缀为MINIBLOG ，如果是miniblog，自动转变为大写。
	viper.SetEnvPrefix("MINIBLOG")

	// 以下两行 将viper.Get(key) key字符串中'.'和'-'替换为'_'
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// 读取配置文件，如果指定了配置文件名，则使用指定的配置文件，否则在注册的搜索路径中搜索
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

}
