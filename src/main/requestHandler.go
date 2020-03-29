package main

import (
	"audigo-stream/src/audioStreamer"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

type ReqHandler struct {
	audioStreamerMap *sync.Map
}

var audioRoute = "/audio/"
var commandRoute = "/command/"

func (reqHandler *ReqHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	url := request.URL.String()
	if strings.HasPrefix(url, audioRoute) {
		targetUrl, id, err := reqHandler.getTargetUrlAndId(url)
		if err != nil {
			responseWriter.WriteHeader(404)
			return
		}

		reqHandler.streamAudio(responseWriter, targetUrl, id)
	} else if strings.HasPrefix(url, commandRoute) {
		id, command := reqHandler.getIdAndCommand(url)
		streamer, ok := reqHandler.audioStreamerMap.Load(id)
		if !ok {
			responseWriter.WriteHeader(404)
			return
		}

		log.Println(streamer, command)

		// TODO: use commands to control browser
	} else {
		responseWriter.WriteHeader(404)
	}

}

func (reqHandler *ReqHandler) streamAudio(responseWriter http.ResponseWriter, targetUrl string, id string) {

	log.Println("streaming audio from", targetUrl)

	streamer := &audioStreamer.AudioStreamer{}

	reqHandler.audioStreamerMap.Store(id, streamer)

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

func (reqHandler *ReqHandler) getIdAndCommand(url string) (string, string) {
	split := strings.Split(url, "/")
	if len(split) <= 3 {
		return split[2], ""
	}
	if len(split) <= 2 {
		return "", ""
	}

	return split[2], split[3]
}

func (reqHandler *ReqHandler) getTargetUrlAndId(url string) (string, string, error) {
	split := strings.Split(url, "/")
	targetUrl := strings.Join(split[3:], "/")
	id := split[2]

	if strings.HasPrefix(id, "http") {
		return "", "", errors.New("invalid id")
	}
	return targetUrl, id, nil
}
