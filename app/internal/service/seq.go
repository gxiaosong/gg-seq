package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/gouez/gg-seq/app/api/seq/v1"
	"github.com/gouez/gg-seq/app/internal/biz"
)

type SeqService struct {
	v1.UnimplementedSeqServer

	uc  *biz.SeqUsecase
	log *log.Helper
}

func NewSeqService(uc *biz.SeqUsecase, logger log.Logger) *SeqService {
	return &SeqService{uc: uc, log: log.NewHelper(logger)}
}

func (seq *SeqService) GetId(context.Context, *v1.GetIdReq) (*v1.GetIdResp, error) {
	return nil, nil
}
