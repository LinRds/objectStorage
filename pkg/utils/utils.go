package utils

import (
	"net/http"
	"strconv"
)

const (
	DATA_SHARDS = 4
	PARITY_SHARDS = 2
	ALL_SHARDS = DATA_SHARDS + PARITY_SHARDS
	BLOCK_SIZE_PER = 8000
	BLOCK_SIZE = BLOCK_SIZE_PER * DATA_SHARDS
)
func GetHashFromHeader(h http.Header) string {
	digest := h.Get("digest")
	if len(digest) < 8 || digest[:7] != "SHA-256" {
		return ""
	}
	return digest[8:]
}

func GetSizeFromHeader(h http.Header) int64 {
	s := h.Get("size")
	size, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return -1
	}
	return size
}