package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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
	AuditRule   string `gorm:"column:audit_rule;not null"`
	Description string `gorm:"column:description;not null"`
	Users       []User `gorm:"many2many:user_projects;"`
	Items       []Item `gorm:"foreignKey:ProjectId"`
	Apikey      string `gorm:"column:apikey"`
	//AccessKey   string `gorm:"column:access_key;not null;uniqueIndex,size:64"`
	//SecretKey   string `gorm:"column:secret_key;not null"`
	HookUrl string `gorm:"column:hook_url;not null"`
}

type UserProject struct {
	UserID    uint `gorm:"primaryKey"`
	ProjectID uint `gorm:"primaryKey"`
	Role      int  `gorm:"column:role"`
}

type ProjectPermit struct {
	ProjectID   uint `json:"project_id"`
	ProjectRole int  `json:"project_role"`
}

type UserResponse struct {
	Name        string `json:"name"`
	ID          uint   `json:"id"` //user_id
	Avatar      string `json:"avatar"`
	Email       string `json:"email"`
	ProjectRole int    `json:"project_role"`
	Role        int    `json:"role"`
}

type ProjectList struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
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
	Comments   []Comment       `gorm:"foreignKey:ItemId;references:ID;constraint:OnDelete:CASCADE"`
	Auditor    uint            `gorm:"column:auditor;"`
	Reason     string          `gorm:"column:reason"`
	Pictures   GormStringSlice `gorm:"type:json"`
	HookUrl    string          `gorm:"column:hook_url;not null"`
	HookId     uint            `gorm:"column:hook_id;not null;index"` //调用方该项目的id
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

type UserInfos struct {
	Name          string          `json:"name"`
	Avatar        string          `json:"avatar"`
	Email         string          `json:"email"`
	UserRole      int             `json:"user_role"`
	ProjectPermit []ProjectPermit `json:"project_permit"`
}

type RemoveItemStatus struct {
	Status int  `json:"status"`
	HookId uint `json:"hook_id"`
}

type RemoveItemsStatus struct {
	Items []RemoveItemStatus `json:"items"`
}

type GormStringSlice []string

func (g GormStringSlice) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *GormStringSlice) Scan(value interface{}) error {
	if value == nil {
		*g = []string{}
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan GormStringSlice: %v", value)
	}

	return json.Unmarshal(b, g)
}
