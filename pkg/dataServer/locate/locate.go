package locate

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/linrds/objectStorage/pkg/rabbitmq"
)

var (
	objects = make(map[string]struct{})
	mu sync.Mutex
	rmu sync.RWMutex
)

func StartLocate() {
	rb := rabbitmq.NewRabbitmq(os.Getenv("RABBIT_SERVER"))
	defer rb.Close()

	rb.Bind("dataServer")
	ch := rb.Consume()

	for msg := range ch {
		hash, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}
		if Exist(hash) {
			rb.SingleSend(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}

func Exist(hash string) bool {
	rmu.RLock()
	_, ok := objects[hash]
	rmu.RUnlock()
	return ok
}

func CollectObjects() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		hash := filepath.Base(files[i])
		objects[hash] = struct{}{}
	}
}

func Add(name string) {
	
}