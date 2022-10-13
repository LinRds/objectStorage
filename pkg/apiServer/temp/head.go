package temp

import (
	"fmt"
	"log"
	"net/http"

	"github.com/linrds/objectStorage/pkg/rs"
	"github.com/linrds/objectStorage/pkg/utils"
)

func head(w http.ResponseWriter, r *http.Request) {
	token := utils.GetNameFromUrl(r.URL)
	stream, err := rs.NewRsResumablePutStreamFromToken(token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	size, err := stream.CurrentSize()
	if size == -1 || err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("content-length", fmt.Sprintf("%d", size))
}