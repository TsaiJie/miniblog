// Copyright 2022 Innkeeper tsai <mengjietsai@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/tsai/miniblog.

package miniblog

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/miniblog/internal/pkg/log"
	"github.com/miniblog/pkg/version/verflag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
)

var cfgFile string

// NewMiniBlogCommand 创建一个cobra.Command对象，之后可以使用Command对象的Execute方法来启动应用程序
func NewMiniBlogCommand() *cobra.Command {
	cmd := &cobra.Command{
		//	指定命令的名字，该名字会出现在帮助信息中
		Use: "miniblog",
		//	命令的简短描述
		Short: "A good Go practical project",
		// 命令的详细描述
		Long: `A good Go practical project, used to create user with basic information.
		Find more miniblog information at:
        	https://github.com/tsaijie/miniblog#readme`,
		// 命令出错时，不打印帮助信息。不需要打印帮助信息，设置为 true 可以保持命令出错时一眼就能看到错误信息
		SilenceUsage: true,
		// 这里设置命令运行时，不需要指定命令行参数
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
		// 指定调用 cmd.Execute() 时，执行的 Run 函数，函数执行失败会返回错误信息
		RunE: func(cmd *cobra.Command, args []string) error {
			// 如果 `--version=true`，则打印版本并退出
			verflag.PrintAndExitIfRequested()

			// 初始化日志
			log.Init(logOptions())
			// Sync 将缓存中的日志刷新到磁盘文件中
			defer log.Sync()

			return run()
		},
	}

	// 以下设置，使得 initConfig 函数在每个命令运行时都会被调用以读取配置
	cobra.OnInitialize(initConfig)
	// 在这里定义标志和配置设置
	// cobra 支持持久性标志，该标志可用于它所分配的命令以及该命令下的每个子命令
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path to the miniblog configuration file. Empty string for no configuration file.")
	// cobra 也支持本地标志，本地标志只能在其所绑定的命令上使用
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// 添加 --version 标志
	verflag.AddFlags(cmd.PersistentFlags())
	return cmd
}

// run 函数是实际的业务代码入口函数.
func run() error {

	// 设置gin的模式
	gin.SetMode(viper.GetString("runmode"))

	// 创建Gin引擎
	g := gin.New()

	// 注册404Handler
	g.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"code": 10003, "message": "Page not found."})
	})

	// 注册/healthz handler
	g.GET("/healthz", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 创建HTTP服务器
	httpsrv := &http.Server{
		Addr:    viper.GetString("addr"),
		Handler: g,
	}
	// 运行HTTP服务器
	// 打印一条日志，用来提示HTTP服务已经起来，方便排查
	log.Infow("Start to listening the incoming requests on http address", "addr", viper.GetString("addr"))

	if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalw(err.Error())
	}
	return nil
}
