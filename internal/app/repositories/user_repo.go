package repositories

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"htst/internal/app/models"
	"htst/internal/pkg/storage"
	"htst/pkg/util"
)

type UserRepository interface {
	Create(user *models.SysUser) error
	Update(user *models.SysUser) error
	DeleteByID(id uint64) error
	DeleteByUsername(username string) error
	FindByID(id uint64) (*models.SysUser, error)
	FindByUsername(username string) (*models.SysUser, error)
	UpdateRole(id uint64, role *string) error
	UpdateDescription(id uint64, description *string) error
	List(offset, limit int, users *[]*models.SysUser) error
	Count(total *int64) error
	Authenticate(username, password string) (*models.SysUser, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(user *models.SysUser) error {
	return storage.DB.Create(user).Error
}

func (r *userRepository) Update(user *models.SysUser) error {
	// 构建更新字段的map
	updates := make(map[string]interface{})

	// 只有非空字段才会被更新
	if user.Username != "" {
		updates["username"] = user.Username
	}
	if user.Password != "" {
		updates["password"] = user.Password
	}
	if user.Roles != nil {
		updates["roles"] = user.Roles
	}
	if user.Description != nil {
		updates["description"] = user.Description
	}

	// 使用 Updates 方法只更新指定字段
	return storage.DB.Model(&models.SysUser{}).Where("id = ?", user.ID).Updates(updates).Error
}

func (r *userRepository) DeleteByID(id uint64) error {
	return storage.DB.Where("id = ?", id).Delete(&models.SysUser{}).Error
}

func (r *userRepository) DeleteByUsername(username string) error {
	return storage.DB.Where("username = ?", username).Delete(&models.SysUser{}).Error
}

func (r *userRepository) FindByID(id uint64) (*models.SysUser, error) {
	var user models.SysUser
	err := storage.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*models.SysUser, error) {
	var user models.SysUser
	err := storage.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateRole(id uint64, role *string) error {
	return storage.DB.Model(&models.SysUser{}).Where("id = ?", id).Update("role", role).Error
}

func (r *userRepository) UpdateDescription(id uint64, description *string) error {
	return storage.DB.Model(&models.SysUser{}).Where("id = ?", id).Update("description", description).Error
}

func (r *userRepository) List(offset, limit int, users *[]*models.SysUser) error {
	return storage.DB.Offset(offset).Limit(limit).Find(users).Error
}

func (r *userRepository) Count(total *int64) error {
	return storage.DB.Model(&models.SysUser{}).Count(total).Error
}

func (r *userRepository) Authenticate(username, password string) (*models.SysUser, error) {
	var user models.SysUser

	// 根据用户名查找用户
	err := storage.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	// 验证密码
	log.Infoln("password: ", password)
	log.Infoln("user.Password: ", user.Password)
	if !util.CheckPassword(password, user.Password) {
		log.Infoln("password not match")
		return nil, errors.New("password not match")
	}
	log.Infoln("password match")

	return &user, nil
}
