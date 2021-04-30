package gorm

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"

	"github.com/minipkg/log"
)

// IDB is the interface for a DB connection
type IDB interface {
	DB() *gorm.DB
	Close() error
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

func (db *DB) Close() error {
	sqlDB, err := db.D.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (db *DB) IsAutoMigrate() bool {
	return db.isAutoMigrate
}

var _ IDB = (*DB)(nil)

// Config for a DB connection
type Config struct {
	Dialect       string
	DSN           string
	IsAutoMigrate bool
	Log           LogConfig
}

type LogConfig struct {
	SlowThreshold             time.Duration
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	LogLevel                  int
}

// New creates a new DB connection
func New(logger log.ILogger, conf Config) (*DB, error) {
	var db *gorm.DB
	var err error

	newLogger := gorm_logger.New(logger, gorm_logger.Config{
		SlowThreshold:             conf.Log.SlowThreshold,
		Colorful:                  conf.Log.Colorful,
		IgnoreRecordNotFoundError: conf.Log.IgnoreRecordNotFoundError,
		LogLevel:                  gorm_logger.LogLevel(conf.Log.LogLevel),
	})

	switch conf.Dialect {
	case "postgres":
		db, err = gorm.Open(postgres.Open(conf.DSN), &gorm.Config{
			Logger: newLogger,
		})
	case "mysql":
		db, err = gorm.Open(mysql.Open(conf.DSN), &gorm.Config{
			Logger: newLogger,
		})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(conf.DSN), &gorm.Config{
			Logger: newLogger,
		})
	}
	if err != nil {
		return nil, err
	}

	db = db.Set("gorm:auto_preload", true)

	dbobj := &DB{
		D:             db,
		isAutoMigrate: conf.IsAutoMigrate,
	}

	return dbobj, nil
}
