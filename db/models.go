package db

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type TransactionType string

const (
	Deposit  TransactionType = "deposit"
	Withdraw TransactionType = "withdraw"
)

// Uso de tags en struct
type WalletModel struct {
	WalletID  string `gorm:"primaryKey"`
	Balance   int64  `gorm:"not null;default:0"` // en centavos (nunca float)
	Currency  string `gorm:"type:varchar(10);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Transactions []TransactionModel `gorm:"foreignKey:WalletID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// Uso de tags en struct
type TransactionModel struct {
	ID           string          `gorm:"primaryKey"`
	WalletID     string          `gorm:"not null;index"`
	Amount       int64           `gorm:"not null"` // en centavos
	CurrencyCode string          `gorm:"not null"`
	Type         TransactionType `gorm:"type:varchar(20);not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Uso de HOOKS
func (w *WalletModel) BeforeCreate(tx *gorm.DB) error {
	log.Println("Executing BeforeCreate (walletModel)")
	return nil

}

func (w *WalletModel) AfterCreate(tx *gorm.DB) error {
	log.Println("Executing AfterCreate (walletModel)")
	return nil

}

func (w *TransactionModel) BeforeCreate(tx *gorm.DB) error {
	log.Println("Executing BeforeCreate (TransactionModel)")
	return nil

}

func (w *TransactionModel) AfterCreate(tx *gorm.DB) error {
	log.Println("Executing AfterCreate (TransactionModel)")
	return nil

}
