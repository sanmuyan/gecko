package model

import "errors"

type Namespace struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type Project struct {
	ID             int       `json:"id" bson:"id"`
	Name           string    `json:"name" bson:"name"`
	NamespacePath  string    `json:"path_with_namespace" bson:"namespace_path"`
	Visibility     string    `json:"visibility" bson:"visibility"`
	URL            string    `json:"http_url_to_repo" bson:"project_url"`
	Namespace      Namespace `json:"namespace" bson:"namespace"`
	CodeFileName   string    `json:"code_file_name,omitempty" bson:"code_file_name"`
	CodeSuffixName string    `json:"code_suffix_name,omitempty" bson:"code_suffix_name"`
	CodeContent    string    `json:"code_content,omitempty" bson:"code_content"`
}

type Projects struct {
	Projects   []*Project `json:"projects"`
	TotalCount int64      `json:"total_count"`
	PageNumber int        `json:"page_number"`
	PageSize   int        `json:"page_size"`
}

type ProjectACL struct {
	AccessLevel int `json:"access_level"`
}

type GitlabWebhook struct {
	EventName string `json:"event_name" binding:"required"`
	ProjectID int    `json:"project_id,omitempty"`
}

type OauthUser struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	IsAdmin        bool   `json:"is_admin"`
	ExpirationTime int64  `json:"expiration_time"`
}

func (o OauthUser) Valid() error {
	err := errors.New("required is nil")
	if o.ID == 0 {
		return err
	}
	if o.Username == "" {
		return err
	}
	return nil
}
