package lib

import (
	"errors"
	"log"
	"time"

	"github.com/bangueco/auction-api/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnknownClaims = errors.New("unknown claims type")
	ErrTokenExpired  = errors.New("token expired")
)

type JWTClaims struct {
	Iss string `json:"iss,omitempty"`
	Sub int64  `json:"sub,omitempty"`
	Aud string `json:"aud,omitempty"`
	Exp int64  `json:"exp,omitempty"`
	Nbf int64  `json:"nbf,omitempty"`
	Iat int64  `json:"iat,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(userId int64) (string, error) {
	cfg := config.Load()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "github.com/bangueco/auction-api",
		"sub": userId,
		"aud": "user",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.TOKEN_SECRET))

	if err != nil {
		log.Printf("error signing token: %s", err)
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(token string) (bool, int64, error) {
	cfg := config.Load()

	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(cfg.TOKEN_SECRET), nil
	})

	if err != nil {
		log.Printf("Cannot parse token: %s", err)
		return false, 0, err
	}

	if !parsedToken.Valid {
		return false, 0, err
	}

	claims, ok := parsedToken.Claims.(*JWTClaims)

	if !ok {
		return false, 0, ErrUnknownClaims
	}

	if time.Now().Unix() > claims.Exp {
		log.Println("Token expired: ", claims)
		return false, 0, ErrTokenExpired
	}

	log.Println("Token verified: ", claims)
	return true, claims.Sub, nil
}
