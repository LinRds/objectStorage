package temp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/linrds/objectStorage/pkg/utils"
)

func post(w http.ResponseWriter, r *http.Request) {
	hash := utils.GetNameFromUrl(r.URL)
	size := utils.GetSizeFromHeader(r.Header)
	if size < 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cmd := exec.Command("uuidgen", "-t")
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	uuid := strings.TrimSuffix(string(output), "\n")
	datPath := fmt.Sprintf("%s/temp/%s.dat", os.Getenv("STORAGE_ROOT"), uuid)
	if f, err := os.Create(datPath); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		f.Close()
	}
	info := &tempInfo{Uuid: string(uuid), Hash: hash, Size: size}
	if err = saveTempInfo(info); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(uuid))
}

func saveTempInfo(info *tempInfo) error {
	b, err := json.Marshal(info)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/temp/%s", os.Getenv("STORAGE_ROOT"), info.Uuid)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.Write(b); err != nil {
		return err
	}
	return nil
}