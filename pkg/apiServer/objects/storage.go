package objects

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/linrds/objectStorage/pkg/apiServer/heartbeat"
	"github.com/linrds/objectStorage/pkg/apiServer/locate"
	"github.com/linrds/objectStorage/pkg/rs"
	"github.com/linrds/objectStorage/pkg/utils"
)

func putStream(hash string, size int64) (*rs.RsPutStream, error) {
	dataServers := heartbeat.ChooseRandomDataServers(utils.ALL_SHARDS, nil)
	if len(dataServers) < utils.ALL_SHARDS {
		err := fmt.Errorf("can not find enough data servers, find %d need %d", utils.ALL_SHARDS, len(dataServers))
		return nil, err
	}
	
	return rs.NewRsPutStream(hash, size, dataServers)
}

func GetStream(hash string, size int64) (*rs.RsGetStream, error) {
	dataServers := locate.Locate(hash)
	lack := utils.ALL_SHARDS - len(dataServers)
	var comple []string
	if lack > 0 {
		comple = heartbeat.ChooseRandomDataServers(lack, dataServers)
	}
	return rs.NewRsGetStream(hash, size, dataServers, comple)
}

func StoreObject(r io.Reader, hash string, size int64) (int, error){
	if locate.Exist(hash) {
		return http.StatusOK, nil
	}
	stream, err := putStream(url.PathEscape(hash), size)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	reader := io.TeeReader(r, stream)
	if hash != utils.CalculateHash(reader) {
		stream.Commit(false)
		err = fmt.Errorf("hash mismatch")
		return http.StatusInternalServerError, err
	}
	stream.Commit(true)
	return http.StatusOK, nil
}

