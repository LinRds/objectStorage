package rs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/linrds/objectStorage/pkg/utils"
)

type resumableToken struct{
	Name string
	Size int64
	Hash string
	Uuids []string
	Servers []string
}

type RsResumablePutStream struct{
	*RsPutStream
	*resumableToken
}

func NewRsResumablePutStream(hash string, dataServers []string, size int64) (*RsResumablePutStream, error) {
	rsPutStream, err := NewRsPutStream(hash, size, dataServers)
	if err != nil {
		log.Panicln(err)
		return nil, err
	}
	uuids := make([]string, utils.ALL_SHARDS)
	servers := make([]string, utils.ALL_SHARDS)
	for i := range rsPutStream.writers {
		uuids[i] = rsPutStream.writers[i].(*TempPutStream).Uuid
		servers[i] = rsPutStream.writers[i].(*TempPutStream).Server
	}
	resumableToken := &resumableToken{Hash: hash, Size: size, Uuids: uuids, Servers: servers}
	return &RsResumablePutStream{rsPutStream, resumableToken}, nil
}

func NewRsResumablePutStreamFromToken(token string) (*RsResumablePutStream, error) {
	b, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	var streamToken resumableToken
	err = json.Unmarshal(b, &streamToken)
	if err != nil {
		return nil, err
	}
	writers := make([]io.Writer, utils.ALL_SHARDS)
	for i := 0; i < utils.ALL_SHARDS; i++ {
		writers[i] = &TempPutStream{Uuid: streamToken.Uuids[i], Server: streamToken.Servers[i]}
	}
	encoder, err := NewEncoder(writers)
	if err != nil {
		return nil, err
	}
	rsPutStream := &RsPutStream{encoder: encoder}
	return &RsResumablePutStream{rsPutStream, &streamToken}, nil
}

func (r *RsResumablePutStream) ToToken() string{
	b, err := json.Marshal(r.resumableToken)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func (r *RsResumablePutStream) CurrentSize() (int64, error) {
	url := fmt.Sprintf("http://%s/temp/%s", r.Servers[0], r.Uuids[0])
	resp, err := http.Head(url)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	size := utils.GetSizeFromHeader(resp.Header)
	if size == -1 {
		return size, fmt.Errorf("get size from header failed")
	}
	size *= utils.DATA_SHARDS
	if size > r.Size {
		size = r.Size
	}
	return size, nil
}