package database

import "gorm.io/gorm"

type DBUserRepository = dbUserRepository

func (r *DBUserRepository) SetDB(db *gorm.DB) {
	r.db = db
}

type DBGroupRepository = dbGroupRepository

func (r *DBGroupRepository) SetDB(db *gorm.DB) {
	r.db = db
}
