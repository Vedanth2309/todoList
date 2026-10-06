package utils

import (
	"time"

	"tracker/config"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GenerateToken(id primitive.ObjectID) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": id.Hex(), "exp": time.Now().Add(7 * 24 * time.Hour).Unix()}).
		SignedString([]byte(config.C.JWTSecret))
}

func ParseToken(s string) (primitive.ObjectID, error) {
	tok, err := jwt.Parse(s, func(*jwt.Token) (any, error) { return []byte(config.C.JWTSecret), nil }, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return primitive.NilObjectID, err
	}
	sub, _ := tok.Claims.(jwt.MapClaims)["sub"].(string)
	return primitive.ObjectIDFromHex(sub)
}
