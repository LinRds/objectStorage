package temp

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// TODO error作为日志输出
func patch(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2] // uuid
	infoPath := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	datPath := infoPath + ".dat"
	tempInfo, err := readTempInfo(infoPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = writeDatFile(datPath, r.Body, tempInfo.Size)
	if err != nil{
		os.Remove(infoPath)
		os.Remove(datPath)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func readTempInfo(path string) (*tempInfo, error) {
	infoFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer infoFile.Close()
	content, err := ioutil.ReadAll(infoFile)
	if err != nil {
		return nil, err
	}
	var info *tempInfo
	json.Unmarshal(content, info)
	return info, nil
}

func writeDatFile(path string, body io.ReadCloser, size int64) error {
	datFile, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer datFile.Close()
	_, err = io.Copy(datFile, body)
	if err != nil {
		return err
	}
	info, err := datFile.Stat()
	if err != nil {
		return err
	}
	if actual := info.Size(); actual > size {
		return fmt.Errorf("actual size %d exceeds %d", actual, size)
	}
	return nil
}