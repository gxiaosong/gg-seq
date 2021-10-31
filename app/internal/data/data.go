package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/gouez/gg-seq/app/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewDB,
	NewData,
	NewGreeterRepo,
)

// Data .
type Data struct {
	db  *gorm.DB
	log *log.Helper
}

// NewData .
func NewData(db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(log.With(logger, "module", "data"))

	d := &Data{
		db:  db,
		log: log,
	}

	return d, func() {

	}, nil
}

//NewDB
func NewDB(conf *conf.Data, logger log.Logger) *gorm.DB {
	log := log.NewHelper(log.With(logger, "module", "data/gorm"))

	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}

	if err := db.AutoMigrate(); err != nil {
		log.Fatal(err)
	}
	return db
}
