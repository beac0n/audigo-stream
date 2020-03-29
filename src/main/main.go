package main

import (
	"log"
	"net/http"
	"sync"
)

func main() {
	reqHandler := &ReqHandler{audioStreamerMap: &sync.Map{}}
	listenAddress := ":8910"
	log.Println("server is listening on", listenAddress)
	if err := http.ListenAndServe(listenAddress, reqHandler); err != nil {
		log.Fatal("ERROR on starting webserver: ", err)
	}
}
