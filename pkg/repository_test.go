package pkg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type FindRepositoryTestCase struct {
	WorkingDirectory string
	GitDirectory     string
	Expected         bool
}

var testDirectories = []string{
	"/dir2/dir2/dir3",
	"/dir2",
	"/dir2/.git",
	"/dir2/dir4/dir5",
	"/dir6/dir7/dir8",
	"/dir9/dir10/dir11/dir8",
	"/dir9/dir10/dir11/dir8/.git",
	"/dir9/dir10/dir11/dir8/dir12/dir13/dir14/dir15/dir16",
}

func setupTest(t *testing.T, testDirectories ...string) {
	for _, testDirectory := range testDirectories {
		err := os.MkdirAll(filepath.FromSlash("/tmp/RepositoryTest"+testDirectory), os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
	}

}

func tearDown(t *testing.T) {
	err := os.RemoveAll(filepath.FromSlash("/tmp/RepositoryTest/"))
	if err != nil {
		t.Fatal(err)
	}
}

func assertRepository(t *testing.T, workingDirectory, gitDirectory string) (bool, error) {
	t.Helper()
	err := os.Chdir(filepath.FromSlash("/tmp/RepositoryTest" + workingDirectory))
	if err != nil {
		return false, err
	}
	gitDir, err := FindRepository(filepath.FromSlash("/tmp/RepositoryTest" + workingDirectory))
	if err != nil {
		return false, err
	}
	if strings.Compare(gitDir, filepath.FromSlash("/tmp/RepositoryTest"+gitDirectory)) == 0 {
		return true, nil
	}
	return false, nil
}

func TestFindRepository(t *testing.T) {
	setupTest(t, testDirectories...)
	defer tearDown(t)
	testCases := map[string]FindRepositoryTestCase{
		"1": {WorkingDirectory: testDirectories[0],
			GitDirectory: testDirectories[1],
			Expected:     true},
		"2": {WorkingDirectory: testDirectories[3],
			GitDirectory: testDirectories[1],
			Expected:     true},
		"3": {WorkingDirectory: testDirectories[4],
			GitDirectory: "",
			Expected:     false},
		"4": {WorkingDirectory: testDirectories[7],
			GitDirectory: testDirectories[5],
			Expected:     true},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := assertRepository(t, testCase.WorkingDirectory, testCase.GitDirectory)
			if err != nil {
				if testCase.Expected || err != NotARepoError {
					t.Error(err)
				}
			}
			if result != testCase.Expected {
				t.Fail()
			}
		})
	}
}
