package rs

import (
	"fmt"
	"io"

	"github.com/linrds/objectStorage/pkg/utils"
)

type ResumableRsGetStream struct {
	*decoder
}

func NewTempGetStream(server, uuid string) (*GetStream, error) {
	return newGetStream(fmt.Sprintf("http://%s/temp/%s", server, uuid)) // 这个函数本质就是读取server上名为name的文件，不要纠结传入的是hash还是uuid
}

func NewRsResumableGetStream(dataServers, uuids []string, size int64) (*ResumableRsGetStream, error) {
	readers := make([]io.Reader, utils.ALL_SHARDS)
	var err error
	for i := 0; i < utils.ALL_SHARDS; i++ {
		readers[i], err = NewTempGetStream(dataServers[i], uuids[i])
		if err != nil {
			return nil, err
		}
	}
	writers := make([]io.Writer, utils.ALL_SHARDS) // writers都为nil，读出时不会写入受损的分片
	decoder, err := NewDecoder(readers, writers, size)
	if err != nil {
		return nil, err
	}
	return &ResumableRsGetStream{decoder: decoder}, nil
}
