package entity

import "gorm.io/gorm"

type Entity interface {
	Register(conn *gorm.DB) error

	Put(conn *gorm.DB) (uint, error)
	Get(conn *gorm.DB) error
}
