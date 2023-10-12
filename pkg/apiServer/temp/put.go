package temp

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/linrds/objectStorage/pkg/apiServer/locate"
	es "github.com/linrds/objectStorage/pkg/elastic"
	"github.com/linrds/objectStorage/pkg/rs"
	"github.com/linrds/objectStorage/pkg/utils"
)

func put(w http.ResponseWriter, r *http.Request) {
	token := utils.GetNameFromUrl(r.URL)
	log.Println(token)
	stream, err := rs.NewRsResumablePutStreamFromToken(token)
	if err != nil {
		log.Printf("retrive stream from token failed, %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	offset, err := utils.GetOffsetFromHeader(r.Header)
	if offset == -1 {
		log.Printf("invalid offset, %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	current, err := stream.CurrentSize()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if offset != current {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	buf := make([]byte, utils.BLOCK_SIZE)
	for {
		n, err := io.ReadFull(r.Body, buf)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			stream.Commit(false)
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		current += int64(n)
		if current > stream.Size {
			stream.Commit(false)
			log.Println("resumable put exceeds file size")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if n != utils.BLOCK_SIZE && current != stream.Size {
			return // 每次写入一个block的字节，不足一个block的数据会被直接丢弃
		}
		stream.Write(buf[:n])
		if (current == stream.Size) {
			getStream, err := rs.NewRsResumableGetStream(stream.Servers, stream.Uuids, stream.Size)
			if err != nil {
				log.Println("can not validate hash")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			hash := url.PathEscape(utils.CalculateHash(getStream))
			log.Println("hash: ", hash, "expected: ", stream.Hash)
			if hash != stream.Hash {
				stream.Commit(false)
				log.Println("get file successful but hash mismatch")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if locate.Exist(url.PathEscape(hash)) {
				stream.Commit(false)
			} else {
				stream.Commit(true)
			}
			err = es.AddVersion(stream.Name, stream.Hash, stream.Size)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}
}