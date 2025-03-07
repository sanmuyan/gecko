package cmd

import (
	"context"
	"gecko/pkg/configpost"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCtx context.Context

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "Gecko Server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := initConfig(cmd)
		if err != nil {
			logrus.Fatalf("init config error: %v", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		configpost.PostInit(rootCtx)
	},
}

var configFile string

const (
	logLevel         = 4
	serverBind       = ":8080"
	syncProjectLimit = 10
	searchProvider   = "es"
	esUrl            = "http://localhost:9200"
	reposPath        = "./repos"
	httpHost         = "http://localhost"
	gitlabUrl        = "http://localhost"
	maxFileSize      = 1024 * 1024 * 10
	maxLienLength    = 1000
	maxSearchTotal   = 1000
)

var (
	builtInDirBlacklist  = []string{"^\\.", "^node_modules$"}
	builtInFileBlacklist = []string{"\\.exe$"}
)

func init() {
	// 初始化命令行参数
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
	rootCmd.PersistentFlags().IntP("log-level", "l", logLevel, "log level")
	rootCmd.PersistentFlags().BoolP("pprof-server", "", false, "enable pprof server")
	rootCmd.Flags().String("server-bind", serverBind, "Server bind addr")
}

func Execute(ctx context.Context) {
	rootCtx = ctx
	if err := rootCmd.Execute(); err != nil {
		logrus.Tracef("cmd execute error: %v", err)
	}
}
