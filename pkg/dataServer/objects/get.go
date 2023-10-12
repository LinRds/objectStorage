package objects

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/linrds/objectStorage/pkg/utils"
	"github.com/linrds/objectStorage/pkg/dataServer/locate"
)

func writeFile(w io.Writer, file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}

func getHashValidatedFile(name string) string {
	files, _ := filepath.Glob(os.Getenv("STROAGE_ROOT") + "/objects/" + name + ".*")
	if len(files) != 1 {
		return ""
	}
	file := files[0]
	h := sha256.New()
	writeFile(h, file)
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	hash := strings.Split(name, ".")[2]
	if d != hash {
		log.Println("object hash mismatch, remove", file)
		os.Remove(file)
		locate.Del(hash)
		return ""
	}
	return file
}


func get(w http.ResponseWriter, r *http.Request) {
	hash := utils.GetNameFromUrl(r.URL)
	file := getHashValidatedFile(hash)
	if file == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	writeFile(w, file)
}