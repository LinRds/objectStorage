package objects

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"fmt"

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
	offset, err := utils.GetOffsetFromHeader(r.Header)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	stream, err := GetStream(meta.Hash, meta.Size)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer stream.Close() // 将恢复分片时重新写入的临时分片转正
	if offset > 0 {
		stream.Seek(offset, io.SeekCurrent)
		w.Header().Set("content-range", fmt.Sprintf("bytes %d-%d/%d", offset, meta.Size-1, meta.Size))
		w.WriteHeader(http.StatusPartialContent)
	}
	_, err = io.Copy(w, stream)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
	}
}