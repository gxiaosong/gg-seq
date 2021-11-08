package data

import (
	"database/sql"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/gouez/gg-seq/app/internal/biz"
	"gorm.io/gorm"
)

type Segment struct {
	Id         uint64 `gorm:"primaryKey;autoIncrement;autoIncrementIncrement=10000000"`
	BizType    string `gorm:"not null;size:60;index"`
	BeginId    uint64 `gorm:"default:0"`
	MaxId      uint64 `gorm:"default:0"`
	Step       uint64 `gorm:"default:0"`
	Delta      uint64 `gorm:"default:1"`
	Remainder  uint64 `gorm:"default:0"`
	CreateTime time.Time
	UpdateTime time.Time
	Version    uint64 `gorm:"default:0"`
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
		log:  log.NewHelper(log.With(logger, "module", "data/repoRepo")),
	}
}

func (repo *seqRepo) GetNextSegment(bizType string) (*biz.Segment, error) {
	var (
		seg Segment
		err error
	)
	defer func() {
		if err != nil {
			repo.log.Errorw(err)
		}
	}()

	err = repo.data.db.Transaction(func(tx *gorm.DB) error {
		var result *gorm.DB
		if result = tx.Where(&Segment{BizType: bizType}).First(&seg); result.Error != nil {
			return result.Error
		}
		newMaxId := seg.MaxId + seg.Step
		oldMaxId := seg.MaxId
		if result = tx.Model(&seg).Where("max_id=? and version=? and biz_type=?", oldMaxId, seg.Version, bizType).Updates(&Segment{
			MaxId:      newMaxId,
			Version:    seg.Version + 1,
			UpdateTime: time.Now(),
		}); result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 1 {
			seg.MaxId = newMaxId
			return nil
		}
		return biz.ErrNotUpdated
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	s := &biz.Segment{
		CurrentId: seg.MaxId - seg.Step,
		MaxId:     seg.MaxId,
		Delta:     seg.Delta,
		Remainder: seg.Remainder,
	}
	s.LoadingId = (atomic.LoadUint64(&s.CurrentId) + seg.Step) * 20 / 100

	return s, nil
}
