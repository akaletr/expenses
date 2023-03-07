package storage

import (
	"cmd/main/main.go/internal/config"
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

func (storage *storage) Provide(entities ...Entity) error {
	for _, entity := range entities {
		err := entity.Register(storage.conn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (storage *storage) Put(entity Entity) (uint, error) {

	return 0, nil
}

func (storage *storage) Get(entity Entity, id uint) (Entity, error) {
	return entity, nil
}
