package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/ajanthan/gitcmd/pkg"
	"io"
	"os"
	"strconv"
)

var nullBytes = []byte("\x00")

type Serializer interface {
}
type GitObject struct {
	Data []byte
	Type string
}

func (gitObject GitObject) ReadObject(sha string) error {

	return nil
}

func (gitObject GitObject) WriteObject(object []byte) error {
	return nil
}

//zlip(<object type> <size(body)>.<body>)
func DecodeGitObject(sha string, gitDirectoryPath string) (GitObject, error) {
	gitObject := GitObject{}
	objectPath := pkg.JoinDir(gitDirectoryPath, ".git", "objects", sha[:2], sha[2:])
	objectFile, err := os.Open(objectPath)
	if err != nil {
		return gitObject, err
	}

	if err != nil {
		return gitObject, err
	}

	var dataBuf bytes.Buffer
	gzReader, err := zlib.NewReader(objectFile)
	if err != nil {
		return gitObject, err
	}

	_, err = io.Copy(&dataBuf, gzReader)
	if err != nil {
		return gitObject, err
	}
	data := dataBuf.Bytes()
	i := bytes.IndexByte(data, ' ')
	if i < 0 {
		return gitObject, fmt.Errorf("invalid object format")
	}
	objectType := string(data[:i])

	j := bytes.Index(data, nullBytes)
	if j < 0 {
		return gitObject, fmt.Errorf("invalid object format")
	}
	sizeData := string(data[i+1 : j])
	s, err := strconv.Atoi(sizeData)
	if err != nil {
		return gitObject, err
	}
	if s != len(data[j+1:]) {
		return gitObject, fmt.Errorf("invalid object format")
	}
	gitObject.Type = objectType
	gitObject.Data = data[j+1:]
	return gitObject, nil
}

func EncodeGitObject(gitObject GitObject, gitDirectoryPath string) error {
	var data []byte
	data = append(data, []byte(gitObject.Type)...)
	data = append(data, ' ')
	size := len(gitObject.Data)
	sizeBytes := strconv.Itoa(size)
	data = append(data, sizeBytes...)
	data = append(data, nullBytes...)
	data = append(data, gitObject.Data...)

	sha1Cal := sha1.New()
	_, err := sha1Cal.Write(data)
	if err != nil {
		return err
	}
	objectName := hex.EncodeToString(sha1Cal.Sum(nil))
	objectFilePath := pkg.JoinDir(gitDirectoryPath, ".git", "objects", string(objectName[:2]))
	err = os.MkdirAll(objectFilePath, os.ModePerm)
	if err != nil {
		return err
	}
	objectFile, err := os.Create(pkg.JoinDir(objectFilePath, string(objectName[2:])))
	if err != nil {
		return err
	}
	defer objectFile.Close()

	zlibWriter := zlib.NewWriter(objectFile)

	_, err = zlibWriter.Write(data)
	if err != nil {
		return err
	}
	zlibWriter.Close()
	return nil
}
