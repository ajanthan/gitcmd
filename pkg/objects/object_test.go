package objects

import (
	"github.com/ajanthan/gitcmd/pkg"
	"os"
	"strings"
	"testing"
)

const testText = "This is a test\n"

func TestDecodeAndEncodeGitObject(t *testing.T) {
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
	}
}
