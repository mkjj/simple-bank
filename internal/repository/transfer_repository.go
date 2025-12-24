package repository

import (
	"simple_bank/internal/models"

	"gorm.io/gorm"
)

type TransferRepository interface {
	Create(transfer *models.Transfer) error
	GetByID(id int64) (*models.Transfer, error)
	GetByAccountID(accountID int64, limit, offset int) ([]models.Transfer, error)
	GetByFromAccountID(fromAccountID int64, limit, offset int) ([]models.Transfer, error)
	GetByToAccountID(toAccountID int64, limit, offset int) ([]models.Transfer, error)
	List(page, pageSize int) ([]models.Transfer, error)
}

type transferRepository struct {
	db *gorm.DB
}

func NewTransferRepository(db *gorm.DB) TransferRepository {
	return &transferRepository{db: db}
}

func (r *transferRepository) Create(transfer *models.Transfer) error {
	return r.db.Create(transfer).Error
}

func (r *transferRepository) GetByID(id int64) (*models.Transfer, error) {
	var transfer models.Transfer
	err := r.db.Preload("FromAccount").Preload("ToAccount").
		First(&transfer, id).Error
	return &transfer, err
}

func (r *transferRepository) GetByAccountID(accountID int64, limit, offset int) ([]models.Transfer, error) {
	var transfers []models.Transfer
	err := r.db.Preload("FromAccount").Preload("ToAccount").
		Where("from_account_id = ? OR to_account_id = ?", accountID, accountID).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&transfers).Error
	return transfers, err
}

func (r *transferRepository) GetByFromAccountID(fromAccountID int64, limit, offset int) ([]models.Transfer, error) {
	var transfers []models.Transfer
	err := r.db.Preload("FromAccount").Preload("ToAccount").
		Where("from_account_id = ?", fromAccountID).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&transfers).Error
	return transfers, err
}

func (r *transferRepository) GetByToAccountID(toAccountID int64, limit, offset int) ([]models.Transfer, error) {
	var transfers []models.Transfer
	err := r.db.Preload("FromAccount").Preload("ToAccount").
		Where("to_account_id = ?", toAccountID).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&transfers).Error
	return transfers, err
}

func (r *transferRepository) List(page, pageSize int) ([]models.Transfer, error) {
	var transfers []models.Transfer
	offset := (page - 1) * pageSize
	err := r.db.Preload("FromAccount").Preload("ToAccount").
		Limit(pageSize).Offset(offset).
		Order("created_at DESC").
		Find(&transfers).Error
	return transfers, err
}
