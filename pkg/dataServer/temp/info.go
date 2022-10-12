package temp

import (
	"strconv"
	"strings"
)

type tempInfo struct {
	Uuid string
	Hash string
	Size int64
}

func (t *tempInfo) GetHash() string {
	return strings.Split(t.Hash, ".")[0]
}

func (t *tempInfo) GetId() int {
	ids := strings.Split(t.Hash, ".")[1]
	id, _ := strconv.Atoi(ids)
	return id
}