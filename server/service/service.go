package service

type SegmentService interface {
	GetNextSegment(bizType string) *Segment
}
