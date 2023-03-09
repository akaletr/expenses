package event

import (
	"gorm.io/gorm"

	"cmd/main/main.go/internal/entity/category"
)

type Event struct {
	gorm.Model
	CategoryID  uint
	Category    category.Category `gorm:"foreignKey:CategoryID"`
	Description string
}

type Test struct {
	gorm.Model
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

func (event *Event) Put(conn *gorm.DB) (uint, error) {
	tx := conn.Create(event)
	return event.ID, tx.Error
}

func (event *Event) Get(conn *gorm.DB) error {
	return nil
}
