package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gouez/gg-seq/comm"
)

const (
	urlPatten = "%s/getSegment?bizType=%s"
)

type httpSegmentService struct {
	url string
}

func NewHttpSegmentService(url string) comm.SegmentService {
	return httpSegmentService{
		url: url,
	}
}

func (service httpSegmentService) GetNextSegment(bizType string) *comm.Segment {

	var (
		resp *http.Response
		err  error
	)

	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	resp, err = http.Get(fmt.Sprintf(urlPatten, service.url, bizType))
	if err == nil {
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		return nil
	}
	var s comm.Segment
	json.Unmarshal(body, &s)
	return &comm.Segment{
		MaxId:     s.MaxId,
		Step:      s.Step,
		Remainder: s.Remainder,
		LodingId:  (s.MaxId - s.Step) * 20 / 100,
		Delta:     s.Delta,
		CurrentId: (s.MaxId - s.Step),
	}
}
