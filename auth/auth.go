package auth

import (
	"github.com/PhasitWo/duchenne-server/config"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type PatientClaims struct {
	PatientId   int `json:"patientId"`
	DeviceId int `json:"deviceId"`
	jwt.RegisteredClaims
}

func GeneratePatientToken(userId int) (string, error) {
	expirationTime := time.Now().Add(90 * 24 * time.Hour)
	claims := &PatientClaims{
		PatientId: userId,
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
