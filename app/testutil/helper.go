package testutil

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DBMock(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
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
			}},
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
