package service

import (
	"database/sql"
)

type SegmentService interface {
	GetNextSegment(bizType string) *Segment
}

type DBSegmentService struct {
}

func NewDBSegmentService(db *sql.DB) DBSegmentService {
	return DBSegmentService{}
}

func (s DBSegmentService) GetNextSegment(bizType string) *Segment {

	return nil
}
