package storage

import (
	"cmd/main/main.go/internal/config"
	"cmd/main/main.go/internal/entity"
)

type Storage interface {
	Start(cfg config.Database) error
	Provide(entities ...entity.Entity) error
	Stop()

	Put(entity entity.Entity) (uint, error)
	Get(entity entity.Entity, id uint) (entity.Entity, error)
}
