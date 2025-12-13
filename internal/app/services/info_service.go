package services

import (
	log "github.com/sirupsen/logrus"
	"htst/internal/app/models"
	"htst/internal/app/repositories"
	"htst/pkg/util"
	"os"
	"path/filepath"
)

type IInfoService interface {
	CreateInfo(info *models.Info) error
	DeleteInfoByID(id uint64) error
	GetInfoByID(id uint64) (*models.Info, error)
	GetInfoListByType(infoType string, page, pageSize int) ([]*models.Info, int64, error)
	IncrementInfoCount(id uint64) error
	GetInfoByMD5(md5 string) (*models.Info, error)
	CheckInfoExistsByMD5(md5 string) (bool, error)
	GetInfoListByQuery(query models.InfoQuery) (int64, []*models.Info, error)
	CleanFile(filePath string)
}

type InfoService struct {
	infoRepo repositories.InfoRepository
}

func NewInfoService() IInfoService {
	return &InfoService{
		infoRepo: repositories.NewInfoRepository(),
	}
}

func (s *InfoService) CreateInfo(info *models.Info) error {
	return s.infoRepo.Create(info)
}

func (s *InfoService) DeleteInfoByID(id uint64) error {
	return s.infoRepo.DeleteByID(id)
}

func (s *InfoService) GetInfoByID(id uint64) (*models.Info, error) {
	return s.infoRepo.FindByID(id)
}

func (s *InfoService) GetInfoListByType(infoType string, page, pageSize int) ([]*models.Info, int64, error) {
	offset := (page - 1) * pageSize

	infos, err := s.infoRepo.ListByType(infoType, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.infoRepo.CountByType(infoType)
	if err != nil {
		return nil, 0, err
	}

	return infos, count, nil
}

func (s *InfoService) IncrementInfoCount(id uint64) error {
	return s.infoRepo.IncrementCount(id)
}

func (s *InfoService) GetInfoByMD5(md5 string) (*models.Info, error) {
	return s.infoRepo.FindByMD5(md5)
}

func (s *InfoService) CheckInfoExistsByMD5(md5 string) (bool, error) {
	return s.infoRepo.ExistsByMD5(md5)
}

func (s *InfoService) GetInfoListByQuery(query models.InfoQuery) (int64, []*models.Info, error) {
	// 设置默认值
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}

	offset := (query.Page - 1) * query.PageSize

	return s.infoRepo.ListByQuery(query, offset, query.PageSize)

}

func (s *InfoService) CleanFile(filePath string) {
	entries, err := os.ReadDir(filePath)
	if err != nil {
		log.Errorf("无法读取目录 %s: %v", filePath, err)
		return
	}

	for _, entry := range entries {
		// os.ReadDir 已经自动过滤了 . 和 .. 目录
		if entry.Name() == "." || entry.Name() == ".." {
			continue
		}

		fullPath := filepath.Join(filePath, entry.Name())
		if entry.IsDir() {
			log.Debugf("跳过目录: %s", fullPath)
		} else {
			log.Debugf("处理文件: %s", fullPath)
			// 添加文件处理逻辑
			md5, _ := util.CalculateFileMD5(fullPath)
			ex, err := s.infoRepo.ExistsByMD5(md5)
			if err != nil || ex == false {
				log.Errorf("文件 %s 未使用，删除该文件", fullPath)
				//_ = os.Remove(fullPath)
				continue
			}
		}
	}
}
