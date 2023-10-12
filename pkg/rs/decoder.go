package rs

import (
	"io"
	"sync"

	"github.com/klauspost/reedsolomon"

	"github.com/linrds/objectStorage/pkg/utils"
)

type decoder struct {
	wg        sync.WaitGroup
	writers   []io.Writer
	readers   []io.Reader
	notCommit []bool
	enc       reedsolomon.Encoder
	size      int64 // size是文件实际大小，避免因为分片向上取整读出填充的内容
	cache     []byte
	cacheSize int
	total     int64
}

func NewDecoder(readers []io.Reader, writers []io.Writer, size int64) (*decoder, error) {
	enc, err := reedsolomon.New(utils.DATA_SHARDS, utils.PARITY_SHARDS)
	if err != nil {
		return nil, err
	}
	return &decoder{
		readers:   readers,
		writers:   writers,
		enc:       enc,
		size:      size,
		notCommit: make([]bool, utils.ALL_SHARDS),
	}, nil
}

func (d *decoder) Read(b []byte) (n int, err error) {
	if d.cacheSize == 0 {
		err = d.getData()
		if err != nil {
			return 0, err
		}
	}
	length := len(b)
	if length > d.cacheSize {
		length = d.cacheSize
	}
	copy(d.cache[:length], b)
	d.cache = d.cache[length:]
	d.cacheSize -= length
	return length, nil
}

func (d *decoder) getData() (err error) {
	if d.total == d.size {
		return io.EOF
	}
	repairIds := make([]int, 0)
	shards := make([][]byte, utils.ALL_SHARDS)
	for i := 0; i < utils.ALL_SHARDS; i++ {
		if d.readers[i] == nil {
			repairIds = append(repairIds, i)
		} else {
			shards[i] = make([]byte, utils.BLOCK_SIZE_PER)
			n, err := io.ReadFull(d.readers[i], shards[i])
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				shards[i] = nil // 这个分片损毁了
			}
			if n != utils.BLOCK_SIZE_PER {
				shards[i] = shards[i][:n]
			}
		}
	}

	err = d.enc.Reconstruct(shards)
	if err != nil {
		return
	}

	for _, id := range repairIds {
		if d.writers[id] != nil {
			d.wg.Add(1)
			go func(i int) {
				defer d.wg.Done()
				_, err := d.writers[i].Write(shards[i])
				if err != nil { // 允许写入被恢复的分片失败
					d.notCommit[i] = true
				}
			}(id)
		}
	}
	d.wg.Wait()

	for i := 0; i < utils.DATA_SHARDS; i++ {
		shardSize := int64(len(shards[i]))
		if d.total+shardSize > d.size {
			shardSize = d.size - d.total
		}
		d.cache = append(d.cache, shards[i][:shardSize]...)
		d.total += shardSize
		d.cacheSize += int(shardSize)
	}

	return nil
}
