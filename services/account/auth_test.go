package account

import (
	"os"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	os.Setenv("AUTH_SECRET", "GLbR3zUjXPbSKLwsSqNDTG3ODNkZYDdF")
	auth := New()

	email := "email1@example.com"
	password := "password"

	err := auth.CreateAccount(&email, &password)

	if err != nil {
		t.Error(err)
	}
}

func TestVerifyLogin(t *testing.T) {
	os.Setenv("AUTH_SECRET", "GLbR3zUjXPbSKLwsSqNDTG3ODNkZYDdF")
	auth := New()

	email := "email1@example.com"
	password := "password"

	err := auth.VerifyLogin(&email, &password)

	if err != nil {
		t.Error(err)
	}
}
