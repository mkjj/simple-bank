package models

import (
	"time"
)

type Account struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Owner     string    `gorm:"type:varchar;not null;index" json:"owner"`
	Balance   int64     `gorm:"type:bigint;not null;default:0" json:"balance"`
	Currency  string    `gorm:"type:varchar;not null" json:"currency"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	Entries   []Entry   `gorm:"foreignKey:AccountID" json:"entries,omitempty"`
	// FromTransfers are transfers where this account is the sender
	FromTransfers []Transfer `gorm:"foreignKey:FromAccountID" json:"from_transfers,omitempty"`
	// ToTransfers are transfers where this account is the receiver
	ToTransfers []Transfer `gorm:"foreignKey:ToAccountID" json:"to_transfers,omitempty"`
}

// TableName specifies the table name for GORM
func (Account) TableName() string {
	return "accounts"
}
