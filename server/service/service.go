package service

import (
	"github.com/gouez/gg-seq/comm"
	"github.com/gouez/gg-seq/server/dao"
	"github.com/gouez/gg-seq/server/data"
)

type DBSegmentService struct {
	dao dao.SegmentDao
}

func NewDBSegmentService(data *data.Data) DBSegmentService {
	return DBSegmentService{
		dao: dao.NewSegmentDao(data),
	}
}

func (service DBSegmentService) GetNextSegment(bizType string) *comm.Segment {
	s := service.dao.GetNextSegment(bizType, 3)
	if s == nil {
		return nil
	}
	return &comm.Segment{
		MaxId:     s.MaxId,
		Step:      s.Step,
		Remainder: s.Remainder,
		LodingId:  (s.MaxId - s.Step) * 20 / 100,
		Delta:     s.Delta,
		CurrentId: (s.MaxId - s.Step),
	}
}
