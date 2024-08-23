package jwt

import (
	"testing"

	"os"

	"github.com/aws/aws-lambda-go/events"
)

const mockToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdmF0YXIiOiIiLCJlbWFpbCI6Imh1YW5neXVhbjNoQGdtYWlsLmNvbSIsImV4cCI6MTcxMzg1MjUxNiwiaXNzIjoiaHR0cDovL25vcnRoLXBhdGguc2l0ZSIsInVzZXJOYW1lIjoiaHVhbmd5dWFuM2gifQ.Aam1vD60vhHtf2mBzRXaBAybjfkuKrS8ZFy3YE9qMo8"

func TestVerifyHeaderAuth(t *testing.T) {

	os.Setenv("JWT_SECRET", "h3OOumyH3vLgUhHve7bLPv8hgNXbxUQr")
	headers := &map[string]string{
		"authorization": mockToken,
	}

	input := events.APIGatewayV2HTTPRequest{
		Body:    "{\"subject\":\"oo2bq0bhal\",\"content\":\"f\",\"categories\":[\"k\",\"5\"],\"images\":[\"o\",\"m\"]}",
		Headers: *headers,
	}

	myclaim, err := VerifyRequest(input)

	if err != nil {
		t.Error("error verifying auth")
	}

	if myclaim.Username != "huangyuan3h" {
		t.Error("error verifying auth")
	}
}

func TestVerifyCookieAuth(t *testing.T) {

	os.Setenv("JWT_SECRET", "h3OOumyH3vLgUhHve7bLPv8hgNXbxUQr")
	cookies := &[]string{
		"Authorization=" + mockToken,
	}

	input := events.APIGatewayV2HTTPRequest{
		Body:    "{\"subject\":\"oo2bq0bhal\",\"content\":\"f\",\"categories\":[\"k\",\"5\"],\"images\":[\"o\",\"m\"]}",
		Cookies: *cookies,
	}

	myclaim, err := VerifyRequest(input)

	if err != nil {
		t.Error("error verifying auth")
	}

	if myclaim.Username != "huangyuan3h" {
		t.Error("error verifying auth")
	}
}
