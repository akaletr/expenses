package user

import (
	"encoding/json"
	"fmt"

	"cmd/main/main.go/internal/entity/wallet"
	"cmd/main/main.go/internal/jsonrpc"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	WalletID  uint          `json:"wallet_id"`
	Wallet    wallet.Wallet `gorm:"foreignKey:WalletID"`
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
	Email     string        `json:"email" gorm:"unique"`
}

func (user *User) Register(conn *gorm.DB) error {
	if !conn.Migrator().HasTable(&user) {
		err := conn.Migrator().CreateTable(&user)
		if err != nil {
			return err
		}

		u := User{
			Model:     gorm.Model{},
			WalletID:  1,
			FirstName: "Dmitrii",
			LastName:  "Poliakov",
			Email:     "akaletr@gmail.com",
		}

		conn.Create(&u)
	}

	err := conn.Migrator().AutoMigrate(user)
	if err != nil {
		return err
	}

	return nil
}

func Get(opt jsonrpc.Options) (json.RawMessage, error) {
	var user User
	opt.Conn.First(&user, opt.UserId)
	fmt.Println(opt.Params)
	return json.Marshal(user)
}

func GetMany(opt jsonrpc.Options) (json.RawMessage, error) {
	var user []User
	opt.Conn.Find(&user)
	return json.Marshal(user)
}

func Create(opt jsonrpc.Options) (json.RawMessage, error) {
	var user User
	err := json.Unmarshal(opt.Params, &user)
	if err != nil {
		return nil, err
	}

	opt.Conn.Create(&user)
	return json.Marshal(user.ID)
}

func Delete(opt jsonrpc.Options) (json.RawMessage, error) {
	var user User
	err := json.Unmarshal(opt.Params, &user)
	if err != nil {
		return nil, err
	}

	opt.Conn.Delete(&user, user.ID)
	return json.Marshal(user.ID)
}
