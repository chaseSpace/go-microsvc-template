package auth

import "github.com/golang-jwt/jwt/v5"

func GenerateJwT(claims jwt.Claims, signKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(signKey))
	return ss, err
}
