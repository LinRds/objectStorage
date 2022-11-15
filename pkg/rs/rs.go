package rs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"errors"

	"github.com/linrds/objectStorage/pkg/utils"
)

var ErrBuildWriteStream = errors.New("build write stream failed")

type RsPutStream struct {
	*encoder
}

type RsGetStream struct {
	*decoder
}

type TempPutStream struct {
	Server string
	Uuid string
}

func NewTempPutStream(name, server string, size int64) (*TempPutStream, error) {
	url := fmt.Sprintf("http://%s/temp/%s", server, url.PathEscape(name))
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, ErrBuildWriteStream
	}
	client := http.Client{}
	request.Header.Set("size", fmt.Sprintf("%d", size))
	resp, err := client.Do(request)
	if err != nil {
		return nil, ErrBuildWriteStream
	}
	defer resp.Body.Close()
	uuid, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &TempPutStream{Server: server, Uuid: string(uuid)}, nil
}

func (t *TempPutStream) Write(b []byte) (n int, err error){
	url := fmt.Sprintf("http://%s/temp/%s", t.Server, t.Uuid)
	reader := bytes.NewReader(b)
	request, err := http.NewRequest("PATCH", url, reader)
	if err != nil {
		return 0, err
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return 0, ErrBuildWriteStream
	}
	defer resp.Body.Close()
	size, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil
	}
	n, err = strconv.Atoi(string(size))
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (t *TempPutStream) Commit(success bool) error {
	url := fmt.Sprintf("http://%s/temp/%s", t.Server, t.Uuid)
	method := "DELETE"
	if success {
		method = "PUT"
	}
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	client := http.Client{}
	_, err = client.Do(request)
	return err
}

func NewRsPutStream(name string, size int64, dataServers []string) (*RsPutStream, error) {
	perShard := (size + utils.DATA_SHARDS - 1) / utils.DATA_SHARDS // 向上取整
	var writers = make([]io.Writer, len(dataServers))
	var err error
	for i := range dataServers {
		writers[i], err = NewTempPutStream(
			name + fmt.Sprintf(".%d", i + 1), // i = 0 lefted to represent invalid
			dataServers[i], 
			perShard,
		)
		if err != nil {
			return nil, err
		}
	}
	encoder, err := NewEncoder()
	if err != nil {
		return nil, err
	}
	encoder.writers = writers
	return &RsPutStream{encoder: encoder}, nil
}

func (rsp *RsPutStream) Commit(success bool) error {
	var err error
	for _, w := range rsp.writers {
		err = w.(*TempPutStream).Commit(success)
	}
	return err
}

func NewRsGetStream(hash string, size int64, dataServers map[string]int, comple []string) (*RsGetStream, error) {
	servers := make([]string, utils.ALL_SHARDS)
	for server, id := range dataServers {
		servers[id - 1] = server // shard id starts from 1
	}
	readers := make([]io.Reader, utils.ALL_SHARDS)
	writers := make([]io.Writer, utils.ALL_SHARDS)
	perShard := (size + utils.DATA_SHARDS - 1) / utils.DATA_SHARDS
	alivable := 0
	for i := 0; i < utils.ALL_SHARDS; i++ {
		if servers[i] != "" {
			hashId := hash + fmt.Sprintf(".%d", i + 1)
			getStream, err := NewGetStream(servers[i], hashId)
			if err == nil {
				readers[i] = getStream
				alivable += 1
			}
		} else {
			servers[i] = comple[0]
			comple = comple[1:]
		}
	}

	if alivable < utils.DATA_SHARDS {
		return nil, fmt.Errorf("not have enough shards for reconstruct data. have %d need %d", alivable, utils.DATA_SHARDS)
	}

	for i := range readers {
		if readers[i] == nil {
			hashId := hash + fmt.Sprintf(".%d", i + 1)
			writer, err := NewTempPutStream(
				hashId,
				servers[i],
				perShard,
			)
			if err == nil {
				writers[i] = writer
			}
		}
	}
	dec, err := NewDecoder(readers, writers, size)
	if err != nil {
		return nil, err
	}
	return &RsGetStream{dec}, nil
}

func (rsg *RsGetStream) Seek(offset int64, whence int) (int64, error) {
	if whence != io.SeekCurrent {
		return -1, fmt.Errorf("invalid whence: %d, only support io.SeekCurrent", whence)
	}
	
	if offset < 0 {
		return -1, fmt.Errorf("only support forward seek")
	}

	for offset > 0 {
		length := int64(utils.BLOCK_SIZE)
		if length > offset {
			length = offset
		}
		buf := make([]byte, length)
		io.ReadFull(rsg, buf)
		offset -= length
	}
	return offset, nil
}

func (rsg *RsGetStream) Close() {
	for i := range rsg.writers {
		if rsg.writers[i] != nil {
			rsg.writers[i].(*TempPutStream).Commit(!rsg.notCommit[i])
		}
	}
}