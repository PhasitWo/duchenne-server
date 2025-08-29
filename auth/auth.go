package auth

import (
	"errors"
	"time"

	"github.com/PhasitWo/duchenne-server/config"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type PatientRefreshClaims struct {
	PatientId int `json:"patientId"`
	jwt.RegisteredClaims
}

func GeneratePatientRefreshToken(patientId int) (string, error) {
	expirationTime := time.Now().Add(30 * 24 * time.Hour)
	claims := &PatientRefreshClaims{
		PatientId: patientId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT_REFRESH_KEY))
}

func ParsePatientRefreshToken(tokenString string) (patientId int, err error) {
	claims := &PatientRefreshClaims{PatientId: -1}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT_REFRESH_KEY), nil
	})
	if err != nil {
		return -1, err
	}
	if !token.Valid {
		return -1, errors.New("invalid token")
	}
	return claims.PatientId, nil
}

type PatientAccessClaims struct {
	PatientId int `json:"patientId"`
	DeviceId  int `json:"deviceId"`
	jwt.RegisteredClaims
}

func GeneratePatientAccessToken(patientId int, deviceId int) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &PatientAccessClaims{
		PatientId: patientId,
		DeviceId:  deviceId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT_KEY))
}

type DoctorClaims struct {
	DoctorId int        `json:"doctorId"`
	Role     model.Role `json:"role"`
	jwt.RegisteredClaims
}

func GenerateDoctorAccessToken(doctorId int, role model.Role) (string, error) {
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &DoctorClaims{
		DoctorId: doctorId,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT_KEY))
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
