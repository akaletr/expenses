package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	WalletID  uint   `json:"wallet_id"`
	Wallet    Wallet `gorm:"foreignKey:WalletID"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (user *User) Register(conn *gorm.DB) error {

	if !conn.Migrator().HasTable(&user) {
		err := conn.Migrator().CreateTable(&user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (user *User) Put(conn *gorm.DB) error {
	tx := conn.Create(user)
	return tx.Error
}

func (user *User) Get(conn *gorm.DB, id uint) error {
	tx := conn.First(user, id)
	return tx.Error
}
