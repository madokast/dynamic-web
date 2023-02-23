package main

import (
	"errors"
	"net/http"

	"github.com/madokast/dynamic-web/web"
)

type HelloRequest struct {
	Name string `json:"name"`
}

type HelloResponse struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

func Hello(request *HelloRequest) (*HelloResponse, error) {
	if len(request.Name) == 0 {
		return nil, errors.New("name is empty")
	} else {
		return &HelloResponse{
			Name: request.Name,
			Msg:  "hello",
		}, nil
	}
}

func main() {
	http.HandleFunc("/hello", web.Do(Hello))
	http.ListenAndServe(":8080", nil)
}
