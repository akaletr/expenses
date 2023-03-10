package event

import (
	"cmd/main/main.go/internal/entity/category"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	CategoryID  uint
	Category    category.Category `gorm:"foreignKey:CategoryID"`
	Description string
	Sum         int
}

func (event *Event) Register(conn *gorm.DB) error {

	if !conn.Migrator().HasTable(&event) {
		err := conn.Migrator().CreateTable(&event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (event *Event) Put(conn *gorm.DB) error {
	tx := conn.Create(event)
	return tx.Error
}

func (event *Event) Get(conn *gorm.DB, id uint) error {
	return nil
}
