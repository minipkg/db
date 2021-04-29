package mock

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	gorm_logger "gorm.io/gorm/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/log"
)

// New creates a new DB connection
func New(logger log.ILogger, conf minipkg_gorm.Config) (*minipkg_gorm.DB, *sqlmock.Sqlmock, error) {
	var db *gorm.DB
	var err error
	var mock sqlmock.Sqlmock
	var dbm *sql.DB

	newLogger := gorm_logger.New(logger, gorm_logger.Config{
		SlowThreshold:             conf.Log.SlowThreshold,
		Colorful:                  conf.Log.Colorful,
		IgnoreRecordNotFoundError: conf.Log.IgnoreRecordNotFoundError,
		LogLevel:                  gorm_logger.LogLevel(conf.Log.LogLevel),
	})

	dbm, mock, err = sqlmock.New() // mock sql.DB
	if err != nil {
		return nil, nil, err
	}

	switch conf.Dialect {
	case "postgres":
		db, err = gorm.Open(postgres.New(postgres.Config{Conn: dbm}), &gorm.Config{
			Logger: newLogger,
		})
	case "mysql":
		db, err = gorm.Open(mysql.New(mysql.Config{Conn: dbm}), &gorm.Config{
			Logger: newLogger,
		})
	case "sqlite":
		db, err = gorm.Open(sqlite.Dialector{Conn: dbm}, &gorm.Config{
			Logger: newLogger,
		})
	}
	if err != nil {
		return nil, nil, err
	}
	// Enable auto preload embeded entities
	db = db.Set("gorm:auto_preload", true)

	dbobj := &minipkg_gorm.DB{D: db}

	return dbobj, &mock, nil
}
