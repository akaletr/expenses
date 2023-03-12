package wallet

import (
	"cmd/main/main.go/internal/jsonrpc"
	"encoding/json"
	"fmt"
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

	return nil
}

func (wallet *Wallet) Put(conn *gorm.DB) error {
	tx := conn.Create(wallet)
	return tx.Error
}

func (wallet *Wallet) Get(conn *gorm.DB, id uint) error {
	tx := conn.First(wallet, id)
	return tx.Error
}

func GetWallet(opt jsonrpc.Options) (json.RawMessage, error) {
	fmt.Println(opt.Params)
	return []byte{}, nil
}
