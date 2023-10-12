package temp

import (
	"log"
	"net/http"
	"os"
	"fmt"

	"github.com/linrds/objectStorage/pkg/utils"
)

func head(w http.ResponseWriter, r *http.Request) {
	uuid := utils.GetNameFromUrl(r.URL)
	f, err := os.Open(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid + ".dat")
	if err != nil {
		log.Println("[dataServer]: ", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	stat, _ := f.Stat()
	w.Header().Set("size", fmt.Sprintf("%d", stat.Size()))
}