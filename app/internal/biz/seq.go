package biz

import (
	"sync/atomic"

	"github.com/go-kratos/kratos/v2/log"
)

type ResultCode int

const (
	Over ResultCode = iota
	Loading
	Normal
)

type Result struct {
	Code ResultCode
	Id   uint64
}

type Segment struct {
	MaxId     uint64
	LoadingId uint64
	CurrentId uint64
	Delta     uint64
	Remainder uint64
}

func (s *Segment) NextId() *Result {
	id := atomic.AddUint64(&s.CurrentId, s.Delta)
	if id > s.MaxId {
		return &Result{Code: Over, Id: id}
	}
	if id >= s.LoadingId {
		return &Result{Code: Loading, Id: id}
	}
	return &Result{Code: Normal, Id: id}
}

func (s *Segment) Useful() bool {
	return atomic.LoadUint64(&s.CurrentId) <= s.MaxId
}

type SeqRepo interface {
	GetNextSegment(bizType string) (*Segment, error)
}

type SegmentService interface {
	GetNextSegment(bizType string) (*Segment, error)
}

type SeqUsecase struct {
	repo SeqRepo
	log  *log.Helper
}

func NewSeqUsecase(repo SeqRepo, logger log.Logger) *SeqUsecase {
	return &SeqUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *SeqUsecase) GetNextSegment(bizType string) (*Segment, error) {
	var (
		s   *Segment
		err error
	)
	if s, err = uc.repo.GetNextSegment(bizType); err != nil {
		return nil, err
	}
	return s, nil
}
