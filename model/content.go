package model

import "gorm.io/plugin/soft_delete"

type Content struct {
	ID          int `json:"id"`
	CreateAt    int `json:"createAt" gorm:"autoCreateTime;not null"`
	UpdateAt    int `json:"updateAt" gorm:"autoUpdateTime;not null"`
	DeletedAt   soft_delete.DeletedAt
	Title       string `json:"title"`
	Body        string `json:"body"`
	IsPublished bool   `json:"isPublished" gorm:"not null"`
	Order       int    `json:"order" gorm:"not null;default:0"`
	DoctorID    int
	Doctor      Doctor `json:"doctor"`
}
