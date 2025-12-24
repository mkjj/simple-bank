package repository

import (
	"simple_bank/internal/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountRepository interface {
	Create(account *models.Account) error
	GetByID(id int64) (*models.Account, error)
	GetByOwner(owner string, limit, offset int) ([]models.Account, error)
	List(page, pageSize int) ([]models.Account, error)
	Update(account *models.Account) error
	Delete(id int64) error
	UpdateBalance(id int64, amount int64) error
	GetForUpdate(id int64) (*models.Account, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

// Create a new account
func (r *accountRepository) Create(account *models.Account) error {
	// Set created_at if not set
	if account.CreatedAt.IsZero() {
		account.CreatedAt = time.Now()
	}
	return r.db.Create(account).Error
}

// Get account by ID
func (r *accountRepository) GetByID(id int64) (*models.Account, error) {
	var account models.Account
	err := r.db.First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// Get accounts by owner with pagination
func (r *accountRepository) GetByOwner(owner string, limit, offset int) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Where("owner = ?", owner).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&accounts).Error
	return accounts, err
}

// List accounts with pagination
func (r *accountRepository) List(page, pageSize int) ([]models.Account, error) {
	var accounts []models.Account
	offset := (page - 1) * pageSize
	err := r.db.Limit(pageSize).Offset(offset).
		Order("created_at DESC").
		Find(&accounts).Error
	return accounts, err
}

// Update account
func (r *accountRepository) Update(account *models.Account) error {
	return r.db.Save(account).Error
}

// Delete account
func (r *accountRepository) Delete(id int64) error {
	return r.db.Delete(&models.Account{}, id).Error
}

// Update account balance (atomic update)
func (r *accountRepository) UpdateBalance(id int64, amount int64) error {
	return r.db.Model(&models.Account{}).
		Where("id = ?", id).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

// Get account for update (with row lock)
func (r *accountRepository) GetForUpdate(id int64) (*models.Account, error) {
	var account models.Account
	err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}
