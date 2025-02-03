package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string    `gorm:"column:name;unique;not null"` // 指定列名为 name，唯一约束且不能为空
	Email    string    `gorm:"column:email;not null"`
	Avatar   string    `gorm:"column:avatar"`
	UserRole int       `gorm:"column:user_role"`
	Projects []Project `gorm:"many2many:user_projects;"`
}
type Project struct {
	gorm.Model
	ProjectName string `gorm:"column:project_name;not null"`
	Users       []User `gorm:"many2many:user_projects;"`
}
type UserProject struct {
	UserID    uint `gorm:"primaryKey"`
	ProjectID uint `gorm:"primaryKey"`
	Role      int  `gorm:"column:role"`
}
type ProjectPermit struct {
	ProjectID   uint   `json:"project_id"`
	ProjectName string `json:"project_name"`
	ProjectRole int    `json:"project_role"`
}
type UserResponse struct {
	Name          string          `json:"name"`
	UserID        uint            `json:"user_id"`
	Avatar        string          `json:"avatar"`
	ProjectPermit []ProjectPermit `json:"project_permit"`
	Role          int             `json:"role"`
}
