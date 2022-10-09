package heartbeat

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/linrds/objectStorage/pkg/rabbitmq"
)

var (
	dataServers = make(map[string]time.Time)
	mu sync.Mutex
	rmu sync.RWMutex
)

func ListenHeartbeat() {
	rb := rabbitmq.NewRabbitmq(os.Getenv("RABBIT_SERVER"))
	defer rb.Close()
	go removeExpiredDataServer()
	rb.Bind("apiServer")
	ch := rb.Consume()
	for msg := range ch {
		dataServer, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			log.Println(err)
		}
		mu.Lock()
		dataServers[dataServer] = time.Now()
		mu.Unlock()
	}
}

func removeExpiredDataServer() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for range ticker.C {
		mu.Lock()
		for dataServer, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, dataServer)
			}
		}
		mu.Unlock()
	}
}