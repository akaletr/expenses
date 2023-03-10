package entity

import (
	"gorm.io/gorm"
)

type Entity interface {
	Register(conn *gorm.DB) error

	Put(conn *gorm.DB) error
	Get(conn *gorm.DB, id uint) error
}
