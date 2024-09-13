package jwt

import (
	"time"

	"errors"
	"os"

	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	Email     string  `json:"email"`
	Username  string  `json:"username"`
	Avatar    string  `json:"avatar"`
	ExpiresAt float64 `json:"exp"`
	Issuer    string  `json:"iss"`
}

const SECRET_KEY = "JWT_SECRET"

func getKey() ([]byte, error) {

	if os.Getenv(SECRET_KEY) == "" {
		return nil, errors.New("secret is not set")
	}
	return []byte(os.Getenv(SECRET_KEY)), nil
}

func CreateToken(in map[string]interface{}) (string, error) {
	secret, err := getKey()
	if err != nil {
		return "", err
	}
	claim := make(jwt.MapClaims, len(in)+2) // add issuer and expires
	claim["iss"] = "http://it-t.xyz"
	claim["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix() // not expired in one year
	for key, value := range in {
		claim[key] = value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (*MyClaims, error) {

	secret, err := getKey()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("failed to parse claims")
	}

	var myClaims MyClaims

	jsonString, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonString, &myClaims); err != nil {
		return nil, err
	}

	return &myClaims, nil
}
