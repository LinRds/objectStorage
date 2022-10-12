package rs

import (
	"fmt"
	"io"
	"net/http"
	"log"
)

type GetStream struct {
	reader io.Reader
}

func (gs *GetStream) Read(p []byte) (n int, err error) {
	return gs.reader.Read(p)
}

func newGetStream(url string) (*GetStream, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dataServer return http code %d", resp.StatusCode)
	}
	return &GetStream{resp.Body}, nil
}

func NewGetStream(server, name string) (*GetStream, error) {
	if server == "" || name == "" {
		return nil, fmt.Errorf("invalid server and name of %s and %s", server, name)
	}
	url := fmt.Sprintf("http://%s/objects/%s", server, name)
	return newGetStream(url)
}