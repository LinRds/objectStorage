package main

import (
	"log"
	"net/http"
	"os"

	"github.com/linrds/objectStorage/pkg/dataServer/heartbeat"
	"github.com/linrds/objectStorage/pkg/dataServer/locate"
	"github.com/linrds/objectStorage/pkg/dataServer/objects"
	"github.com/linrds/objectStorage/pkg/dataServer/temp"
)

func main() {
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}