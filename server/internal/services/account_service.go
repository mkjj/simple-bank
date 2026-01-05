package services

import (
	"context"
	"errors"
	"simple_bank/server/internal/models"
	"simple_bank/server/internal/repositories"
)

type AccountService interface {
	CreateAccount(ctx context.Context, owner, currency string, initialBalance int64) (*models.Account, error)
	GetAccount(ctx context.Context, id int64) (*models.Account, error)
	GetAccountsByOwner(ctx context.Context, owner string, page, pageSize int) ([]models.Account, error)
	ListAccounts(ctx context.Context, page, pageSize int) ([]models.Account, error)
	UpdateAccount(ctx context.Context, id int64, balance int64) (*models.Account, error)
	DeleteAccount(ctx context.Context, id int64) error
}

type accountService struct {
	repo *repositories.Repository
}

func NewAccountService(repo *repositories.Repository) AccountService {
	return &accountService{
		repo: repo,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, owner, currency string, initialBalance int64) (*models.Account, error) {
	if owner == "" {
		return nil, errors.New("owner cannot be empty")
	}
	if currency == "" {
		currency = "USD"
	}
	if initialBalance < 0 {
		return nil, errors.New("initial balance cannot be negative")
	}

	// Check if owner already has an account with this currency
	accounts, err := s.repo.Account.GetByOwner(owner, 100, 0)
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		if account.Currency == currency {
			return nil, errors.New("owner already has an account with this currency")
		}
	}

	// Create account instance
	account := &models.Account{
		Owner:    owner,
		Balance:  initialBalance,
		Currency: currency,
	}

	// Save to database
	err = s.repo.Account.Create(account)
	if err != nil {
		return nil, err
	}

	// Create initial entry
	entry := &models.Entry{
		AccountID: account.ID, // Now account.ID is defined
		Amount:    initialBalance,
	}
	err = s.repo.Entry.Create(entry)
	if err != nil {
		// Rollback account creation
		s.repo.Account.Delete(account.ID)
		return nil, err
	}

	return account, nil
}

func (s *accountService) GetAccount(ctx context.Context, id int64) (*models.Account, error) {
	return s.repo.Account.GetByID(id)
}

func (s *accountService) GetAccountsByOwner(ctx context.Context, owner string, page, pageSize int) ([]models.Account, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	return s.repo.Account.GetByOwner(owner, pageSize, offset)
}

func (s *accountService) ListAccounts(ctx context.Context, page, pageSize int) ([]models.Account, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.repo.Account.List(page, pageSize)
}

func (s *accountService) UpdateAccount(ctx context.Context, id int64, balance int64) (*models.Account, error) {
	account, err := s.repo.Account.GetByID(id)
	if err != nil {
		return nil, err
	}

	account.Balance = balance
	err = s.repo.Account.Update(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *accountService) DeleteAccount(ctx context.Context, id int64) error {
	return s.repo.Account.Delete(id)
}
