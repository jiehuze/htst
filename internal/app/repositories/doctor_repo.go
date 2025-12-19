package repositories

import (
	"htst/internal/app/models"
	"htst/internal/pkg/storage"
)

type DoctorInfoRepository interface {
	Create(doctor *models.DoctorInfo) error
	DeleteByID(id uint64) error
	FindByID(id uint64) (*models.DoctorInfo, error)
	Update(id uint64, doctor *models.DoctorInfo) error
	ListByQuery(query models.DoctorInfoQuery, offset, limit int) (int64, []*models.DoctorInfo, error)
	ExistsByName(name string) (bool, error)
}

type doctorInfoRepository struct{}

func NewDoctorInfoRepository() DoctorInfoRepository {
	return &doctorInfoRepository{}
}

func (r *doctorInfoRepository) Create(doctor *models.DoctorInfo) error {
	return storage.DB.Create(doctor).Error
}

func (r *doctorInfoRepository) DeleteByID(id uint64) error {
	return storage.DB.Where("id = ?", id).Delete(&models.DoctorInfo{}).Error
}

func (r *doctorInfoRepository) FindByID(id uint64) (*models.DoctorInfo, error) {
	var doctor models.DoctorInfo
	err := storage.DB.Where("id = ?", id).First(&doctor).Error
	if err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (r *doctorInfoRepository) Update(id uint64, doctor *models.DoctorInfo) error {
	return storage.DB.Model(&models.DoctorInfo{}).
		Where("id = ?", id).
		Updates(doctor).Error
}

func (r *doctorInfoRepository) ListByQuery(query models.DoctorInfoQuery, offset, limit int) (int64, []*models.DoctorInfo, error) {
	var count int64
	var doctors []*models.DoctorInfo
	db := storage.DB.Model(&models.DoctorInfo{})

	// 添加查询条件
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}

	if query.Department != "" {
		db = db.Where("department = ?", query.Department)
	}

	if query.Title != "" {
		db = db.Where("title = ?", query.Title)
	}

	err := db.Count(&count).Error
	if err != nil {
		return 0, nil, err
	}

	err = db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&doctors).Error
	return count, doctors, err
}

func (r *doctorInfoRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := storage.DB.Model(&models.DoctorInfo{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
