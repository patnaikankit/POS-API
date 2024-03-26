package middlewares

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func SplitToken(headerToken string) string {
	parseToken := strings.SplitAfter(headerToken, " ")
	tokenString := parseToken[1]
	return tokenString
}

func AuthenticateToken(tokenString string) error {
	// token check
	_, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return err
	}

	return nil
}
