package model

import "gorm.io/datatypes"

type ActivityLog struct {
	ID         int            `json:"id"`
	Claims     datatypes.JSON `json:"claims" gorm:"not null"`
	Method     string         `json:"method" gorm:"not null"`
	RequestURL string         `json:"requestURL" gorm:"not null"`
	Status     int            `json:"status" gorm:"not null"`
	Data       datatypes.JSON `json:"data"`
	CreateAt   int            `json:"createAt" gorm:"autoCreateTime;not null"`
}
