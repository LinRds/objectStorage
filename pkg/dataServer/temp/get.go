package temp

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/linrds/objectStorage/pkg/utils"
)

func get(w http.ResponseWriter, r *http.Request) {
	uuid := utils.GetNameFromUrl(r.URL)
	f, e := os.Open(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid + ".dat")
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}