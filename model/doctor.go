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
	Specialist *string `json:"specialist"`
	Role       Role    `json:"role" gorm:"not null"`
}

type TrimDoctor struct {
	Doctor
	Username *string `json:"username" gorm:"-"`
	Password *string `json:"password" gorm:"-"`
}

type UpdateProfileRequest struct {
	FirstName  string  `json:"firstName" binding:"required"`
	MiddleName *string `json:"middleName"`
	LastName   string  `json:"lastName" binding:"required"`
	Username   string  `json:"username" binding:"required"`
	Password   string  `json:"password" binding:"required"`
	Specialist *string `json:"specialist"`
}

type CreateDoctorRequest struct {
	FirstName  string  `json:"firstName" binding:"required"`
	MiddleName *string `json:"middleName"`
	LastName   string  `json:"lastName" binding:"required"`
	Username   string  `json:"username" binding:"required,max=20"`
	Password   string  `json:"password" binding:"required"`
	Role       Role    `json:"role" binding:"required"`
	Specialist *string `json:"specialist"`
}
