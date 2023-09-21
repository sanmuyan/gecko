package service

import "gecko/pkg/model"

func storeProject(project *model.Project) {
	Projects.Store(project.ID, &model.Project{
		ID:                project.ID,
		PathWithNamespace: project.PathWithNamespace,
	})
}
