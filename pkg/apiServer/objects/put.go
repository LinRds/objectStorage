package objects

import (
	"log"
	"net/http"
	"strings"

	"github.com/linrds/objectStorage/pkg/utils"
	es "github.com/linrds/objectStorage/pkg/elastic"
)

func put(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	hash := utils.GetHashFromHeader(r.Header)
	size := utils.GetSizeFromHeader(r.Header)
	if len(hash) == 0 {
		log.Println("Hash is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	statusCode, err := StoreObject(r.Body, hash, size)
	if err != nil {
		log.Println(err)
		w.WriteHeader(statusCode)
		return
	}
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
		return
	}
	err = es.AddVersion(name, hash, size)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}