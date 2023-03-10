package category

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Title       string
	Description string
}

func (category *Category) Register(conn *gorm.DB) error {

	if !conn.Migrator().HasTable(&category) {
		err := conn.Migrator().CreateTable(&category)
		if err != nil {
			return err
		}
	}

	return nil
}

func (category *Category) Put(conn *gorm.DB) error {
	tx := conn.Create(category)
	return tx.Error
}

func (category *Category) Get(conn *gorm.DB, id uint) error {
	tx := conn.First(category, id)
	return tx.Error
}
