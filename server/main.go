package main

import (
	"github.com/gouez/gg-seq/server/config"
	"github.com/gouez/gg-seq/server/controller"
	"github.com/gouez/gg-seq/server/data"
	"github.com/gouez/gg-seq/server/server"
	"github.com/gouez/gg-seq/server/service"
)

func main() {
	conf := config.NewConfigFromFile("config.json")
	data := data.NewData(conf)
	idgen := service.NewIdGeneratorFactory(data)
	server.RunHttpServer(conf.Server, controller.GetHandlers(idgen))
}
