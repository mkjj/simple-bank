package services

import (
	"simple_bank/server/internal/repositories"

	"gorm.io/gorm"
)

type Services struct {
	Account  AccountService
	Transfer TransferService
}

func NewServices(repo *repositories.Repository, db *gorm.DB) *Services {
	return &Services{
		Account:  NewAccountService(repo),
		Transfer: NewTransferService(repo, db),
	}
}
