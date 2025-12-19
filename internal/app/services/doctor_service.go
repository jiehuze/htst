package services

import (
	"htst/internal/app/models"
	"htst/internal/app/repositories"
)

type IDoctorInfoService interface {
	CreateDoctorInfo(doctor *models.DoctorInfo) (*models.DoctorInfo, error)
	GetDoctorInfoByID(id uint64) (*models.DoctorInfo, error)
	UpdateDoctorInfo(id uint64, input *models.DoctorInfo) error
	DeleteDoctorInfo(id uint64) error
	ListDoctorInfos(query *models.DoctorInfoQuery) (int64, []*models.DoctorInfo, error)
}

type DoctorInfoService struct {
	repo repositories.DoctorInfoRepository
}

func NewDoctorInfoService() IDoctorInfoService {
	return &DoctorInfoService{
		repo: repositories.NewDoctorInfoRepository(),
	}
}

//type CreateDoctorInfoRequest struct {
//	Name                 string  `json:"name" binding:"required"`
//	Title                string  `json:"title" binding:"required"`
//	Department           string  `json:"department" binding:"required"`
//	Degree               *string `json:"degree,omitempty"`
//	AvatarURL            *string `json:"avatar_url,omitempty"`
//	Introduction         *string `json:"introduction,omitempty"`
//	DetailedDescription  *string `json:"detailed_description,omitempty"`
//	Specialties          *string `json:"specialties,omitempty"`
//	ResearchAchievements *string `json:"research_achievements,omitempty"`
//}
//
//type UpdateDoctorInfoRequest struct {
//	Name                 *string `json:"name,omitempty"`
//	Title                *string `json:"title,omitempty"`
//	Department           *string `json:"department,omitempty"`
//	Degree               *string `json:"degree,omitempty"`
//	AvatarURL            *string `json:"avatar_url,omitempty"`
//	Introduction         *string `json:"introduction,omitempty"`
//	DetailedDescription  *string `json:"detailed_description,omitempty"`
//	Specialties          *string `json:"specialties,omitempty"`
//	ResearchAchievements *string `json:"research_achievements,omitempty"`
//}

func (s *DoctorInfoService) CreateDoctorInfo(doctor *models.DoctorInfo) (*models.DoctorInfo, error) {
	err := s.repo.Create(doctor)
	if err != nil {
		return nil, err
	}
	return doctor, nil
}

func (s *DoctorInfoService) GetDoctorInfoByID(id uint64) (*models.DoctorInfo, error) {
	return s.repo.FindByID(id)
}

func (s *DoctorInfoService) UpdateDoctorInfo(id uint64, input *models.DoctorInfo) error {
	// 构建更新字段映射
	updates := make(map[string]interface{})

	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.Department != "" {
		updates["department"] = input.Department
	}
	if input.Degree != "" {
		updates["degree"] = input.Degree
	}
	if input.AvatarURL != "" {
		updates["avatar_url"] = input.AvatarURL
	}
	if input.Introduction != "" {
		updates["introduction"] = input.Introduction
	}
	if input.DetailedDescription != "" {
		updates["detailed_description"] = input.DetailedDescription
	}
	if input.Specialties != "" {
		updates["specialties"] = input.Specialties
	}
	if input.ResearchAchievements != "" {
		updates["research_achievements"] = input.ResearchAchievements
	}

	// 执行更新
	doctor := &models.DoctorInfo{}
	for k, v := range updates {
		switch k {
		case "name":
			if name, ok := v.(string); ok {
				doctor.Name = name
			}
		case "title":
			if title, ok := v.(string); ok {
				doctor.Title = title
			}
		case "department":
			if department, ok := v.(string); ok {
				doctor.Department = department
			}
		default:
			// 其他字段保持原样处理
		}
	}

	return s.repo.Update(id, doctor)
}

func (s *DoctorInfoService) DeleteDoctorInfo(id uint64) error {
	return s.repo.DeleteByID(id)
}

func (s *DoctorInfoService) ListDoctorInfos(query *models.DoctorInfoQuery) (int64, []*models.DoctorInfo, error) {
	// 计算分页参数
	page := query.Page
	pageSize := query.PageSize

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.repo.ListByQuery(*query, offset, pageSize)
}
