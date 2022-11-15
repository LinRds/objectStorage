package locate

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/linrds/objectStorage/pkg/rabbitmq"
	"github.com/linrds/objectStorage/pkg/types"
)

var (
	objects = make(map[string]int)
	mu sync.Mutex
	rmu sync.RWMutex
)

func StartLocate() {
	rb := rabbitmq.NewRabbitmq(os.Getenv("RABBITMQ_SERVER"))
	defer rb.Close()

	rb.Bind("dataServers")
	ch := rb.Consume()

	for msg := range ch {
		hash, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}
		if id, ok := Exist(hash); ok {
			if id > 0 {
				rb.SingleSend(
					msg.ReplyTo, 
					types.LocateMessage{Id: id - 1, Addr: os.Getenv("LISTEN_ADDRESS")},
				)
			}
			
		}
	}
}

func Exist(hash string) (id int, ok bool) {
	rmu.RLock()
	id, ok = objects[hash]
	rmu.RUnlock()
	return
}

func CollectObjects() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		file := strings.Split(filepath.Base(files[i]), ".")
		if len(file) != 3 {
			panic(files[i])
		}
		hash := file[0]
		id, err := strconv.Atoi(file[1])
		if err != nil {
			panic(err)
		}
		objects[hash] = id
	}
}

func Add(hash string, id int) {
	mu.Lock()
	objects[hash] = id
	mu.Unlock()
}

func Del(hash string) {
	mu.Lock()
	delete(objects, hash)
	mu.Unlock()
}