package util

import (
	"encoding/json"
	"fmt"
	"gecko/pkg/config"
	"gecko/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func ReadFileList(path string, files map[string]string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		subPath, err := os.ReadDir(path)
		if err != nil {
			return err
		}
	Look:
		for _, t := range subPath {
			if t.IsDir() {
				for _, r := range config.Conf.DirBlacklist {
					matchString, err := regexp.MatchString(r, t.Name())
					if err != nil {
						logrus.Errorf("regex error: %s", err.Error())
						continue
					}
					if matchString {
						continue Look
					}
				}
			}
			if err = ReadFileList(fmt.Sprint(path, "/", t.Name()), files); err != nil {
				return err
			}
		}
	} else {
		ext := filepath.Ext(fi.Name())
		for _, regex := range config.Conf.FileBlacklist {
			matchString, err := regexp.MatchString(regex, fi.Name())
			if err != nil {
				logrus.Errorf("regex error: %s", err.Error())
				continue
			}
			if matchString {
				return nil
			}
		}
		if fi.Size() == 0 || fi.Size() > int64(config.Conf.MaxFileSize) {
			return nil
		}
		files[path] = ext
	}
	return nil
}

func BuilderGitURL(repoURL, user, token string) string {
	var httpScheme string
	if strings.Contains(repoURL, "https") {
		httpScheme = "https://"
	} else {
		httpScheme = "http://"
	}
	return fmt.Sprintf("%s%s:%s@%s", httpScheme, user, token, strings.Replace(repoURL, httpScheme, "", 1))
}

func GetPage(c *gin.Context) (int, int) {
	pageNumber := c.Query("page_number")
	pageInt, err := strconv.Atoi(pageNumber)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	pageSize := c.Query("page_size")
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt > 20 || pageSizeInt < 1 {
		pageSizeInt = 10
	}
	return pageInt, pageSizeInt

}

func GetUser(c *gin.Context) (*model.OauthUser, error) {
	userStr, err := json.Marshal(c.Keys["userToken"])
	user := &model.OauthUser{}
	err = json.Unmarshal(userStr, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func MinPreview(content string) string {
	// 只保留前10行, 行不超过1000个字符
	lines := strings.Split(content, "\n")
	if len(lines) > 10 {
		lines = lines[:10]
	}
	for _, line := range lines {
		if len(line) > config.Conf.MaxLineLength {
			return "..."
		}
	}
	return strings.Join(lines, "\n")
}
