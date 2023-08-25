package request

import (
	"context"
	"encoding/json"
	"fmt"
	"gecko/pkg/config"
	"gecko/pkg/model"
	"github.com/sanmuyan/dao/request"
	"golang.org/x/oauth2"
	"io"
)

func GetOauthUser(conf *oauth2.Config, code string) (*model.OauthUser, error) {
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	client := conf.Client(context.Background(), token)
	resp, err := client.Get(fmt.Sprint(config.Conf.GitlabURL, "/api/v4/user"))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request error: %v", resp.Status)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	user := &model.OauthUser{}
	err = json.Unmarshal(body, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetProjectACL(projectID, userID int) (*model.ProjectACL, error) {
	reqConfig := request.Request{
		Config: &request.Options{
			URL:    fmt.Sprint(config.Conf.GitlabURL, "/api/v4/projects/", projectID, "/members/", userID, "/?private_token=", config.Conf.GitlabToken),
			Method: "GET",
		},
	}
	res, err := reqConfig.Request()
	if err != nil {
		return nil, err
	}
	var acl *model.ProjectACL
	err = json.Unmarshal(res.Body, &acl)
	if err != nil {
		return nil, err
	}
	return acl, nil
}

func GetGroupACL(groupID, userID int) (*model.ProjectACL, error) {
	reqConfig := request.Request{
		Config: &request.Options{
			URL:    fmt.Sprint(config.Conf.GitlabURL, "/api/v4/groups/", groupID, "/members/", userID, "/?private_token=", config.Conf.GitlabToken),
			Method: "GET",
		},
	}
	res, err := reqConfig.Request()
	if err != nil {
		return nil, err
	}
	var acl *model.ProjectACL
	err = json.Unmarshal(res.Body, &acl)
	if err != nil {
		return nil, err
	}
	return acl, nil
}
