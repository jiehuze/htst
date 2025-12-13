package services

import (
	"htst/internal/app/models"
	"htst/internal/app/repositories"
)

type IUserService interface {
	CreateUser(user *models.SysUser) error
	UpdateUser(user *models.SysUser) error
	DeleteUserByID(id uint64) error
	DeleteUserByUsername(username string) error
	GetUserByID(id uint64) (*models.SysUser, error)
	GetUserByUsername(username string) (*models.SysUser, error)
	UpdateUserRole(id uint64, role *string) error
	UpdateUserDescription(id uint64, description *string) error
	GetUserList(page, pageSize int) ([]*models.SysUser, int64, error)
	Authenticate(username, password string) (*models.SysUser, error)
}

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService() IUserService {
	return &UserService{
		userRepo: repositories.NewUserRepository(),
	}
}

func (s *UserService) CreateUser(user *models.SysUser) error {
	return s.userRepo.Create(user)
}

func (s *UserService) UpdateUser(user *models.SysUser) error {
	return s.userRepo.Update(user)
}

func (s *UserService) DeleteUserByID(id uint64) error {
	return s.userRepo.DeleteByID(id)
}

func (s *UserService) DeleteUserByUsername(username string) error {
	return s.userRepo.DeleteByUsername(username)
}

func (s *UserService) GetUserByID(id uint64) (*models.SysUser, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) GetUserByUsername(username string) (*models.SysUser, error) {
	return s.userRepo.FindByUsername(username)
}

func (s *UserService) UpdateUserRole(id uint64, role *string) error {
	return s.userRepo.UpdateRole(id, role)
}

func (s *UserService) UpdateUserDescription(id uint64, description *string) error {
	return s.userRepo.UpdateDescription(id, description)
}

func (s *UserService) GetUserList(page, pageSize int) ([]*models.SysUser, int64, error) {
	var users []*models.SysUser
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := s.userRepo.Count(&total); err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := s.userRepo.List(offset, pageSize, &users); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *UserService) Authenticate(username, password string) (*models.SysUser, error) {
	return s.userRepo.Authenticate(username, password)
}
