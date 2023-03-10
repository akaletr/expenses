package wallet

import (
	"encoding/json"

	"cmd/main/main.go/internal/jsonrpc"

	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	Name string `json:"name"`
	Sum  int    `json:"sum"`
}

func (wallet *Wallet) Register(conn *gorm.DB) error {
	if !conn.Migrator().HasTable(&wallet) {
		err := conn.Migrator().CreateTable(&wallet)
		if err != nil {
			return err
		}
	}

	err := conn.Migrator().AutoMigrate(wallet)
	if err != nil {
		return err
	}
	return nil
}

func Get(opt jsonrpc.Options) (json.RawMessage, error) {
	var wallet Wallet
	opt.Conn.First(&wallet, 1)
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
