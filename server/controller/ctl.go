package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gouez/gg-seq/comm"
	"github.com/gouez/gg-seq/server/data"
	"github.com/gouez/gg-seq/server/service"
)

func GetHandlers(data *data.Data, idGeneratorFactory comm.IdGeneratorFactory) map[string]http.HandlerFunc {
	handlers := make(map[string]http.HandlerFunc)

	handlers["/get"] = func(rw http.ResponseWriter, r *http.Request) {
		bizType := r.URL.Query().Get("bizType")
		rw.Write([]byte(strconv.FormatUint(idGeneratorFactory.GetIdGenerator(bizType).GetId(), 10)))
	}

	handlers["/get/batch"] = func(rw http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		bizType := values.Get("bizType")
		size := values.Get("size")
		isize, _ := strconv.Atoi(size)
		ids := idGeneratorFactory.GetIdGenerator(bizType).GetIds(isize)
		v, err := json.Marshal(ids)
		if err != nil {
			log.Fatalln(err)
		}
		rw.Write([]byte(v))
	}

	handlers["/getSegment"] = func(rw http.ResponseWriter, r *http.Request) {
		bizType := r.URL.Query().Get("bizType")
		service := service.NewDBSegmentService(data)
		s := service.GetNextSegment(bizType)
		v, err := json.Marshal(s)
		if err != nil {
			log.Fatalln(err)
		}
		rw.Write([]byte(v))
	}

	return handlers
}
