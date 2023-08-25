package cmd

import (
	"gecko/pkg/config"
	"gecko/pkg/search"
	"gecko/server/controller"
	"gecko/server/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
	"runtime"
)

var cmdReady bool

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "Gecko Server",
	Run: func(cmd *cobra.Command, args []string) {
		cmdReady = true
	},
	Example: "gecko -c config.yaml",
}

var configFile string

const (
	logLevel         = 4
	serverBind       = "0.0.0.0:8080"
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
	defaultConfig := "./config/config.yaml"
	if os.Getenv("ENV_NAME") != "" {
		defaultConfig = "./config/config-" + os.Getenv("ENV_NAME") + ".yaml"
	}
	rootCmd.Flags().StringVarP(&configFile, "config", "c", defaultConfig, "config file")
	rootCmd.Flags().IntP("log-level", "l", logLevel, "log level")
	rootCmd.Flags().String("server-bind", serverBind, "server bind addr")
}

func initConfig() error {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File)
			return frame.Function, fileName
		},
	})

	if len(configFile) > 0 {
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()
		if err != nil {
			return err
		}
	}
	_ = viper.BindPFlag("log_level", rootCmd.Flags().Lookup("log-level"))
	_ = viper.BindPFlag("server_bind", rootCmd.Flags().Lookup("server-bind"))

	config.Conf.SyncProjectLimit = syncProjectLimit
	config.Conf.SearchProvider = searchProvider
	config.Conf.EsURL = esUrl
	config.Conf.ReposPath = reposPath
	config.Conf.GitlabURL = gitlabUrl
	config.Conf.HTTPHost = httpHost
	config.Conf.MaxFileSize = maxFileSize
	config.Conf.MaxLineLength = maxLienLength
	config.Conf.MaxSearchTotal = maxSearchTotal

	err := viper.Unmarshal(&config.Conf)
	if err != nil {
		return err
	}

	logrus.SetLevel(logrus.Level(config.Conf.LogLevel))
	gin.SetMode(gin.ReleaseMode)
	if logrus.Level(config.Conf.LogLevel) >= logrus.DebugLevel {
		gin.SetMode(gin.DebugMode)
		logrus.SetReportCaller(true)
	}

	config.Conf.DirBlacklist = append(config.Conf.DirBlacklist, builtInDirBlacklist...)
	config.Conf.FileBlacklist = append(config.Conf.FileBlacklist, builtInFileBlacklist...)
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
	if cmdReady {
		err := initConfig()
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Debugf("config %+v", config.Conf)

		initConfigPost()
	}
}

func initConfigPost() {
	search.Init()
	go service.NewService().Init()
	controller.RunServer(config.Conf.ServerBind)
}
