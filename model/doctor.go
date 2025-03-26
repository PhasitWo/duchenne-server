package model

// Doctor roles
type Role string

const (
	ROOT  Role = "root"
	ADMIN Role = "admin"
	USER  Role = "user"
)

type Doctor struct {
	ID         int     `json:"id"`
	FirstName  string  `json:"firstName" gorm:"not null"`
	MiddleName *string `json:"middleName"` // nullable
	LastName   string  `json:"lastName" gorm:"not null"`
	Username   string  `json:"username" gorm:"unique;not null"`
	Password   string  `json:"password" gorm:"not null"`
	Role       Role    `json:"role" gorm:"not null"`
}

type TrimDoctor struct {
	Doctor
	Username *string `json:"username" gorm:"-"`
	Password *string `json:"password" gorm:"-"`
}
