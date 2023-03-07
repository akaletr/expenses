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

		//conn.Migrator().CreateConstraint(&category, "category_fk")
	}

	return nil
}
