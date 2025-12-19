package repositories

import (
	"htst/internal/app/models"
	"htst/internal/pkg/storage"
)

type ContactRepository interface {
	Create(contact *models.Contact) error
	Update(id uint64, contact *models.Contact) error
	DeleteByID(id uint64) error
	FindByID(id uint64) (*models.Contact, error)
	One() (*models.Contact, error)
	ExistsByEmail(email string) (bool, error)
}

type contactRepository struct{}

func NewContactRepository() ContactRepository {
	return &contactRepository{}
}

func (r *contactRepository) Create(contact *models.Contact) error {
	return storage.DB.Create(contact).Error
}

func (r *contactRepository) Update(id uint64, contact *models.Contact) error {
	return storage.DB.Model(&models.Contact{}).
		Where("id = ?", id).
		Updates(contact).Error
}

func (r *contactRepository) DeleteByID(id uint64) error {
	return storage.DB.Where("id = ?", id).Delete(&models.Contact{}).Error
}

func (r *contactRepository) FindByID(id uint64) (*models.Contact, error) {
	var contact models.Contact
	err := storage.DB.Where("id = ?", id).First(&contact).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

func (r *contactRepository) One() (*models.Contact, error) {
	var contact models.Contact
	db := storage.DB.Model(&models.Contact{})

	err := db.First(&contact).Error
	if err != nil {
		return nil, err
	}

	return &contact, nil
}

func (r *contactRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := storage.DB.Model(&models.Contact{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
