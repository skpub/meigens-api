package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(tokenString string) (string, error) {
	secret := os.Getenv("SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithJSONNumber())
	if err != nil {
		return "", fmt.Errorf("unauthorized. (invalid token)")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, _ := claims["exp"].(json.Number).Int64()
		if exp < time.Now().Unix() {
			return "", fmt.Errorf("unauthorized. (your token is expired)")
		} else {
			// if claims["exp"] > time.Now().Add(time.Hour * 24 * 3).Unix() {
			// }
			// Authorized
			return claims["user_id"].(string), nil
		}
	}
	return "", fmt.Errorf("unauthorized. (invalid token)")
}
