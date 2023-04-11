package wallet

import (
	"encoding/json"
	"time"

	"cmd/main/main.go/internal/entity/user"
	"cmd/main/main.go/internal/jsonrpc"

	"gorm.io/gorm"
)

type Currency int

const (
	USD Currency = iota
	ARS
	RUB
)

type Wallet struct {
	gorm.Model
	Name     string    `json:"name"`
	Currency Currency  `json:"currency"`
	Sum      float64   `json:"sum"`
	UserID   uint      `json:"user_id"`
	User     user.User `gorm:"foreignKey:UserID"`
}

func (wallet *Wallet) Register(conn *gorm.DB) error {
	if !conn.Migrator().HasTable(&wallet) {
		err := conn.Migrator().CreateTable(&wallet)
		if err != nil {
			return err
		}

		w1 := Wallet{
			Model:    gorm.Model{},
			Name:     "Main",
			Currency: 0,
			Sum:      1000,
			UserID:   1,
		}

		conn.Create(&w1)
	}

	err := conn.Migrator().AutoMigrate(wallet)
	if err != nil {
		return err
	}

	return nil
}

func Get(opt jsonrpc.Options) (json.RawMessage, error) {
	time.Sleep(time.Millisecond * 100)
	var wallet Wallet
	err := json.Unmarshal(opt.Params, &wallet)
	if err != nil {
		return nil, err
	}
	opt.Conn.Where("user_id = ?", 1).First(&wallet)
	return json.Marshal(wallet)
}

func GetMany(opt jsonrpc.Options) (json.RawMessage, error) {
	var wallet []Wallet
	opt.Conn.Find(&wallet)
	return json.Marshal(wallet)
}

func Create(opt jsonrpc.Options) (json.RawMessage, error) {
	var wallet Wallet
	err := json.Unmarshal(opt.Params, &wallet)
	if err != nil {
		return nil, err
	}

	opt.Conn.Create(&wallet)
	return json.Marshal(wallet.ID)
}

func Delete(opt jsonrpc.Options) (json.RawMessage, error) {
	var wallet Wallet
	err := json.Unmarshal(opt.Params, &wallet)
	if err != nil {
		return nil, err
	}

	opt.Conn.Delete(&wallet, wallet.ID)
	return json.Marshal(wallet.ID)
}
