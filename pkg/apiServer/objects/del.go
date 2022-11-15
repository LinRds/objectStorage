package objects

import (
	"log"
	"net/http"
	"strings"

	es "github.com/linrds/objectStorage/pkg/elastic"
)

func del(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	version, e := es.SearchLatestVersion(name)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	e = es.PutMetadata(name, version.Version+1, 0, "")
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}