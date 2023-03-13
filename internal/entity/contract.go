package entity

import (
	"gorm.io/gorm"
)

type Entity interface {
	Register(conn *gorm.DB) error
}
