package database_test

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func dbMock(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New()=_, _, %v; want _, _, nil", err)
	}

	gormDB, err := gorm.Open(
		mysql.Dialector{
			Config: &mysql.Config{
				DriverName:                "mysql",
				Conn:                      db,
				SkipInitializeWithVersion: true,
			},
		},
		&gorm.Config{
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		t.Fatalf("gorm.Open(mysql, %v)=_, %v; want _, nil", db, err)
	}

	gormDB.Debug()

	return mock, gormDB
}

func toDriverValues[T any](t *testing.T, v ...T) []driver.Value {
	t.Helper()
	result := make([]driver.Value, len(v))
	for i, e := range v {
		result[i] = driver.Value(e)
	}
	return result
}
