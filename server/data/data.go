package data

import (
	"database/sql"

	"github.com/gouez/gg-seq/server/config"
)

type Data struct {
	DB     map[string]*sql.DB
	Config *config.Config
}

func NewData(config *config.Config) *Data {

	m := make(map[string]*sql.DB)

	for _, value := range config.Database {
		m[value.Name] = NewDB(value)
	}

	data := &Data{
		DB: m,
	}
	return data
}

func (d *Data) Close() {
	for _, value := range d.DB {
		value.Close()
	}
}
