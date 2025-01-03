package auth

import (
	"time"

	"github.com/PhasitWo/duchenne-server/config"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type PatientClaims struct {
	PatientId int `json:"patientId"`
	DeviceId  int `json:"deviceId"`
	jwt.RegisteredClaims
}

func GeneratePatientToken(patientId int, deviceId int) (string, error) {
	expirationTime := time.Now().Add(90 * 24 * time.Hour)
	claims := &PatientClaims{
		PatientId: patientId,
		DeviceId: deviceId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT_KEY))
}

type DoctorClaims struct {
	DoctorId int `json:"doctorId"`
	Role model.Role `json:"role"`
	jwt.RegisteredClaims
}

func GenerateDoctorToken(doctorId int, role model.Role) (string, error) {
	expirationTime := time.Now().Add(300*24 * time.Hour)
	claims := &DoctorClaims{
		DoctorId: doctorId,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT_KEY))
}

// unused
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// unused
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
