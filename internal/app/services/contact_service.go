package services

import (
	"htst/internal/app/models"
	"htst/internal/app/repositories"
)

type IContactService interface {
	CreateContact(contact *models.Contact) (*models.Contact, error)
	GetContactByID(id uint64) (*models.Contact, error)
	UpdateContact(id uint64, input *models.Contact) error
	DeleteContact(id uint64) error
	FirstContact() (*models.Contact, error)
}

type ContactService struct {
	repo repositories.ContactRepository
}

func NewContactService() IContactService {
	return &ContactService{
		repo: repositories.NewContactRepository(),
	}
}

func (s *ContactService) CreateContact(contact *models.Contact) (*models.Contact, error) {
	err := s.repo.Create(contact)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (s *ContactService) GetContactByID(id uint64) (*models.Contact, error) {
	return s.repo.FindByID(id)
}

func (s *ContactService) UpdateContact(id uint64, input *models.Contact) error {
	// 构建更新字段映射
	updates := make(map[string]interface{})

	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Phone != "" {
		updates["phone"] = input.Phone
	}
	if input.Email != "" {
		updates["email"] = input.Email
	}
	if input.Address != "" {
		updates["address"] = input.Address
	}

	// 执行更新
	contact := &models.Contact{}
	for k, v := range updates {
		switch k {
		case "name":
			if name, ok := v.(string); ok {
				contact.Name = name
			}
		case "phone":
			if phone, ok := v.(string); ok {
				contact.Phone = phone
			}
		case "email":
			if email, ok := v.(string); ok {
				contact.Email = email
			}
		case "address":
			if address, ok := v.(string); ok {
				contact.Address = address
			}
		default:
			// 其他字段保持原样处理
		}
	}

	return s.repo.Update(id, contact)
}

func (s *ContactService) DeleteContact(id uint64) error {
	return s.repo.DeleteByID(id)
}

func (s *ContactService) FirstContact() (*models.Contact, error) {
	return s.repo.One()
}
