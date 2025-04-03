package model

import "gorm.io/plugin/soft_delete"

type Question struct {
	ID        int     `json:"id"`
	Topic     string  `json:"topic" gorm:"not null"`
	Question  string  `json:"question" gorm:"not null"`
	CreateAt  int     `json:"createAt" gorm:"not null"`
	Answer    *string `json:"answer"`   // nullable
	AnswerAt  *int    `json:"answerAt"` // nullable
	PatientID int     `json:"-" gorm:"not null"`
	Patient   Patient `json:"patient"`
	DoctorID  *int    `json:"-"`
	Doctor    *Doctor `json:"doctor"` // nullable
	DeletedAt soft_delete.DeletedAt
}

type SafeQuestion struct {
	Question
	Doctor TrimDoctor `json:"doctor"`
}

type QuestionTopic struct {
	ID        int         `json:"id"`
	Topic     string      `json:"topic"`
	CreateAt  int         `json:"createAt"`
	AnswerAt  *int        `json:"answerAt"` // nullable
	PatientID int         `json:"-"`
	Patient   Patient     `json:"patient"`
	DoctorID  int         `json:"-"`
	Doctor    *TrimDoctor `json:"doctor"` // nullable
}
