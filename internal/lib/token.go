package lib

import (
	"log"
	"time"

	"github.com/bangueco/auction-api/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Iss string `json:"iss,omitempty"`
	Sub int64  `json:"sub,omitempty"`
	Aud string `json:"aud,omitempty"`
	Exp int64  `json:"exp,omitempty"`
	Nbf int64  `json:"nbf,omitempty"`
	Iat int64  `json:"iat,omitempty"`
}

func GenerateToken(userId int64) (string, error) {
	cfg := config.Load()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "github.com/bangueco/auction-api",
		"sub": userId,
		"aud": "user",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now(),
	})

	tokenString, err := token.SignedString([]byte(cfg.TOKEN_SECRET))

	if err != nil {
		log.Printf("error signing token: %s", err)
		return "", err
	}

	return tokenString, nil
}
