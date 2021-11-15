package data

import (
	"context"
	"database/sql"

	"github.com/gouez/gg-seq/server/config"
)

const (
	DB1 = "db1"
)

type Data struct {
	DB map[string]*sql.DB
}

func NewData(database []config.Database) *Data {

	m := make(map[string]*sql.DB)

	for _, value := range database {
		m[value.Name] = NewDB(value)
	}
	return &Data{
		DB: m,
	}
}

func NewDB(database config.Database) *sql.DB {

	db, err := sql.Open("mysql", database.URL)

	if err != nil {
		panic(err)
	}

	return db
}

func (d *Data) Close() {
	for _, value := range d.DB {
		value.Close()
	}
}

func (d *Data) GetDB(ctx context.Context, defaultDBName string) {

}
