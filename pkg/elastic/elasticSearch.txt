package elastic

import (
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type MetaData struct {
	Name string
	Version string
	Size string
	Hash string
}

var esCli *elasticsearch.Client

func InitEsClient() (err error) {
	if esCli != nil {
		return
	}
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}
	esCli, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Printf("Error creating the client: %s", err)
	} else {
		log.Println(esCli.Info())
	}
	return
}

func getMetaData(name, version string) (meta MetaData, err error) {
	if esCli == nil {
		err = fmt.Errorf("elasticSearch client not inited, get metadata failed")
		return
	}
	req := esapi.IndexRequest{
		Index: "metadata",
		Body: strings.NewReader(`{"name": "` + name + `"}`),
	}
	return
}

func PutMetaData(name, version, size, hash string) (err error) {
	if esCli == nil {
		err = fmt.Errorf("elasticSearch client not inited, put metadata failed")
		return
	}
	req := esapi.IndexRequest{
		Index: "metadata",
		Body: strings.NewReader(`{"name": "` + name + `"}`),
	}
}