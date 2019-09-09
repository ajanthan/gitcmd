package pkg

import (
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
)

func JoinDir(parent ...string) string {
	var finalPath string
	for _, subDir := range parent {
		finalPath = finalPath + string(os.PathSeparator) + subDir
	}
	return finalPath
}

func IsEmptyDirectory(directoryName string) (bool, error) {
	subDir, err := ioutil.ReadDir(directoryName)
	if err != nil {
		return false, err
	}
	if len(subDir) == 0 {
		return true, nil
	}
	return false, nil
}

func GetGitConfigFile(gitDir string) string {
	return JoinDir(gitDir, gitConfigFile)
}

func WriteDefaultGitConfig(baseRepository string) error {
	cf := ini.Empty()
	sc := cf.Section("core")
	_, err := sc.NewKey("repositoryformatversion", "0")
	if err != nil {
		return err
	}
	_, err = sc.NewKey("filemode", "false")
	if err != nil {
		return err
	}
	_, err = sc.NewKey("bare", "false")
	if err != nil {
		return err
	}
	err = cf.SaveTo(JoinDir(baseRepository, "config"))
	if err != nil {
		return err
	}
	return nil
}
