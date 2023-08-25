package service

import (
	"gecko/pkg/model"
	"gecko/pkg/search"
)

func (s *Service) SearchCode(project *model.Project, pageNumber, pageSize int) (*model.Projects, error) {
	codeProjects, err := search.Client.SearchCode(project, pageNumber, pageSize)
	if err != nil {
		return nil, err
	}
	return codeProjects, nil
}

func (s *Service) Projects() *model.Projects {
	var projects []*model.Project
	Projects.Range(func(key, value any) bool {
		if project, ok := value.(*model.Project); ok {
			projects = append(projects, project)
		}
		return true
	})
	return &model.Projects{
		Projects:   projects,
		TotalCount: int64(len(projects)),
		PageSize:   len(projects),
		PageNumber: 1,
	}
}
