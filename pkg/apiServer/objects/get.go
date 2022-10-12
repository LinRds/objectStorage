package objects

import (
	"io"
	"log"
	"net/http"
	"strconv"

	es "github.com/linrds/objectStorage/pkg/elastic"
	"github.com/linrds/objectStorage/pkg/utils"
)

func get(w http.ResponseWriter, r *http.Request) {
	name := utils.GetNameFromUrl(r.URL)
	versionId := r.URL.Query()["version"]
	version := 0
	var err error
	if len(versionId) != 0 {
		version, err = strconv.Atoi(versionId[0])
		if err != nil {		//有问题就报错
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	meta, err := es.GetMetadata(name, version)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if meta.Hash == "" {			
		w.WriteHeader(http.StatusNotFound) //空的就是没找到，删除操作也是将hash置空
		return
	}
	stream, err := GetStream(meta.Hash, meta.Size)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(w, stream)
}