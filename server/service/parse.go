package service

import (
	"gecko/pkg/model"
	"gecko/pkg/search"
	"gecko/pkg/util"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"unicode/utf8"
)

// ParseFile
// 1. 解析一个仓库的所有代码文件
// 2. 过滤掉不符合的文件
// 3. 读取代码内容存储到搜索引擎
func (s *Service) ParseFile(project *model.Project, projectPath string) {
	// 项目有更新直接删除，在全量更新
	err := search.Client.DeleteProject(project.ID)
	if err != nil {
		logrus.Errorf("delete project error: %v", err)
	}
	files := make(map[string]string)
	err = util.ReadFileList(projectPath, files)
	if err != nil {
		logrus.Errorf("read file list error: %v", err)
	}
	for file, suffixName := range files {
		projectFile := strings.Replace(file, projectPath+"/", "", 1)
		logrus.Debugf("update code: %s %s", project.Name, projectFile)
		readFile, err := os.ReadFile(file)
		if err != nil {
			logrus.Errorf("read file error: %v", err)
			continue
		}
		if !utf8.Valid(readFile) {
			logrus.Debugf("file is not text file: %s", file)
			continue
		}
		project.CodeFileName = projectFile
		project.CodeSuffixName = suffixName
		if project.CodeSuffixName == "" {
			project.CodeSuffixName = projectFile
		}
		project.CodeContent = string(readFile)
		err = search.Client.UpdateCode(project)
		if err != nil {
			logrus.Errorf("update code error: %v", err)
			continue
		}
	}
}
