package repository

import "gorm.io/gorm"

type Repository struct {
	Account  AccountRepository
	Entry    EntryRepository
	Transfer TransferRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Account:  NewAccountRepository(db),
		Entry:    NewEntryRepository(db),
		Transfer: NewTransferRepository(db),
	}
}
