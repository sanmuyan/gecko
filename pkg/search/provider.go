package search

import (
	"gecko/pkg/config"
	"gecko/pkg/model"
	"github.com/sirupsen/logrus"
)

type Provider interface {
	UpdateCode(*model.Project) error
	SearchCode(*model.Project, int, int) (*model.Projects, error)
	DeleteProject(int) error
	Set(string)
}

var Client Provider

func Init() {
	var err error
	switch config.Conf.SearchProvider {
	case "es":
		Client, err = NewEsClient(config.Conf.EsURL)
	}
	if err != nil {
		logrus.Fatalf("create search client error: %v", err)
	}
}
