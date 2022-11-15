package main

import (
	"log"
	"net/http"
	"os"
	
	"github.com/linrds/objectStorage/pkg/apiServer/heartbeat"
	"github.com/linrds/objectStorage/pkg/apiServer/locate"
	"github.com/linrds/objectStorage/pkg/apiServer/objects"
	"github.com/linrds/objectStorage/pkg/apiServer/temp"
	"github.com/linrds/objectStorage/pkg/apiServer/versions"
)

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}