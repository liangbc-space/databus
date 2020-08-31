package utils

import (
	"demo/system"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtClaims struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Mobile   string `json:"mobile"`
	jwt.StandardClaims
}

const JWTTokenExpireDuration = time.Hour * 24

func GenerateToken(id uint, username, mobile string) (token string, err error) {
	claims := JwtClaims{
		id,
		username,
		mobile,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(JWTTokenExpireDuration).Unix(),
			Issuer:    "梁彬城",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenClaims.SignedString([]byte(system.SystemConfig.JWTSecret))
	return token, err
}

func ParseToken(token string) (claims *JwtClaims, err error) {

	tokenClaims, err := jwt.ParseWithClaims(token, new(JwtClaims), func(token *jwt.Token) (interface{}, error) {
		return []byte(system.SystemConfig.JWTSecret), nil
	})

	if err != nil {
		return claims, err
	}

	claims, ok := tokenClaims.Claims.(*JwtClaims)
	if !ok || !tokenClaims.Valid {
		return claims, err
	}

	return claims, err
}
