package http

import (
	"encoding/json"

	"net/http"

	"api.it-t.xyz/utils/errors"
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
