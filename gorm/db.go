package gorm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/minipkg/log"
)

// IDB is the interface for a DB connection
type IDB interface {
	DB() *gorm.DB
	IsAutoMigrate() bool
}

// DB is the struct for a DB connection
type DB struct {
	D             *gorm.DB
	isAutoMigrate bool
}

func (db *DB) DB() *gorm.DB {
	return db.D
}

func (db *DB) IsAutoMigrate() bool {
	return db.isAutoMigrate
}

var _ IDB = (*DB)(nil)

// Config for a DB connection
type Config struct {
	Dialect       string
	DSN           string
	IsLogMode     bool
	IsAutoMigrate bool
}

// New creates a new DB connection
func New(conf Config, logger log.ILogger) (*DB, error) {
	db, err := gorm.Open(conf.Dialect, conf.DSN)
	if err != nil {
		return nil, err
	}
	db.SetLogger(logger)
	// Enable Logger, show detailed log
	db.LogMode(conf.IsLogMode)
	// Enable auto preload embeded entities
	db = db.Set("gorm:auto_preload", true)

	dbobj := &DB{
		D:             db,
		isAutoMigrate: conf.IsAutoMigrate,
	}

	return dbobj, nil
}
