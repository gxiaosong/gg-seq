package dao

import (
	"context"
	"database/sql"
	"time"

	"github.com/gouez/gg-seq/server/data"
)

type Segment struct {
	Id         uint64
	BizType    string
	BeginId    uint64
	MaxId      uint64
	Step       uint64
	Delta      uint64
	Remainder  uint64
	CreateTime time.Time
	UpdateTime time.Time
	Version    uint64
}

type SegmentDao struct {
	data *data.Data
}

func NewSegmentDao(data *data.Data) SegmentDao {

	return SegmentDao{
		data: data,
	}
}

func (dao SegmentDao) GetNextSegment(bizType string) *Segment {
	var db = dao.data.DB[data.DB1]
	ctx := context.Background()
	tx0, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
	})
	if err != nil {
		return nil
	}
	ctx = context.WithValue(ctx, data.TX, tx0)
	return nil
}

func (dao SegmentDao) GetSegmentByBizType(ctx context.Context, bizType string) (*Segment, error) {
	var (
		s   Segment
		row *sql.Row
		db  = dao.data.DB[data.DB1]
	)
	if row = db.QueryRow("select * from segment_id where biz_type = ? ", bizType); row.Err() != nil {
		return nil, row.Err()
	}
	if err := row.Scan(&s.Id, &s.BizType, &s.BeginId, &s.MaxId, &s.Step, &s.Delta, &s.Remainder, &s.CreateTime, &s.UpdateTime, &s.Version); err != nil {
		return nil, err
	}
	return &s, nil
}

func (dao SegmentDao) UpdateMaxId(ctx context.Context, id, uint64, newMaxId uint64, oldMaxId uint64, version uint64, bizType string) (int64, error) {
	var (
		result sql.Result
		err    error
		db     = dao.data.DB[data.DB1]
	)
	sql := `
		update segment_id set max_id= ?,update_time=now(), version=version+1
		where id=? and max_id=? and version=? and biz_type=?
	`
	result, err = db.Exec(sql, newMaxId, id, oldMaxId, version, bizType)
	if err != nil {
		return 0, err
	}
	rows, _ := result.RowsAffected()
	return rows, nil
}
