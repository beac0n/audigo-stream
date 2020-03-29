package main

import (
	"audigo-stream/src/audioStreamer"
	"io"
	"log"
	"net/http"
	"strings"
)

type ReqHandler struct {
}

var audioRoute = "/audio/"
var commandRoute = "/command/"

func (reqHandler *ReqHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	url := request.URL.String()
	if strings.HasPrefix(url, audioRoute) {
		reqHandler.streamAudio(responseWriter, url)
	} else if strings.HasPrefix(url, commandRoute) {
		// TODO: use commands to control browser
	} else {
		responseWriter.WriteHeader(404)
	}

}

func (reqHandler *ReqHandler) streamAudio(responseWriter http.ResponseWriter, url string) {
	targetUrl := strings.Replace(url, audioRoute, "", 1)

	log.Println("streaming audio from", targetUrl)

	// TODO: put streamer in synchronized map to access it with commands
	streamer := &audioStreamer.AudioStreamer{}
	defer streamer.Cleanup()
	audioStreamReader, err := streamer.Create(targetUrl)
	if err != nil {
		responseWriter.WriteHeader(400)
		log.Println("ERROR - create audio stream: ", err)
		return
	}

	responseWriter.WriteHeader(200)
	responseWriter.Header().Set("Content-Type", "audio/mpeg")
	responseWriter.Header().Set("Transfer-Encoding", "chunked")
	responseWriter.Header().Set("Connection", "keep-alive")

	if _, err := io.Copy(responseWriter, audioStreamReader); err != nil {
		log.Println("ERROR - ServeHTTP io.Copy: ", err)
	}
}

func main() {
	reqHandler := &ReqHandler{}
	listenAddress := ":8910"
	log.Println("server is listening on", listenAddress)
	if err := http.ListenAndServe(listenAddress, reqHandler); err != nil {
		log.Fatal("ERROR on starting webserver: ", err)
	}
}
