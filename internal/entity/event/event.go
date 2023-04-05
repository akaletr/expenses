package event

import (
	"cmd/main/main.go/internal/entity/wallet"
	"encoding/json"
	"fmt"

	"cmd/main/main.go/internal/entity/category"
	"cmd/main/main.go/internal/entity/user"
	"cmd/main/main.go/internal/jsonrpc"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	CategoryID  uint              `json:"category_id"`
	Category    category.Category `gorm:"foreignKey:CategoryID"`
	UserID      uint              `json:"user_id"`
	User        user.User         `gorm:"foreignKey:UserID"`
	Description string            `json:"description"`
	Sum         int               `json:"sum"`
}

func (event *Event) Register(conn *gorm.DB) error {
	if !conn.Migrator().HasTable(&event) {
		err := conn.Migrator().CreateTable(&event)
		if err != nil {
			return err
		}
		return nil
	}

	err := conn.Migrator().AutoMigrate(event)
	if err != nil {
		return err
	}
	return nil
}

func Get(opt jsonrpc.Options) (json.RawMessage, error) {
	var event Event
	opt.Conn.Where("user_id = ?", opt.UserId).First(&event, 1)
	return json.Marshal(event)
}

func GetMany(opt jsonrpc.Options) (json.RawMessage, error) {
	var event []Event

	opt.Conn.Where("user_id = ?", opt.UserId).Find(&event)
	return json.Marshal(event)
}

func Create(opt jsonrpc.Options) (json.RawMessage, error) {
	var event Event
	err := json.Unmarshal(opt.Params, &event)
	if err != nil {
		return nil, err
	}

	var w wallet.Wallet
	opt.Conn.First(&w, opt.UserId)

	event.UserID = opt.UserId
	fmt.Println(event.Sum)
	opt.Conn.Create(&event)

	s := w.Sum - event.Sum
	opt.Conn.Model(&w).Where("id = ?", opt.UserId).Update("sum", s)

	return json.Marshal(event.ID)
}

func Delete(opt jsonrpc.Options) (json.RawMessage, error) {
	var event Event
	err := json.Unmarshal(opt.Params, &event)
	if err != nil {
		return nil, err
	}

	opt.Conn.Delete(&event, event.ID)
	return json.Marshal(event.ID)
}
