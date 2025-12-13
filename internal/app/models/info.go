package models

import (
	"time"
)

// InfoQuery 查询参数结构体
type InfoQuery struct {
	Type      string `form:"type"`
	Title     string `form:"title"`
	Format    string `form:"format"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
}

type Info struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Type            string    `gorm:"size:50;not null;index" json:"type"`
	Title           string    `gorm:"size:255;not null" json:"title"`
	Format          string    `gorm:"size:20;not null" json:"format"`
	FilePath        string    `gorm:"size:512;not null" json:"file_path"`
	Md5             string    `gorm:"size:32;not null;uniqueIndex" json:"md5"`
	Size            uint64    `gorm:"default:0" json:"size"`
	Count           uint32    `gorm:"default:0" json:"count"`
	CoreDescription string    `gorm:"type:text" json:"core_description,omitempty"` // 可为空
	PermissionNote  string    `gorm:"type:text" json:"permission_note,omitempty"`  // 可为空
	Remark          string    `gorm:"type:text" json:"remark,omitempty"`           // 可为空
	CreatedAt       time.Time `json:"created_at"`
}

func (Info) TableName() string {
	return "info"
}
