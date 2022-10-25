package temp

import (
	"strings"
	"strconv"
	"os"
	"net/url"

	"github.com/linrds/objectStorage/pkg/dataServer/locate"
	"github.com/linrds/objectStorage/pkg/utils"
)

func (t *tempInfo) hash() string {
	s := strings.Split(t.Hash, ".")
	return s[0]
}

func (t *tempInfo) id() int {
	s := strings.Split(t.Hash, ".")
	id, _ := strconv.Atoi(s[1])
	return id
}

func commitTempObject(datFile string, tempinfo *tempInfo) {
	f, _ := os.Open(datFile)
	d := url.PathEscape(utils.CalculateHash(f))
	f.Close()
	os.Rename(datFile, os.Getenv("STORAGE_ROOT") + "/objects/" + tempinfo.Hash + "." + d)
	locate.Add(tempinfo.hash(), tempinfo.id())
}