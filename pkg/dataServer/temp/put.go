package temp

import (
	"fmt"
	"net/http"
	"os"

	"github.com/linrds/objectStorage/pkg/utils"
)

func put(w http.ResponseWriter, r *http.Request) {
	name := utils.GetNameFromUrl(r.URL)
	infoPath := os.Getenv("STORAGE_ROOT") + "/temp/" + name
	datPath := infoPath + ".dat"
	tempInfo, err := readTempInfo(infoPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	os.Remove(infoPath)
	err = validateFileSize(datPath, tempInfo.Size)
	if err != nil {
		os.Remove(datPath)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	commitTempObject(datPath, tempInfo)
}

func validateFileSize(path string, expected int64) error {
	datFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer datFile.Close()
	datInfo, err := datFile.Stat()
	if err != nil {
		return err
	}
	if actual := datInfo.Size(); actual != expected {
		return fmt.Errorf("actual size %d not equal to expected size %d", actual, expected)
	}
	return nil
}