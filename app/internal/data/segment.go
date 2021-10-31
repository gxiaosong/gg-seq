package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gouez/gg-seq/app/internal/biz"
)

type Segment struct {
	Id uint64 `gorm:"primaryKey;autoIncrement;autoIncrementIncrement=10000000"`
}

func (Segment) TableName() string {
	return "segment"
}

type seqRepo struct {
	data *Data
	log  *log.Helper
}

func NewSeqRepo(data *Data, logger log.Logger) biz.SeqRepo {
	return &seqRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
