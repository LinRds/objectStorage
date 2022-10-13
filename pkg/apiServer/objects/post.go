package objects

import (
	"log"
	"net/http"
	"net/url"

	"github.com/linrds/objectStorage/pkg/apiServer/heartbeat"
	"github.com/linrds/objectStorage/pkg/apiServer/locate"
	es "github.com/linrds/objectStorage/pkg/elastic"
	"github.com/linrds/objectStorage/pkg/rs"
	"github.com/linrds/objectStorage/pkg/utils"
)

func post(w http.ResponseWriter, r *http.Request) {
	name := utils.GetNameFromUrl(r.URL)
	size := utils.GetSizeFromHeader(r.Header)
	if size == -1 {
		log.Println("invalid size of -1")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	if locate.Exist(hash) {
		err := es.AddVersion(name, hash, size)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	}

	dataServers := heartbeat.ChooseRandomDataServers(utils.ALL_SHARDS, nil)
	if len(dataServers) != utils.ALL_SHARDS {
		log.Println("can not find enough data servers")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	stream, err := rs.NewRsResumablePutStream(url.PathEscape(hash), dataServers, size)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	stream.Name = name

	w.Header().Set("location", "/temp/" + url.PathEscape(stream.ToToken()))
	w.WriteHeader(http.StatusCreated)
}