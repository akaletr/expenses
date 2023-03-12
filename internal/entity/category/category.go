package category

import (
	"cmd/main/main.go/internal/entity/user"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      uint      `json:"user_id"`
	User        user.User `gorm:"foreignKey:UserID""`
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
