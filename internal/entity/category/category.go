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

func (category *Category) Put(conn *gorm.DB) (uint, error) {
	tx := conn.Create(category)
	return category.ID, tx.Error
}

func (category *Category) Get(conn *gorm.DB) error {
	return nil
}
