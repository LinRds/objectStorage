package rs

import (
	"io"
	"log"
	"sync"

	"github.com/klauspost/reedsolomon"
	"github.com/linrds/objectStorage/pkg/utils"
)

type encoder struct {
	wg sync.WaitGroup
	writers []io.Writer
	enc reedsolomon.Encoder
	cache []byte
}

func NewEncoder(writers []io.Writer) (*encoder, error){
	enc, err := reedsolomon.New(utils.DATA_SHARDS, utils.PARITY_SHARDS)
	if err != nil {
		return nil, err
	}
	return &encoder{enc: enc, writers: writers}, err
}

//TODO: error 处理
func (e *encoder) Write(b []byte) (n int, err error) {
	length := len(b)
	current := 0
	for length != 0 {
		next := utils.BLOCK_SIZE - len(e.cache)
		if next > length {
			next = length
		}
		e.cache = append(e.cache, b[current:current+next]...)
		if len(e.cache) == utils.BLOCK_SIZE {
			e.Flush()
		}
		current += next
		length -= next
	}
	return len(b), nil
}

func (e *encoder) Flush() {
	if len(e.cache) == 0 {
		return
	}
	shards, err := e.enc.Split(e.cache)
	if err != nil {
		log.Printf("Flush failed when Split shard: %v", err)
		return
	}
	err = e.enc.Encode(shards)
	if err != nil {
		log.Printf("Flush failed when encode shard: %v", err)
		return
	}
	for i := range e.writers {
		e.wg.Add(1)
		go func(i int) {
			defer e.wg.Done()
			e.writers[i].Write(shards[i])
		}(i)
	}
	e.wg.Wait()
	// remember to clear cache
	e.cache = []byte{}
}