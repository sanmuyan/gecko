package service

import (
	"gecko/pkg/config"
	"github.com/sirupsen/logrus"
	"os"
)

func (s *Service) Init() {
	// 启动时同步一次
	if _, err := os.Stat(config.Conf.ReposPath); os.IsNotExist(err) {
		logrus.Infof("create repos path: %s", config.Conf.ReposPath)
		err := os.MkdirAll(config.Conf.ReposPath, os.ModePerm)
		if err != nil {
			logrus.Fatalf("create repos path error: %v", err)
		}
	}
	s.GetProjectsAndSync()
}
