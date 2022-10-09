package heartbeat

import (
	"os"
	"time"

	"github.com/linrds/objectStorage/pkg/rabbitmq"
)
func StartHeartbeat() {
	rb := rabbitmq.NewRabbitmq(os.Getenv("RABBIT_SERVER"))
	defer rb.Close()
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for range ticker.C {
		rb.BroadCast("apiServer", os.Getenv("LISTEN_ADDRESS"))
	}
}