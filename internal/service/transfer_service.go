package service

import (
	"context"
	"errors"
	"fmt"
	"simple_bank/internal/models"
	"simple_bank/internal/repository"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransferService interface {
	CreateTransfer(ctx context.Context, fromAccountID, toAccountID, amount int64) (*models.Transfer, error)
	GetTransfer(ctx context.Context, id int64) (*models.Transfer, error)
	ListTransfers(ctx context.Context, accountID int64, page, pageSize int) ([]models.Transfer, error)
}

type transferService struct {
	repo *repository.Repository
	db   *gorm.DB
}

func NewTransferService(repo *repository.Repository, db *gorm.DB) TransferService {
	return &transferService{
		repo: repo,
		db:   db,
	}
}

// CreateTransfer performs a money transfer between two accounts
func (s *transferService) CreateTransfer(ctx context.Context, fromAccountID, toAccountID, amount int64) (*models.Transfer, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	if fromAccountID == toAccountID {
		return nil, errors.New("cannot transfer to the same account")
	}

	var result *models.Transfer

	// Use transaction to ensure data consistency
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Get both accounts with FOR UPDATE lock
		var fromAccount, toAccount models.Account

		// Lock from account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&fromAccount, fromAccountID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("from account not found")
			}
			return err
		}

		// Lock to account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&toAccount, toAccountID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("to account not found")
			}
			return err
		}

		// Check if accounts have the same currency
		if fromAccount.Currency != toAccount.Currency {
			return errors.New("accounts have different currencies")
		}

		// Check if from account has sufficient balance
		if fromAccount.Balance < amount {
			return errors.New("insufficient balance")
		}

		// Update balances
		if err := tx.Model(&fromAccount).
			Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		if err := tx.Model(&toAccount).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		// Create transfer record
		transfer := &models.Transfer{
			FromAccountID: fromAccountID,
			ToAccountID:   toAccountID,
			Amount:        amount,
			CreatedAt:     time.Now(),
		}

		if err := tx.Create(transfer).Error; err != nil {
			return err
		}

		// Create entry for from account (negative amount)
		fromEntry := &models.Entry{
			AccountID: fromAccountID,
			Amount:    -amount,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(fromEntry).Error; err != nil {
			return err
		}

		// Create entry for to account (positive amount)
		toEntry := &models.Entry{
			AccountID: toAccountID,
			Amount:    amount,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(toEntry).Error; err != nil {
			return err
		}

		result = transfer
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *transferService) GetTransfer(ctx context.Context, id int64) (*models.Transfer, error) {
	return s.repo.Transfer.GetByID(id)
}

func (s *transferService) ListTransfers(ctx context.Context, accountID int64, page, pageSize int) ([]models.Transfer, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.repo.Transfer.GetByAccountID(accountID, pageSize, offset)
}
