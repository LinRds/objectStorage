package temp

import (
	"net/http"
	"os"

	"github.com/linrds/objectStorage/pkg/utils"
)

func del(w http.ResponseWriter, r *http.Request) {
	uuid := utils.GetNameFromUrl(r.URL)
	infoPath := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	datPath := infoPath + ".dat"
	os.Remove(infoPath)
	os.Remove(datPath)
}