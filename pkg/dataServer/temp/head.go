package temp

import (
	"fmt"
	"net/http"
	"os"

	"github.com/linrds/objectStorage/pkg/utils"
)

func head(w http.ResponseWriter, r *http.Request) {
	uuid := utils.GetNameFromUrl(r.URL)
	f, err := os.Open(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	stat, _ := f.Stat()
	w.Header().Set("content-length", fmt.Sprintf("%d", stat.Size()))
}