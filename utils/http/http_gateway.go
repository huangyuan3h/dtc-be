package http

import (
	"encoding/json"
	"os"
	"time"

	"net/http"

	"utils/errors"

	"github.com/aws/aws-lambda-go/events"
)

func Ok(obj any, code int) (events.APIGatewayProxyResponse, error) {

	jsonData, err := json.Marshal(obj)
	if err != nil {
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonData),
		StatusCode: code,
	}, nil
}

func ResponseWithHeader(obj any, code int, header map[string]string) (events.APIGatewayProxyResponse, error) {

	jsonData, err := json.Marshal(obj)
	if err != nil {
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonData),
		StatusCode: code,
		Headers:    header,
	}, nil
}

func isProduction() bool {
	stage := os.Getenv("SST_STAGE")
	return stage == "production"
}

func Auth(token string) (events.APIGatewayProxyResponse, error) {
	domain := ""
	sameSite := http.SameSiteNoneMode

	if isProduction() {
		domain = ".it-t.xyz"
		sameSite = http.SameSiteLaxMode
	}

	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Path:     "/",
		Domain:   domain,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Secure:   true,
		HttpOnly: false,
		SameSite: sameSite,
	}

	return ResponseWithHeader(map[string]string{"Authorization": token}, http.StatusOK, map[string]string{
		"Set-Cookie":   cookie.String(),
		"Content-Type": "application/json",
	})
}
