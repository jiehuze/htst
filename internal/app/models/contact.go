package models

import "time"

// Contact 联系人信息实体
type Contact struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;size:100;not null" json:"name"`
	Phone     string    `gorm:"column:phone;size:20" json:"phone,omitempty"`
	Email     string    `gorm:"column:email;size:150" json:"email,omitempty"`
	Address   string    `gorm:"column:address" json:"address,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Contact) TableName() string {
	return "contacts"
}
