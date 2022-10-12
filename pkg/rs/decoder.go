package rs

import (
	"io"

	"github.com/klauspost/reedsolomon"

	"github.com/linrds/objectStorage/pkg/utils"
)

type decoder struct {
	writers []io.Writer
	readers []io.Reader
	enc reedsolomon.Encoder
	size int64
	cache []byte
	cacheSize int
	total int64
}

func NewDecoder(readers []io.Reader, writers []io.Writer, size int64) (*decoder, error) {
	enc, err := reedsolomon.New(utils.DATA_SHARDS, utils.PARITY_SHARDS, nil)
	if err != nil {
		return nil, err
	}
	return &decoder{
		readers: readers,
		writers: writers,
		enc: enc,
		size: size,
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

	for i := range repairIds {
		_, err = d.writers[i].Write(shards[repairIds[i]])
		if err != nil {
			return
		}
	}

	for i := 0; i < utils.DATA_SHARDS; i++ {
		shardSize := int64(len(shards[i]))
		if d.total + shardSize > d.size {
			shardSize = d.size - d.total
		}
		d.cache = append(d.cache, shards[i][:shardSize]...)
		d.total += shardSize
		d.cacheSize += int(shardSize)
	}

	return nil
}

