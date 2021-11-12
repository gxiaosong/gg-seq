package service

type ResultCode int

const (
	OVER ResultCode = iota
	NORMAL
	Loading
)

type Result struct {
	Code ResultCode
	Id   uint64
}

type Segment struct {
	MaxId     uint64
	Step      int
	Delta     int
	Remainder int
}

func NewSegment() *Segment {
	return nil
}

func (s *Segment) NextId() Result {

	return Result{}
}
