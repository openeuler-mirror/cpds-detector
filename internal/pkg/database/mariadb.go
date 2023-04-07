package database

import (
	"cpds/cpds-detector/internal/core"

	"gorm.io/gorm"
)

type mariadb struct {
	db *gorm.DB
}

func New(db *gorm.DB) Database {
	return &mariadb{
		db: db,
	}
}

func (m *mariadb) Init() error {
	if err := m.db.AutoMigrate(&core.Rule{}); err != nil {
		return err
	}

	if err := m.db.AutoMigrate(&core.Analysis{}); err != nil {
		return err
	}

	return nil
}
