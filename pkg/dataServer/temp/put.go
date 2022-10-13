package temp

import (
	"fmt"
	"net/http"
	"os"

	"github.com/linrds/objectStorage/pkg/dataServer/locate"
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
	err = validateFIleSize(datPath, tempInfo.Size)
	if err != nil {
		os.Remove(datPath)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	commitTempObject(datPath, tempInfo)
}

func validateFIleSize(path string, expected int64) error {
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

func commitTempObject(datFile string, tempInfo *tempInfo) {
	os.Rename(datFile, os.Getenv("STORAGE_ROOT")+"/objects/"+tempInfo.Hash)
	locate.Add(tempInfo.GetHash(), tempInfo.GetId())
}