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

func (dao SegmentDao) GetNextSegment(bizType string, retryCount int) *Segment {
	var (
		db     = dao.data.DB[data.DB1]
		row    *sql.Row
		s      Segment
		err    error
		result sql.Result
		tx0    *sql.Tx
	)

	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	for i := 0; i < retryCount; i++ {
		tx0, err = db.BeginTx(context.Background(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
		})
		if row = tx0.QueryRow("select * from segment_id where biz_type = ? ", bizType); row.Err() != nil {
			err = row.Err()
			tx0.Rollback()
			return nil
		}
		if err = row.Scan(&s.Id, &s.BizType, &s.BeginId, &s.MaxId, &s.Step, &s.Delta, &s.Remainder, &s.CreateTime, &s.UpdateTime, &s.Version); err != nil {
			tx0.Rollback()
			return nil
		}
		sql := `
			update segment_id set max_id= ?,update_time=now(), version=version+1
			where id=? and max_id=? and version=? and biz_type=?
		`
		newMaxId := s.MaxId + s.Step
		result, err = db.Exec(sql, newMaxId, s.Id, s.MaxId, s.Version, bizType)
		if err != nil {
			tx0.Rollback()
			return nil
		}
		tx0.Commit()
		rows, _ := result.RowsAffected()
		if rows > 0 {
			s.MaxId = newMaxId
			return &s
		}
	}
	return nil
}

func (dao SegmentDao) GetSegmentByBizType(bizType string) (*Segment, error) {
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

func (dao SegmentDao) UpdateMaxId(id, uint64, newMaxId uint64, oldMaxId uint64, version uint64, bizType string) (int64, error) {
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
