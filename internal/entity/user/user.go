package user

import (
	"cmd/main/main.go/internal/entity/wallet"
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
	}

	err := conn.Migrator().AutoMigrate(user)
	if err != nil {
		return err
	}
	return nil
}
