package category

import (
	"encoding/json"

	"cmd/main/main.go/internal/entity/user"
	"cmd/main/main.go/internal/jsonrpc"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      uint      `json:"user_id"`
	User        user.User `gorm:"foreignKey:UserID"`
}

func (category *Category) Register(conn *gorm.DB) error {
	if !conn.Migrator().HasTable(&category) {
		err := conn.Migrator().CreateTable(&category)
		if err != nil {
			return err
		}
		c1 := Category{
			Model:       gorm.Model{},
			Title:       "Продукты",
			Description: "Магазин, киоск",
			UserID:      1,
		}
		c2 := Category{
			Model:       gorm.Model{},
			Title:       "Проезд",
			Description: "Такси, Общественный транспорт",
			UserID:      1,
		}

		conn.Create(&c1)
		conn.Create(&c2)
	}

	err := conn.Migrator().AutoMigrate(category)
	if err != nil {
		return err
	}

	return nil
}

func Get(opt jsonrpc.Options) (json.RawMessage, error) {
	var category Category
	opt.Conn.First(&category, 1)
	return json.Marshal(category)
}

func GetMany(opt jsonrpc.Options) (json.RawMessage, error) {
	var category []Category
	opt.Conn.Where("user_id = ?", opt.UserId).Find(&category)
	return json.Marshal(category)
}

func Create(opt jsonrpc.Options) (json.RawMessage, error) {
	var category Category
	err := json.Unmarshal(opt.Params, &category)
	if err != nil {
		return nil, err
	}

	category.UserID = opt.UserId
	opt.Conn.Create(&category)
	return json.Marshal(category.ID)
}

func Delete(opt jsonrpc.Options) (json.RawMessage, error) {
	var category Category
	err := json.Unmarshal(opt.Params, &category)
	if err != nil {
		return nil, err
	}

	opt.Conn.Delete(&category, category.ID)
	return json.Marshal(category.ID)
}
