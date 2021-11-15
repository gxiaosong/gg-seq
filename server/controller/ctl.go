package controller

import (
	"net/http"

	"github.com/gouez/gg-seq/comm"
	"github.com/gouez/gg-seq/server/service"
)

var (
	Handlers           = make(map[string]http.HandlerFunc)
	idGeneratorFactory = service.NewIdGeneratorFactory()
)

func init() {
	Handlers["/get"] = func(rw http.ResponseWriter, r *http.Request) {
		bizType := r.URL.Query().Get("bizType")
		rw.Write(comm.I64tob(idGeneratorFactory.GetIdGenerator(bizType).GetId()))
	}
}
