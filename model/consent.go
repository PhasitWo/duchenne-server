package model

import (
	"gorm.io/plugin/soft_delete"
)

type Consent struct {
	ID        int    `json:"id"`
	Slug      string `json:"slug" gorm:"unique;not null"`
	CreateAt  int    `json:"createAt" gorm:"autoCreateTime;not null"`
	UpdateAt  int    `json:"updateAt" gorm:"autoUpdateTime;not null"`
	DeletedAt soft_delete.DeletedAt
	Body      string `json:"body" gorm:"not null"`
}

type UpsertConsentRequest struct {
	Slug string `json:"slug" binding:"required"`
	Body string `json:"body" binding:"required"`
}
