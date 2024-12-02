package auth

import (
	"github.com/PhasitWo/duchenne-server/config"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Claims struct {
	User_id int `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(user_id int) (string, error) {
	expirationTime := time.Now().Add(90 * 24 * time.Hour)
	claims := &Claims{
		User_id: user_id,
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
