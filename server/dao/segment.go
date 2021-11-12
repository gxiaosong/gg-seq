package dao

import (
	"database/sql"
	"time"
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
	db *sql.DB
}

func NewSegmentDao(db *sql.DB) SegmentDao {

	return SegmentDao{
		db: db,
	}
}

func (dao SegmentDao) GetSegmentByBizType(bizType string) (*Segment, error) {
	var (
		s   Segment
		row *sql.Row
	)
	if row = dao.db.QueryRow("select * from segment_id where biz_type = ? ", bizType); row.Err() != nil {
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
	)
	sql := `
		update segment_id set max_id= ?,update_time=now(), version=version+1
		where id=? and max_id=? and version=? and biz_type=?
	`
	result, err = dao.db.Exec(sql, newMaxId, id, oldMaxId, version, bizType)
	if err != nil {
		return 0, err
	}
	rows, _ := result.RowsAffected()
	return rows, nil
}
