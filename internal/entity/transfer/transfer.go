package transfer

import (
	"cmd/main/main.go/internal/entity/subwallet"
	"cmd/main/main.go/internal/jsonrpc"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Transfer struct {
	gorm.Model
	Name     string              `json:"name"`
	Sum      float64             `json:"sum"`
	Multiply float64             `json:"multiply"`
	FromID   uint                `json:"from_id"`
	From     subwallet.SubWallet `gorm:"foreignKey:FromID"`
	ToID     uint                `json:"to_id"`
	To       subwallet.SubWallet `gorm:"foreignKey:ToID"`
}

func (sw *Transfer) Register(conn *gorm.DB) error {
	if !conn.Migrator().HasTable(&sw) {
		err := conn.Migrator().CreateTable(&sw)
		if err != nil {
			return err
		}
	}

	err := conn.Migrator().AutoMigrate(sw)
	if err != nil {
		return err
	}

	return nil
}

func Get(opt jsonrpc.Options) (json.RawMessage, error) {
	var transfer Transfer
	err := json.Unmarshal(opt.Params, &transfer)
	if err != nil {
		return nil, err
	}
	opt.Conn.First(&transfer, transfer.ID)
	return json.Marshal(transfer)
}

func GetMany(opt jsonrpc.Options) (json.RawMessage, error) {
	var transfer []Transfer
	time.Sleep(time.Millisecond * 100)
	type P struct {
		WalletID int `json:"wallet_id"`
	}

	var p P

	err := json.Unmarshal(opt.Params, &p)
	if err != nil {
		return nil, err
	}
	opt.Conn.Where("wallet_id = ?", p.WalletID).Find(&transfer)
	return json.Marshal(transfer)
}

func Create(opt jsonrpc.Options) (json.RawMessage, error) {
	var transfer Transfer
	err := json.Unmarshal(opt.Params, &transfer)
	if err != nil {
		return nil, err
	}

	opt.Conn.Create(&transfer)
	return json.Marshal(transfer.ID)
}

func Delete(opt jsonrpc.Options) (json.RawMessage, error) {
	var transfer Transfer
	err := json.Unmarshal(opt.Params, &transfer)
	if err != nil {
		return nil, err
	}

	opt.Conn.Delete(&transfer, transfer.ID)
	return json.Marshal(transfer.ID)
}
