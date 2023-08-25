package service

import (
	"bytes"
	"fmt"
	"gecko/pkg/config"
	"gecko/pkg/model"
	"gecko/pkg/util"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

func (s *Service) GitClone(project *model.Project) {
	fullPath := fmt.Sprint(config.Conf.ReposPath, "/", project.Name, project.ID)
	clonePath := fmt.Sprint(project.Name, project.ID)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", util.BuilderGitURL(project.URL, config.Conf.GitlabUser, config.Conf.GitlabToken), clonePath)
		cmd.Dir = config.Conf.ReposPath
		logrus.Infof("clone project: %s", project.NamespacePath)
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		err := cmd.Run()
		if err != nil {
			logrus.Errorf("clone project %s error: %s", project.NamespacePath, stderr)
			return
		}
		logrus.Infof("clone project %s message: %s", project.NamespacePath, stdout)
		s.ParseFile(project, fullPath)
	}
}

func (s *Service) GitPull(project *model.Project) {
	fullPath := fmt.Sprint(config.Conf.ReposPath, "/", project.Name, project.ID)
	cmd := exec.Command("git", "pull")
	cmd.Dir = fullPath
	logrus.Infof("pull project: %s", project.NamespacePath)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		logrus.Errorf("pull project %s error: %s", project.NamespacePath, stderr)
		return
	}
	logrus.Infof("pull project %s message: %s", project.NamespacePath, stdout)
	if strings.HasPrefix(stdout.String(), "Already up-to-date") || strings.HasPrefix(stdout.String(), "Already up to date") {
		return
	}
	s.ParseFile(project, fullPath)
}
