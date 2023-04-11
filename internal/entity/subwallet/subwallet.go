package subwallet

import (
	"cmd/main/main.go/internal/entity/wallet"
	"encoding/json"
	"time"

	"cmd/main/main.go/internal/jsonrpc"

	"gorm.io/gorm"
)

type Currency int

const (
	USD Currency = iota
	ARS
	RUB
)

type SubWallet struct {
	gorm.Model
	Name     string        `json:"name"`
	Currency Currency      `json:"currency"`
	Sum      float64       `json:"sum"`
	WalletID uint          `json:"wallet_id"`
	Wallet   wallet.Wallet `gorm:"foreignKey:WalletID"`
}

func (sw *SubWallet) Register(conn *gorm.DB) error {
	if !conn.Migrator().HasTable(&sw) {
		err := conn.Migrator().CreateTable(&sw)
		if err != nil {
			return err
		}
		sw1 := SubWallet{
			Model:    gorm.Model{},
			Name:     "USD",
			Currency: 0,
			Sum:      1000,
			WalletID: 1,
		}

		sw2 := SubWallet{
			Model:    gorm.Model{},
			Name:     "ARS",
			Currency: 1,
			Sum:      0,
			WalletID: 1,
		}

		sw3 := SubWallet{
			Model:    gorm.Model{},
			Name:     "RUB",
			Currency: 1,
			Sum:      0,
			WalletID: 1,
		}

		conn.Create(&sw1)
		conn.Create(&sw2)
		conn.Create(&sw3)
	}

	err := conn.Migrator().AutoMigrate(sw)
	if err != nil {
		return err
	}

	return nil
}

func Get(opt jsonrpc.Options) (json.RawMessage, error) {
	var sw SubWallet
	err := json.Unmarshal(opt.Params, &sw)
	if err != nil {
		return nil, err
	}
	opt.Conn.First(&sw, sw.ID)
	return json.Marshal(sw)
}

func GetMany(opt jsonrpc.Options) (json.RawMessage, error) {
	var sw []SubWallet
	time.Sleep(time.Millisecond * 100)
	type P struct {
		WalletID int `json:"wallet_id"`
	}

	var p P

	err := json.Unmarshal(opt.Params, &p)
	if err != nil {
		return nil, err
	}
	opt.Conn.Where("wallet_id = ?", p.WalletID).Find(&sw)
	return json.Marshal(sw)
}

func Create(opt jsonrpc.Options) (json.RawMessage, error) {
	var sw SubWallet
	err := json.Unmarshal(opt.Params, &sw)
	if err != nil {
		return nil, err
	}

	opt.Conn.Create(&sw)
	return json.Marshal(sw.ID)
}

func Delete(opt jsonrpc.Options) (json.RawMessage, error) {
	var sw SubWallet
	err := json.Unmarshal(opt.Params, &sw)
	if err != nil {
		return nil, err
	}

	opt.Conn.Delete(&sw, sw.ID)
	return json.Marshal(sw.ID)
}
