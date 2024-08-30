package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"services/token_manager"
	"time"

	errs "utils/errors"

	awsHttp "utils/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type VerifyTokenBody struct {
	Token string `json:"token" validate:"required,len=26"`
}

type VerifyTokenResponse struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var verifyTokenBody VerifyTokenBody
	err := json.Unmarshal([]byte(request.Body), &verifyTokenBody)
	if err != nil {
		return errs.New(errs.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	errStruct := validate.Struct(verifyTokenBody)

	if errStruct != nil {
		firstErr := errStruct.(validator.ValidationErrors)[0]
		var errMessage string
		switch t := firstErr.StructField(); t {
		case "Email":
			errMessage = errs.NotValidEmail
		}

		return errs.New(errMessage, http.StatusBadRequest).GatewayResponse()
	}

	tokenService := token_manager.New()
	token, err := tokenService.SearchToken(&verifyTokenBody.Token)
	if err != nil {
		return errs.New(err.Error(), http.StatusBadRequest).GatewayResponse()
	}

	if token.IsConsumed == "true" {
		return awsHttp.Ok(VerifyTokenResponse{Email: "", Message: errs.TokenConsumed}, http.StatusBadRequest)
	}

	if token.ExpireAt < fmt.Sprint(time.Now().Unix()) {
		return awsHttp.Ok(VerifyTokenResponse{Email: "", Message: errs.TokenHasExpired}, http.StatusBadRequest)
	}

	return awsHttp.Ok(VerifyTokenResponse{Email: token.ConsumedBy, Message: ""}, http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
