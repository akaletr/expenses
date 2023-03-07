package storage

import (
	"cmd/main/main.go/internal/config"

	"gorm.io/gorm"
)

type Storage interface {
	Start(cfg config.Database) error
	Provide(entities ...Entity) error

	Put(entity Entity) (uint, error)
	Get(entity Entity, id uint) (Entity, error)
}

type Entity interface {
	Register(conn *gorm.DB) error
}
