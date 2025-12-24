package repository

import (
	"simple_bank/internal/models"

	"gorm.io/gorm"
)

type EntryRepository interface {
	Create(entry *models.Entry) error
	GetByID(id int64) (*models.Entry, error)
	GetByAccountID(accountID int64, limit, offset int) ([]models.Entry, error)
	List(page, pageSize int) ([]models.Entry, error)
}

type entryRepository struct {
	db *gorm.DB
}

func NewEntryRepository(db *gorm.DB) EntryRepository {
	return &entryRepository{db: db}
}

func (r *entryRepository) Create(entry *models.Entry) error {
	return r.db.Create(entry).Error
}

func (r *entryRepository) GetByID(id int64) (*models.Entry, error) {
	var entry models.Entry
	err := r.db.Preload("Account").First(&entry, id).Error
	return &entry, err
}

func (r *entryRepository) GetByAccountID(accountID int64, limit, offset int) ([]models.Entry, error) {
	var entries []models.Entry
	err := r.db.Preload("Account").
		Where("account_id = ?", accountID).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&entries).Error
	return entries, err
}

func (r *entryRepository) List(page, pageSize int) ([]models.Entry, error) {
	var entries []models.Entry
	offset := (page - 1) * pageSize
	err := r.db.Preload("Account").
		Limit(pageSize).Offset(offset).
		Order("created_at DESC").
		Find(&entries).Error
	return entries, err
}