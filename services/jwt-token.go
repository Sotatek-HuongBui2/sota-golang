package services

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}

func CreateToken(userId string) (tokenStr string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "VTCCanteen",
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			NotBefore: time.Now().Unix(),
		},
	})

	tokenStr, err = token.SignedString([]byte(os.Getenv("HMAC_SECRET")))
	return
}

func VerifyToken(tokenStr string) (claims *TokenClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid token")
		}

		return []byte(os.Getenv("HMAC_SECRET")), nil
	})

	if err != nil {
		return
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("Invalid token")
	}
}
