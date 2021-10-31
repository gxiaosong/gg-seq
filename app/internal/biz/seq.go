package biz

import (
	"github.com/go-kratos/kratos/v2/log"
)

type Segment struct {
}

type SeqRepo interface {
}

type SeqUsecase struct {
	repo SeqRepo
	log  *log.Helper
}

func NewSeqUsecase(repo SeqRepo, logger log.Logger) *SeqUsecase {
	return &SeqUsecase{repo: repo, log: log.NewHelper(logger)}
}
