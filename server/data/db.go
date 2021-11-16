package data

import (
	"database/sql"

	"github.com/gouez/gg-seq/server/config"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DB1 = "db1"
)

func NewDB(database config.Database) *sql.DB {

	db, err := sql.Open("mysql", database.URL)

	if err != nil {
		panic(err)
	}

	return db
}
