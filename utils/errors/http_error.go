package errors

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type HttpError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func New(message string, code int) *HttpError {
	return &HttpError{
		Message: message,
		Code:    code,
	}
}

type ErrorResponseMethod interface {
	ToString() string
	GatewayResponse() (events.APIGatewayProxyResponse, error)
}

func (e HttpError) ToString() string {

	jsonData, err := json.Marshal(e)
	if err != nil {
		return ""
	}

	return string(jsonData)
}

func (e HttpError) GatewayResponse() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: e.Code,
		Body:       e.ToString(),
	}, nil
}

func StringToErrorMessage(body string) (HttpError, error) {

	var httpError HttpError
	err := json.Unmarshal([]byte(body), &httpError)
	return httpError, err
}
