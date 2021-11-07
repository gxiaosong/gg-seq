package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/gouez/gg-seq/app/api/seq/v1"
	"github.com/gouez/gg-seq/app/internal/biz"
)

type SeqService struct {
	v1.UnimplementedSeqServer

	uc        *biz.SeqUsecase
	idFactory *biz.IdGeneratorFactory
	log       *log.Helper
}

func NewSeqService(uc *biz.SeqUsecase, idFactory *biz.IdGeneratorFactory, logger log.Logger) *SeqService {
	return &SeqService{
		uc:        uc,
		idFactory: idFactory,
		log:       log.NewHelper(logger),
	}
}

func (seq *SeqService) GetId(ctx context.Context, req *v1.GetIdReq) (*v1.GetIdResp, error) {
	idGen := seq.idFactory.GetGenerator(req.BizType)
	return &v1.GetIdResp{
		Id: idGen.NextId(),
	}, nil
}
