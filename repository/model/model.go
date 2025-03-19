package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name     string    `gorm:"column:name;unique;not null"` // 指定列名为 name，唯一约束且不能为空
	Email    string    `gorm:"column:email;not null"`
	Avatar   string    `gorm:"column:avatar"`
	UserRole int       `gorm:"column:user_role"`
	Projects []Project `gorm:"many2many:user_projects;"`
	History  []History `gorm:"foreignKey:UserID"`
}
type Project struct {
	gorm.Model
	ProjectName string `gorm:"column:project_name;not null"`
	Logo        string `gorm:"column:logo;not null"`
	AudioRule   string `gorm:"column:audio_rule;not null"`
	Users       []User `gorm:"many2many:user_projects;"`
	Items       []Item `gorm:"foreignKey:ProjectId"`
	Apikey      string `gorm:"column:apikey;not null"`
}
type UserProject struct {
	UserID    uint `gorm:"primaryKey"`
	ProjectID uint `gorm:"primaryKey"`
	Role      int  `gorm:"column:role"`
}
type ProjectPermit struct {
	ProjectID uint `json:"project_id"`
	//ProjectName string `json:"project_name"`
	ProjectRole int `json:"project_role"`
}
type UserResponse struct {
	Name          string          `json:"name"`
	UserID        uint            `json:"user_id"`
	Avatar        string          `json:"avatar"`
	ProjectPermit []ProjectPermit `json:"project_permit"`
	Role          int             `json:"role"`
}
type ProjectList struct {
	ProjectId   uint   `json:"project_id"`
	ProjectName string `json:"project_name"`
}
type Item struct {
	gorm.Model
	Status     int             `gorm:"column:status;not null"`
	ProjectId  uint            `gorm:"column:project_id;not null"`
	Author     string          `gorm:"column:author;not null"`
	Tags       GormStringSlice `gorm:"type:json"`
	PublicTime time.Time       `gorm:"column:public_time;not null"`
	Content    string          `gorm:"column:content;not null"`
	Title      string          `gorm:"column:title;not null"`
	Comments   []Comment       `gorm:"foreignKey:ItemId"`
	Auditor    uint            `gorm:"column:auditor;not null"`
	Reason     string          `gorm:"column:reason"`
	Pictures   GormStringSlice `gorm:"type:json"`
	HookUrl    string          `gorm:"column:hook_url;not null"`
	HookId     int             `gorm:"column:hook_id;not null"`
}

type Comment struct {
	gorm.Model
	Content  string          `gorm:"column:content;not null"`
	Pictures GormStringSlice `gorm:"type:json"`
	ItemId   uint            `gorm:"not null;index"`
}
type History struct {
	gorm.Model
	UserID uint `gorm:"index"`
	ItemId uint `gorm:"index"`
}
type GormStringSlice []string

func (g GormStringSlice) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *GormStringSlice) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), g)
}
