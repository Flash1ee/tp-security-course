package main

import (
	"log"
	"net/http"

	"http-proxy/pkg/proxy"
	"http-proxy/pkg/utils"
)

func main() {
	srv := proxy.New(utils.GetLogger())
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
