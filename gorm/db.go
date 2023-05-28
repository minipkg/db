package gorm

import (
	"context"
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
	Model(value interface{}) (*DB, error)
	WithContext(ctx context.Context) *DB
	ModelWithContext(ctx context.Context, model interface{}) (*DB, error)
	SchemeInit(model interface{}) (*DB, error)
	SchemeInitWithContext(ctx context.Context, model interface{}) (*DB, error)
	GormTx(gormDB *gorm.DB) *gorm.DB
}

// DB is the struct for a DB connection
type DB struct {
	GormDB        *gorm.DB
	isAutoMigrate bool
}

func (db *DB) DB() *gorm.DB {
	return db.GormDB
}

func (db *DB) Model(value interface{}) (*DB, error) {
	gormDB := db.GormDB.Model(value)

	if err := db.statementParse(value); err != nil {
		return nil, err
	}
	return &DB{
		GormDB:        gormDB,
		isAutoMigrate: db.isAutoMigrate,
	}, nil
}

func (db *DB) SchemeInit(model interface{}) (*DB, error) {
	err := db.statementParse(model)
	db.GormDB.Statement.Table = ""
	return db, err
}

func (db *DB) WithContext(ctx context.Context) *DB {
	return &DB{
		GormDB:        db.GormDB.WithContext(ctx),
		isAutoMigrate: db.isAutoMigrate,
	}
}

func (db *DB) ModelWithContext(ctx context.Context, model interface{}) (*DB, error) {
	d, err := db.Model(model)
	if err != nil {
		return nil, err
	}
	return d.WithContext(ctx), nil
}

func (db *DB) SchemeInitWithContext(ctx context.Context, model interface{}) (*DB, error) {
	return db.WithContext(ctx).SchemeInit(model)
}

func (db *DB) Close() error {
	sqlDB, err := db.GormDB.DB()
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
func New(logger log.Logger, conf Config) (*DB, error) {
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

	dbobj := &DB{
		GormDB:        db,
		isAutoMigrate: conf.IsAutoMigrate,
	}

	return dbobj, nil
}

func (db *DB) statementParse(model interface{}) error {
	return db.GormDB.Statement.Parse(model)
}

func (db *DB) GormTx(gormDB *gorm.DB) *gorm.DB {
	gormDB.Statement.Schema = db.DB().Statement.Schema
	return gormDB
}
