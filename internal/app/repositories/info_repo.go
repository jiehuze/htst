package repositories

import (
	"gorm.io/gorm"
	"htst/internal/app/models"
	"htst/internal/pkg/storage"
	"time"
)

type InfoRepository interface {
	Create(info *models.Info) error
	DeleteByID(id uint64) error
	FindByID(id uint64) (*models.Info, error)
	ListByType(infoType string, offset, limit int) ([]*models.Info, error)
	CountByType(infoType string) (int64, error)
	IncrementCount(id uint64) error
	FindByMD5(md5 string) (*models.Info, error)
	ExistsByMD5(md5 string) (bool, error)
	ExistsByTitle(title string) (bool, error)
	ListByQuery(query models.InfoQuery, offset, limit int) (int64, []*models.Info, error)
}

type infoRepository struct{}

func NewInfoRepository() InfoRepository {
	return &infoRepository{}
}

func (r *infoRepository) Create(info *models.Info) error {
	return storage.DB.Create(info).Error
}

func (r *infoRepository) DeleteByID(id uint64) error {
	return storage.DB.Where("id = ?", id).Delete(&models.Info{}).Error
}

func (r *infoRepository) FindByID(id uint64) (*models.Info, error) {
	var info models.Info
	err := storage.DB.Where("id = ?", id).First(&info).Error
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (r *infoRepository) ListByType(infoType string, offset, limit int) ([]*models.Info, error) {
	var infos []*models.Info
	err := storage.DB.Where("type = ?", infoType).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&infos).Error
	return infos, err
}

func (r *infoRepository) CountByType(infoType string) (int64, error) {
	var count int64
	err := storage.DB.Model(&models.Info{}).Where("type = ?", infoType).Count(&count).Error
	return count, err
}

/**
 * 递增计数
 */
func (r *infoRepository) IncrementCount(id uint64) error {
	return storage.DB.Model(&models.Info{}).
		Where("id = ?", id).
		UpdateColumn("count", gorm.Expr("count + ?", 1)).Error
}

func (r *infoRepository) FindByMD5(md5 string) (*models.Info, error) {
	var info models.Info
	err := storage.DB.Where("md5 = ?", md5).First(&info).Error
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (r *infoRepository) ExistsByMD5(md5 string) (bool, error) {
	var count int64
	err := storage.DB.Model(&models.Info{}).Where("md5 = ?", md5).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *infoRepository) ExistsByTitle(title string) (bool, error) {
	var count int64
	err := storage.DB.Model(&models.Info{}).Where("title = ?", title).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *infoRepository) ListByQuery(query models.InfoQuery, offset, limit int) (int64, []*models.Info, error) {
	var count int64
	var infos []*models.Info
	db := storage.DB.Model(&models.Info{})

	// 添加查询条件
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}

	if query.Title != "" {
		db = db.Where("title LIKE ?", "%"+query.Title+"%")
	}

	if query.Format != "" {
		db = db.Where("format = ?", query.Format)
	}

	// 时间范围查询
	if query.StartTime != "" && query.EndTime != "" {
		startTime, err1 := time.Parse("2006-01-02", query.StartTime)
		endTime, err2 := time.Parse("2006-01-02", query.EndTime)
		if err1 == nil && err2 == nil {
			db = db.Where("created_at BETWEEN ? AND ?", startTime, endTime)
		}
	}

	err := db.Count(&count).Error

	err = db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&infos).Error
	return count, infos, err
}
