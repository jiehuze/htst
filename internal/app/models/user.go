package models

import "time"

type SysUser struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username    string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password    string    `gorm:"size:100;not null" json:"password"` // bcrypt 加密后存储
	Roles       *string   `gorm:"size:50" json:"roles"`
	Status      int       `gorm:"size:50" json:"status"`
	Description *string   `gorm:"size:255" json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (SysUser) TableName() string {
	return "sys_user"
}
