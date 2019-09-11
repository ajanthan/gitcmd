package pkg

import (
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const gitDirectory = ".git"
const gitConfigFile = "config"

type Repository struct {
	workTreeDirectory string
	gitDirectoryPath  string
}

var NotARepoError = errors.New("not a git repository")

func NewRepository(path string, force bool) (*Repository, error) {
	repository := &Repository{}
	repository.workTreeDirectory = path
	repository.gitDirectoryPath = path + string(os.PathSeparator) + gitDirectory

	gitDirectory, err := os.Stat(repository.gitDirectoryPath)
	if !force && (err != nil || !gitDirectory.IsDir()) {
		return repository, fmt.Errorf("%s is not a git repository", repository.workTreeDirectory)
	}

	gitConfigFile, err := os.Stat(GetGitConfigFile(repository.gitDirectoryPath))
	if (err != nil && os.IsNotExist(err)) || gitConfigFile.IsDir() {
		if !force {
			return repository, fmt.Errorf("git config file is missing")
		}
	}

	if !force {
		cf, err := ini.Load(gitConfigFile)

		if err != nil {
			return repository, fmt.Errorf("error in loading git config file:%#v", err)
		}
		version, err := cf.Section("core").Key("repositoryformatversion").Int()
		if err != nil {
			return repository, fmt.Errorf("malformed version:%#v", err)
		}
		if version != 0 {
			return repository, fmt.Errorf("unsupported repository version:%d", version)
		}
	}
	return repository, nil
}

func (repo Repository) Create() error {
	//worktree path should
	// 1. Exist
	// 2. Be a directory
	// 3. Be empty
	workTreeFile, err := os.Stat(repo.workTreeDirectory)
	if err != nil && os.IsNotExist(err) {
		os.Mkdir(repo.workTreeDirectory, os.ModePerm)
	} else if !workTreeFile.IsDir() {
		return fmt.Errorf("%s is not a directory ", repo.workTreeDirectory)
	}
	isEmptyDirs, err := IsEmptyDirectory(repo.workTreeDirectory)
	if err != nil || !isEmptyDirs {
		return fmt.Errorf("%s is not an empty directory ", repo.workTreeDirectory)
	}

	//Generate default directories
	// 1. branches
	// 2. objects
	// 3. refs/tags
	// 4. refs/heads
	err = os.MkdirAll(JoinDir(repo.gitDirectoryPath, "branches"), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Mkdir(JoinDir(repo.gitDirectoryPath, "objects"), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(JoinDir(JoinDir(repo.gitDirectoryPath, "refs"), "tags"), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(JoinDir(JoinDir(repo.gitDirectoryPath, "refs"), "heads"), os.ModePerm)
	if err != nil {
		return err
	}
	// write .git/description file
	defaultDescription := []byte("Unnamed repository; edit this file 'description' to name the repository.\n")
	err = ioutil.WriteFile(JoinDir(repo.gitDirectoryPath, "description"), defaultDescription, os.ModePerm)
	if err != nil {
		return err
	}
	// write .git/HEAD file
	defaultHEADRef := []byte("ref: refs/heads/master\n")
	err = ioutil.WriteFile(JoinDir(repo.gitDirectoryPath, "HEAD"), defaultHEADRef, os.ModePerm)
	if err != nil {
		return err
	}
	//Write a default .git/config file
	if err := WriteDefaultGitConfig(repo.gitDirectoryPath); err != nil {
		return err
	}
	return nil
}
func (repo Repository) GetWorkTreeDirectory() string {
	return repo.workTreeDirectory
}
func (repo Repository) GetGitDirectoryPath() string {
	return repo.gitDirectoryPath
}
func FindRepository(workingDirectory string) (string, error) {
	gitDir, err := os.Open(path.Join(workingDirectory, ".git"))
	if strings.Compare(path.Clean("/"), path.Clean(workingDirectory)) == 0 {
		return "", NotARepoError
	}
	if err != nil && os.IsNotExist(err) {
		return FindRepository(path.Clean(path.Join(workingDirectory, "..")))
	}
	fileInfo, err := gitDir.Stat()
	if err != nil {
		return "", err
	}
	if !fileInfo.IsDir() {
		return FindRepository(path.Clean(path.Join(workingDirectory, "..")))
	}
	return workingDirectory, nil
}
