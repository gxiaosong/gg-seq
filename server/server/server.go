package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gouez/gg-seq/server/config"
)

func RunHttpServer(server config.Server, handlers map[string]http.HandlerFunc) {

	for key, element := range handlers {
		http.HandleFunc(key, element)
	}
	log.Printf("http server listens on [:%d] \n", server.Port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", server.Port), nil))
}
