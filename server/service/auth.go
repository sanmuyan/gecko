package service

import (
	"fmt"
	"gecko/pkg/config"
	"gecko/pkg/model"
	"gecko/pkg/request"
	"github.com/sanmuyan/xpkg/xjwt"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"time"
)

func (s *Service) IsUserAccess(project *model.Project, userID int) bool {
	if project.Namespace.Kind == "group" {
		acl, err := request.GetGroupACL(project.Namespace.ID, userID)
		if err != nil {
			return false
		}
		if acl.AccessLevel > 10 {
			return true
		}
	}
	acl, err := request.GetProjectACL(project.ID, userID)
	if err != nil {
		return false
	}
	if acl.AccessLevel > 10 {
		return true
	}
	return false
}

func (s *Service) OauthConfig() *oauth2.Config {
	scopes := []string{"read_user"}
	redirectURL := fmt.Sprint(config.Conf.HTTPHost, "/oauth/callback")
	return &oauth2.Config{
		ClientID:     config.Conf.OAuthClientID,
		ClientSecret: config.Conf.OAuthClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprint(config.Conf.GitlabURL, "/oauth/authorize"),
			TokenURL: fmt.Sprint(config.Conf.GitlabURL, "/oauth/token"),
		},
		Scopes:      scopes,
		RedirectURL: redirectURL,
	}
}

func (s *Service) Login() string {
	return s.OauthConfig().AuthCodeURL("state")
}

func (s *Service) OauthCallback(code string) any {
	user, err := request.GetOauthUser(s.OauthConfig(), code)
	if err != nil {
		logrus.Errorf("get oauth user error: %v", err)
	}
	if err = user.Valid(); err != nil {
		return err
	}
	user.ExpirationTime = time.Now().UTC().Add(time.Hour * 24 * 7).Unix()
	token, err := xjwt.CreateToken(user, config.Conf.TokenKey)
	if err != nil {
		logrus.Errorf("create token error: %v", err)
	}
	return token
}
