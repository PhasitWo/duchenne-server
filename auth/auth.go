package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key")

type User struct {
    ID       uint64 `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
}

type Claims struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func GenerateToken(user User) (string, error) {
    expirationTime := time.Now().Add(1 * time.Hour)
    claims := &Claims{
        Username: user.Username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func VerifyPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}
