package objects

import (
	"compress/gzip"
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

func writeCompactedFile(w io.Writer, file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		log.Println(err)
		return
	}
	defer gr.Close()
	io.Copy(w, gr)
}

func getHashValidatedFile(name string) string {
	files, _ := filepath.Glob(os.Getenv("STROAGE_ROOT") + "/objects/" + name + ".*")
	if len(files) != 1 {
		return ""
	}
	file := files[0]
	h := sha256.New()
	writeCompactedFile(h, file)
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
	writeCompactedFile(w, file)
}