package token_manager

import (
	"testing"
)

func TestCreateNew(t *testing.T) {
	email := "aaaab@qq.com"

	tokenClient := New()

	token, err := tokenClient.CreateToken(&email)

	if err != nil {
		t.Error(err)
	}
	print(token)
}

func TestConsume(t *testing.T) {
	token := "01J69GKKJ2EC2E5WE79Z3QMXFC"

	tokenClient := New()

	_, err := tokenClient.ConsumeToken(&token)

	if err != nil {
		t.Error(err)
	}

}
