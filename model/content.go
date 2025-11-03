package model

import "gorm.io/plugin/soft_delete"

type ContentType string

const (
	ARTICLE ContentType = "article"
	VIDEO   ContentType = "video"
)

type Content struct {
	ID            int `json:"id"`
	CreateAt      int `json:"createAt" gorm:"autoCreateTime;not null"`
	UpdateAt      int `json:"updateAt" gorm:"autoUpdateTime;not null"`
	DeletedAt     soft_delete.DeletedAt
	Title         string      `json:"title" gorm:"not null"`
	Body          string      `json:"body" gorm:"not null"`
	IsPublished   bool        `json:"isPublished" gorm:"not null"`
	Order         int         `json:"order" gorm:"not null;default:1"`
	ContentType   ContentType `json:"contentType" gorm:"not null;default:'article'"`
	CoverImageURL *string     `json:"coverImageURL"`
}

type CreateContentRequest struct {
	Title         string  `json:"title" binding:"required"`
	Body          string  `json:"body" binding:"required"`
	IsPublished   bool    `json:"isPublished"`
	Order         int     `json:"order" binding:"required"`
	CoverImageURL *string `json:"coverImageURL"`
	ContentType   string  `json:"contentType" binding:"oneof=article video"`
}
