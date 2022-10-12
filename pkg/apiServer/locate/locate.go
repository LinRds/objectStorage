package locate

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/linrds/objectStorage/pkg/rabbitmq"
	"github.com/linrds/objectStorage/pkg/utils"
	"github.com/linrds/objectStorage/pkg/types"
)



func Locate(hash string) map[string]int {
	rb := rabbitmq.NewRabbitmq(os.Getenv("RABBIT_SERVER"))
	rb.BroadCast("dataServer", hash)
	ch := rb.Consume()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	go func () {
		for range ctx.Done() {}
		rb.Close()
	}()
	var dataServers = make(map[string]int)
	for i := 0; i < utils.ALL_SHARDS; i++ {
		msg := <- ch
		if len(msg.Body) == 0 {
			return dataServers
		}
		var info types.LocateMessage
		json.Unmarshal(msg.Body, &info)
		dataServers[info.Addr] = info.Id
	}
	return dataServers
}

func Exist (hash string) bool {
	dataServers := Locate(hash)
	return len(dataServers) >= utils.DATA_SHARDS
}