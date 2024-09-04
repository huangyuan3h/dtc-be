package account

import (
	"os"
	"testing"
)

func Test_encrypt_password(t *testing.T) {
	os.Setenv("AUTH_SECRET", "GLbR3zUjXPbSKLwsSqNDTG3ODNkZYDdF")
	password := []byte("password")

	ciphertext, err := encrypt(password)
	if err != nil {
		t.Error("Error encrypting password:", err)
		return
	}

	decryptedtext, err := decrypt(ciphertext)
	if err != nil {
		t.Error("Error decrypting password:", err)
		return
	}

	// Verify the password
	if string(decryptedtext) != string(password) {
		t.Error("Password is incorrect")
		return
	}
}
