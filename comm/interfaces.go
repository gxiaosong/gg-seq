package comm

type IdGenerator interface {
	GetId() uint64
	GetIds(size int) []uint64
}

type IdGeneratorFactory interface {
	GetIdGenerator(bizType string) IdGenerator
}

type SegmentService interface {
	GetNextSegment(bizType string) *Segment
}
