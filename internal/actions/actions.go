package actions

import (
	"cmd/main/main.go/internal/entity/event"
	"cmd/main/main.go/internal/entity/subwallet"
	"cmd/main/main.go/internal/entity/transfer"
	"cmd/main/main.go/internal/entity/wallet"
	"cmd/main/main.go/internal/jsonrpc"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

type TransferParams struct {
	From     uint    `json:"from"`
	To       uint    `json:"to"`
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

	t := transfer.Transfer{
		Name:     "transfer",
		Sum:      params.Sum * params.Multiply,
		Multiply: params.Multiply,
		FromID:   params.From,
		ToID:     params.To,
	}

	opt.Conn.Create(&t)
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

	var w wallet.Wallet
	opt.Conn.Where("user_id = ?", opt.UserId).First(&w)

	var subWalletFrom subwallet.SubWallet
	opt.Conn.Find(&subWalletFrom, params.SubWalletID)

	var transfers []transfer.Transfer
	opt.Conn.Where("to_id = ?", params.SubWalletID).Find(&transfers)

	sum := params.Sum
	walletSum := w.Sum
	subWalletSum := subWalletFrom.Sum

	fmt.Println(params.Sum)

	for _, tr := range transfers {
		switch {
		case sum == 0:
			break
		case sum < tr.Sum:
			walletSum = walletSum - sum/tr.Multiply
			subWalletSum = subWalletSum - sum

			opt.Conn.Model(&tr).Where("id = ?", tr.ID).Update("sum", tr.Sum-sum)
			break
		case sum >= tr.Sum:
			walletSum = walletSum - tr.Sum/tr.Multiply
			subWalletSum = subWalletSum - tr.Sum
			sum = sum - tr.Sum
			opt.Conn.Model(&tr).Where("id = ?", tr.ID).Update("sum", 0)
			opt.Conn.Delete(&tr)
		}
	}

	e := event.Event{
		Model:       gorm.Model{},
		CategoryID:  params.CategoryID,
		UserID:      opt.UserId,
		SubWalletID: params.SubWalletID,
		Description: "",
		Sum:         params.Sum,
	}

	opt.Conn.Create(&e)
	opt.Conn.Model(&w).Where("id = ?", subWalletFrom.WalletID).Update("sum", walletSum)
	opt.Conn.Model(&subWalletFrom).Where("id = ?", subWalletFrom.ID).Update("sum", subWalletSum)

	return json.Marshal("")
}
