package util

import (
	"gecko/pkg/config"
	"testing"
)

func TestReadFileList(t *testing.T) {
	config.Conf.MaxFileSize = 1000000000
	config.Conf.DirBlacklist = []string{"^\\.", "node_modules"}
	config.Conf.FileBlacklist = []string{"\\.exe$"}
	files := make(map[string]string)
	err := ReadFileList("./", files)
	if err != nil {
		t.Error(err)
	}
	t.Log(len(files))
	for k, v := range files {
		t.Log(k, v)
	}
}
