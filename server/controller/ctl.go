package controller

import (
	"net/http"
	"strconv"

	"github.com/gouez/gg-seq/comm"
)

func GetHandlers(idGeneratorFactory comm.IdGeneratorFactory) map[string]http.HandlerFunc {
	handlers := make(map[string]http.HandlerFunc)

	handlers["/get"] = func(rw http.ResponseWriter, r *http.Request) {
		bizType := r.URL.Query().Get("bizType")
		rw.Write([]byte(strconv.FormatUint(idGeneratorFactory.GetIdGenerator(bizType).GetId(), 10)))
	}
	return handlers
}
