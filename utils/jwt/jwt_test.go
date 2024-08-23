package jwt

import (
	"os"
	"testing"
)

func TestJWT(t *testing.T) {

	secret := "your_secret_key"
	os.Setenv(SECRET_KEY, secret)
	in := map[string]interface{}{
		"email":    "123qwe@qq.com",
		"username": "admin",
	}

	tokenString, err := CreateToken(in)
	if err != nil {
		t.Errorf("Error creating token: %s", err.Error())
		return
	}

	_, err = VerifyToken(tokenString)

	if err != nil {
		t.Errorf("Error verifying token: %s", err.Error())
		return
	}

	_, err = VerifyToken("error.jwt.token")
	if err == nil {
		t.Errorf("Error verifying token: %s", err.Error())
		return
	}
}
