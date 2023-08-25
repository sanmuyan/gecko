package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"gecko/pkg/config"
	"gecko/pkg/model"
	"github.com/sanmuyan/dao/request"
	"strconv"
)

func GetProjectsTotal() (int, error) {
	reqConfig := request.Request{
		Config: &request.Options{
			URL:    fmt.Sprint(config.Conf.GitlabURL, "/api/v4/projects?private_token=", config.Conf.GitlabToken, "&per_page=1"),
			Method: "GET",
		},
	}
	res, err := reqConfig.Request()
	if err != nil {
		return 0, err
	}
	if total, ok := res.Header["X-Total"]; ok {
		totalInt, _ := strconv.Atoi(total[0])
		return totalInt, nil
	}
	return 0, nil
}

func GetProjects(page, pageSize int) ([]*model.Project, error) {
	reqConfig := request.Request{
		Config: &request.Options{
			URL:    fmt.Sprint(config.Conf.GitlabURL, "/api/v4/projects?private_token=", config.Conf.GitlabToken, "&per_page=", pageSize, "&page=", page),
			Method: "GET",
		},
	}
	res, err := reqConfig.Request()
	if err != nil {
		return nil, err
	}
	var projects []*model.Project
	err = json.Unmarshal(res.Body, &projects)
	if err != nil {
		return nil, err
	}
	if len(projects) == 0 {
		return nil, errors.New("projects not found")
	}
	return projects, nil
}

func GetProject(projectID int) (*model.Project, error) {
	reqConfig := request.Request{
		Config: &request.Options{
			URL:    fmt.Sprint(config.Conf.GitlabURL, "/api/v4/projects/", projectID, "/?private_token=", config.Conf.GitlabToken),
			Method: "GET",
		},
	}
	res, err := reqConfig.Request()
	if err != nil {
		return nil, err
	}
	var project *model.Project
	err = json.Unmarshal(res.Body, &project)
	if err != nil {
		return nil, err
	}
	if project.ID == 0 {
		return nil, errors.New("project not found")
	}
	return project, nil
}
