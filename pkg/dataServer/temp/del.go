package temp

import (
	"net/http"
	"os"
	"strings"
)

func del(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	infoPath := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	datPath := infoPath + ".dat"
	os.Remove(infoPath)
	os.Remove(datPath)
}