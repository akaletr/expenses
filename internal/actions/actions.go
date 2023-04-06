package actions

import (
	"encoding/json"

	"cmd/main/main.go/internal/entity/event"
	"cmd/main/main.go/internal/entity/subwallet"
	"cmd/main/main.go/internal/entity/wallet"
	"cmd/main/main.go/internal/jsonrpc"

	"gorm.io/gorm"
)

type TransferParams struct {
	From     int     `json:"from"`
	To       int     `json:"to"`
	Sum      float64 `json:"sum"`
	Multiply float64 `json:"multiply"`
}

func Transfer(opt jsonrpc.Options) (json.RawMessage, error) {
	var params TransferParams
	err := json.Unmarshal(opt.Params, &params)
	if err != nil {
		return nil, err
	}

	var walletFrom subwallet.SubWallet
	var walletTo subwallet.SubWallet
	opt.Conn.Find(&walletFrom, params.From)
	opt.Conn.Find(&walletTo, params.To)

	opt.Conn.Model(&walletFrom).Where("id = ?", params.From).Update("sum", walletFrom.Sum-params.Sum)
	opt.Conn.Model(&walletTo).Where("id = ?", params.To).Update("sum", walletTo.Sum+params.Sum*params.Multiply)
	return json.Marshal("")
}

type EventParams struct {
	CategoryID  uint    `json:"category_id"`
	SubWalletID uint    `json:"sub_wallet_id"`
	Description string  `json:"description"`
	Sum         float64 `json:"sum"`
}

func Event(opt jsonrpc.Options) (json.RawMessage, error) {
	var params EventParams
	err := json.Unmarshal(opt.Params, &params)
	if err != nil {
		return nil, err
	}

	var subWalletFrom subwallet.SubWallet

	var w wallet.Wallet
	opt.Conn.Where("user_id = ?", opt.UserId).First(&w)

	opt.Conn.Find(&subWalletFrom, params.SubWalletID)

	e := event.Event{
		Model:       gorm.Model{},
		CategoryID:  params.CategoryID,
		UserID:      opt.UserId,
		SubWalletID: params.SubWalletID,
		Description: "",
		Sum:         params.Sum,
	}

	opt.Conn.Create(&e)
	opt.Conn.Model(&w).Where("id = ?", subWalletFrom.WalletID).Update("sum", w.Sum-params.Sum/384)
	opt.Conn.Model(&subWalletFrom).Where("id = ?", subWalletFrom.ID).Update("sum", subWalletFrom.Sum-params.Sum)

	return json.Marshal("")
}
