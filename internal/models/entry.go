package models

import (
	"time"
)

type Entry struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	AccountID int64     `gorm:"type:bigint;not null;index" json:"account_id"`
	Amount    int64     `gorm:"type:bigint;not null" json:"amount"` // Can be negative or positive
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	Account   Account   `gorm:"foreignKey:AccountID" json:"account,omitempty"`
}

// TableName specifies the table name for GORM
func (Entry) TableName() string {
	return "entries"
}
