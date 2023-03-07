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
