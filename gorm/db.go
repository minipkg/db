package gorm

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"

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
func New(logger log.ILogger, conf Config) (*DB, error) {
	newLogger := gorm_logger.New(logger, gorm_logger.Config{
		SlowThreshold:             0,
		Colorful:                  false,
		IgnoreRecordNotFoundError: false,
		LogLevel:                  0,
	})

	db, err := gorm.Open(postgres.Open(conf.DSN), &gorm.Config{
		Logger: logger,
	})
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
