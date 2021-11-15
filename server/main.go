package main

import (
	"github.com/gouez/gg-seq/server/config"
	"github.com/gouez/gg-seq/server/controller"
	"github.com/gouez/gg-seq/server/server"
)

func main() {
	conf := config.NewConfigFromFile("config.json")
	server.RunHttpServer(conf.Server, controller.Handlers)
}
