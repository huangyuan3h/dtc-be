package main

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	input := events.APIGatewayV2HTTPRequest{
		Body: "{\"token\":\"01J6F5BQD9AC7TY7X5RX13ZWDQ\"}",
	}
	result, _ := Handler(input)

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", result.StatusCode)
	}
}
