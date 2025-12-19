package models

import (
	"time"
)

// DoctorInfo 医生信息实体
type DoctorInfo struct {
	ID                   uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name                 string    `gorm:"column:name;size:50;not null" json:"name"`
	Title                string    `gorm:"column:title;size:50;not null" json:"title"`
	Department           string    `gorm:"column:department;size:100;not null" json:"department"`
	Degree               string    `gorm:"column:degree;size:50" json:"degree,omitempty"`
	AvatarURL            string    `gorm:"column:avatar_url;size:255" json:"avatar_url,omitempty"`
	Introduction         string    `gorm:"column:introduction;size:500" json:"introduction,omitempty"`
	DetailedDescription  string    `gorm:"column:detailed_description" json:"detailed_description,omitempty"`
	Specialties          string    `gorm:"column:specialties;size:500" json:"specialties,omitempty"`
	ResearchAchievements string    `gorm:"column:research_achievements" json:"research_achievements,omitempty"`
	CreatedAt            time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (DoctorInfo) TableName() string {
	return "doctor_info"
}

// DoctorInfoQuery 查询参数结构体
type DoctorInfoQuery struct {
	Name       string `form:"name"`
	Department string `form:"department"`
	Title      string `form:"title"`
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
}
