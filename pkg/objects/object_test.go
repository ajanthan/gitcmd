package objects

import (
	"github.com/ajanthan/gitcmd/pkg"
	"os"
	"path"
	"strings"
	"testing"
)

const testText = "This is a test\n"

func setupTestCase(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	err = os.Rename(path.Join(wd, "test", "git"), path.Join(wd, "test", ".git"))
	if err != nil {
		t.Fatal(err)
	}
}
func tearDownTestCase(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	err = os.Rename(path.Join(wd, "test", ".git"), path.Join(wd, "test", "git"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecodeAndEncodeGitObject(t *testing.T) {
	setupTestCase(t)
	defer tearDownTestCase(t)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	gitObject, err := DecodeGitObject("0527e6bd2d76b45e2933183f1b506c7ac49f5872", wd+"/test")
	if err != nil {
		t.Fatal(err)
	}

	if strings.Compare(testText, string(gitObject.Data)) != 0 {
		t.Error("The content of the README.md didn't match")
		t.Fail()
	}
	readmeFileGitObject := GitObject{Type: "blob", Data: []byte(testText)}
	err = os.MkdirAll(pkg.JoinDir("/tmp", "TestDecodeAndEncodeGitObject"), os.ModePerm)
	defer os.RemoveAll(pkg.JoinDir("/tmp", "TestDecodeAndEncodeGitObject"))
	if err != nil {
		t.Fatal(err)
	}
	err = EncodeGitObject(readmeFileGitObject, pkg.JoinDir("/tmp", "TestDecodeAndEncodeGitObject"))
	if err != nil {
		t.Fatal(err)
	}

	gitObject1, err := DecodeGitObject("0527e6bd2d76b45e2933183f1b506c7ac49f5872", pkg.JoinDir("/tmp", "TestDecodeAndEncodeGitObject"))
	if err != nil {
		t.Fatal(err)
	}

	if strings.Compare(testText, string(gitObject1.Data)) != 0 {
		t.Error("The content of the README.md didn't match")
		t.Fail()
	}
}
