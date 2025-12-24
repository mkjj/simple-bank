package models

import (
	"time"

	"gorm.io/gorm"
)

type Transfer struct {
	ID            int64     `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	FromAccountID int64     `gorm:"type:bigint;not null;index" json:"from_account_id"`
	ToAccountID   int64     `gorm:"type:bigint;not null;index" json:"to_account_id"`
	Amount        int64     `gorm:"type:bigint;not null" json:"amount"` // Must be positive
	CreatedAt     time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	// Define composite index
	FromAccount Account `gorm:"foreignKey:FromAccountID" json:"from_account,omitempty"`
	ToAccount   Account `gorm:"foreignKey:ToAccountID" json:"to_account,omitempty"`
}

// TableName specifies the table name for GORM
func (Transfer) TableName() string {
	return "transfers"
}

// BeforeCreate validates the transfer before creation
func (t *Transfer) BeforeCreate(tx *gorm.DB) error {
	if t.Amount <= 0 {
		return gorm.ErrInvalidValue
	}
	if t.FromAccountID == t.ToAccountID {
		return gorm.ErrInvalidValue
	}
	return nil
}
