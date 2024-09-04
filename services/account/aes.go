package account

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"os"
)

const SECRET_KEY = "AUTH_SECRET"

func getkey() ([]byte, error) {

	if os.Getenv(SECRET_KEY) == "" {
		return nil, errors.New("secret is not set")
	}
	return []byte(os.Getenv(SECRET_KEY)), nil
}

func encrypt(data []byte) ([]byte, error) {
	key, err := getkey()

	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func decrypt(data []byte) ([]byte, error) {
	key, err := getkey()
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
