package service

import (
	"fmt"
	"gecko/pkg/config"
	"gecko/pkg/model"
	"gecko/pkg/request"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

// 1. 获取所用项目
// 2. 判断在项目的本地仓库是否存在，使用pull/clone 同步代码
// 3. 解析仓库代码

var Projects = new(sync.Map)

func (s *Service) GetProjectsAndSync() {
	projectTotal, err := request.GetProjectsTotal()
	if err != nil {
		logrus.Errorf("get project total error: %v", err)
		time.Sleep(3 * time.Second)
		s.GetProjectsAndSync()
		return
	}
	logrus.Infof("project total: %d", projectTotal)
	pageSize := config.Conf.SyncProjectLimit
	pool := make(chan bool, config.Conf.SyncProjectLimit)

	startTime := time.Now()
	wg := new(sync.WaitGroup)
	for page := 1; page*pageSize < projectTotal+pageSize; page++ {
		projects, err := request.GetProjects(page, pageSize)
		if err != nil {
			logrus.Errorf("get gitlab projects error: %v", err)
		}
		for _, project := range projects {
			wg.Add(1)
			Projects.Store(project.ID, project)
			go func(project *model.Project) {
				pool <- true
				defer func() {
					<-pool
					wg.Done()
				}()
				s.SyncProject(project)
			}(project)
		}
	}
	wg.Wait()
	endTime := time.Now()
	logrus.Infof("%d project sync completed", projectTotal)
	logrus.Infof("project sync time: %v", endTime.Sub(startTime))
}

func (s *Service) GetProject(projectID int) (*model.Project, error) {
	project, err := request.GetProject(projectID)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s *Service) SyncProject(project *model.Project) {
	Projects.Store(project.ID, project)
	fullPath := fmt.Sprint(config.Conf.ReposPath, "/", project.Name, project.ID)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		s.GitClone(project)
	} else {
		s.GitPull(project)
	}
}
