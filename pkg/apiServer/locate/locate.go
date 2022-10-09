package locate

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/linrds/objectStorage/pkg/rabbitmq"
	"github.com/linrds/objectStorage/pkg/utils"
)

func Locate(hash string) []string {
	rb := rabbitmq.NewRabbitmq(os.Getenv("RABBIT_SERVER"))
	rb.BroadCast("dataServer", hash)
	ch := rb.Consume()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	go func () {
		for range ctx.Done() {}
		rb.Close()
	}()
	var dataServers = make([]string, 0)
	for msg := range ch {
		ds, _ := strconv.Unquote(string(msg.Body))
		if ds != "" {
			dataServers = append(dataServers, ds)
		}
	}
	return dataServers
}

func Exist (hash string) bool {
	dataServers := Locate(hash)
	return len(dataServers) >= utils.DATA_SHARDS
}