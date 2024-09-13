package main

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	input := events.APIGatewayV2HTTPRequest{
		Body: "{\"email\":\"huangyuan3h@gmail.com\",\"password\":\"P@$$uu0rd123\"}",
	}
	result, _ := Handler(input)

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", result.StatusCode)
	}
}
