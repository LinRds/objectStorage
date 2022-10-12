package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	if len(digest) < 8 || digest[:8] != "SHA-256=" {
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

func CalculateHash(r io.Reader) string {
	h := sha256.New()
	io.Copy(h, r)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func GetNameFromUrl(url *url.URL) string {
	name := strings.Split(url.EscapedPath(), "/")[2]
	return name
}