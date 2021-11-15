package comm

import "sync/atomic"

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
	Step      uint64
	Delta     uint64
	CurrentId uint64
	Remainder uint64
	LodingId  uint64
}

func (s *Segment) NextId() Result {
	id := atomic.AddUint64(&s.CurrentId, s.Delta)
	if id > s.MaxId {
		return Result{Code: OVER, Id: id}
	}
	if id >= s.LodingId {
		return Result{Code: Loading, Id: id}
	}
	return Result{Code: NORMAL, Id: id}
}

func (s *Segment) Useful() bool {
	return atomic.LoadUint64(&s.CurrentId) <= s.MaxId
}
