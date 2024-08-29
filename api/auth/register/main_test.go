package main

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestGetTitleContent(t *testing.T) {
	id := "idxxid"
	locale := "en"
	title, content, err := GetTitleContent(locale, id)

	if err != nil {
		t.Error(err)
	}
	t.Error(title)
	t.Log(content)
}

func TestHandler(t *testing.T) {
	input := events.APIGatewayV2HTTPRequest{
		Body: "{\"email\":\"huangyuan3h@gmail.com\",\"locale\":\"en\"}",
	}
	result, _ := Handler(input)

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", result.StatusCode)
	}
}
