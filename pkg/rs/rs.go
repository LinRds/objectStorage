package rs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/linrds/objectStorage/pkg/utils"
)

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
		return nil, err
	}
	client := http.Client{}
	request.Header.Set("size", fmt.Sprintf("%d", size))
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
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
	client := http.Client{}
	resp, err := client.Do(request)
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
	var writers = make([]io.Writer, 0, len(dataServers))
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

func (rs *RsPutStream) Commit(success bool) error {
	var err error
	for _, w := range rs.writers {
		err = w.(*TempPutStream).Commit(success)
	}
	return err
}

func NewRsGetStream(hash string, size int64, dataServers map[string]int, comple []string) (*RsGetStream, error) {
	servers := make([]string, utils.ALL_SHARDS)
	for server, id := range dataServers {
		servers[id] = server
	}
	readers := make([]io.Reader, utils.ALL_SHARDS)
	writers := make([]io.Writer, utils.ALL_SHARDS)
	for i := 0; i < utils.ALL_SHARDS; i++ {
		hashId := hash + fmt.Sprintf(".%d", i)
		if len(servers[i]) == 0 {
			writer, err := NewTempPutStream(
				hashId,
				comple[0],
				size,
			)
			if err != nil {
				return nil, err
			}
			writers[i] = writer
			comple = comple[1:]
		} else {
			getStream, err := NewGetStream(servers[i], hashId)
			if err != nil {
				return nil, err
			}
			readers[i] = getStream
		}
	}
	dec, err := NewDecoder(readers, writers, size)
	if err != nil {
		return nil, err
	}
	return &RsGetStream{dec}, nil
}