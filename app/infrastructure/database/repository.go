package database

import (
	"context"
	"fmt"
	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type dbRepository struct {
	db *gorm.DB
}

type dbTransaction struct {
	db *gorm.DB
}

type Config struct {
	User     string
	Password string
	Host     string
	DBName   string
	Debug    bool
}

func NewDBRepository(ctx context.Context, config Config) repository.Repository {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err.Error())
	}

	db = db.Session(&gorm.Session{Context: ctx})

	if config.Debug {
		db = db.Debug()
	}

	return &dbRepository{db: db}
}

func (r *dbRepository) RunTransaction(f func(repository.Transaction) error) error {
	tx := r.db.Begin()

	if err := f(&dbTransaction{db: tx}); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *dbRepository) User() repository.UserRepositoryQuery {
	return &dbUserRepository{}
}
func (tx *dbTransaction) User() repository.UserRepositoryCommand {
	return &dbUserRepository{}
}
func (r *dbRepository) Group() repository.GroupRepositoryQuery {
	return &dbGroupRepository{}
}
func (tx *dbTransaction) Group() repository.GroupRepositoryCommand {
	return &dbGroupRepository{}
}
