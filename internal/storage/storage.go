package storage

import (
	"cmd/main/main.go/internal/config"
	"cmd/main/main.go/internal/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type storage struct {
	conn *gorm.DB
}

func New() (Storage, error) {
	return &storage{}, nil
}

func (storage *storage) Start(cfg config.Database) error {

	conn, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		return err
	}

	storage.conn = conn
	return nil
}

func (storage *storage) Provide(entities ...entity.Entity) error {
	for _, e := range entities {
		err := e.Register(storage.conn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (storage *storage) Stop() {
	return
}

func (storage *storage) Put(entity entity.Entity) (uint, error) {

	return 0, nil
}

func (storage *storage) Get(entity entity.Entity, id uint) (entity.Entity, error) {
	return entity, nil
}

func (storage *storage) GetDB() *gorm.DB {
	return storage.conn
}
