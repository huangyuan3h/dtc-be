package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"services/token_manager"
	errs "utils/errors"
	awsHttp "utils/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type CreateAccountBody struct {
	Token    string `json:"token" validate:"required,len=26"`
	Name     string `json:"name" validate:"required,min=6,max=50"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type CreateAccountResponse struct {
	Message string `json:"message"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var createAccountBody CreateAccountBody
	err := json.Unmarshal([]byte(request.Body), &createAccountBody)
	if err != nil {
		return errs.New(errs.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	errStruct := validate.Struct(createAccountBody)

	if errStruct != nil {
		firstErr := errStruct.(validator.ValidationErrors)[0]
		var errMessage string
		switch t := firstErr.StructField(); t {
		case "Token":
			errMessage = errs.TokenIdInvalid
		case "Name":
			errMessage = errs.UseNameInvalid
		case "Password":
			errMessage = errs.PasswordError
		}

		return errs.New(errMessage, http.StatusBadRequest).GatewayResponse()
	}

	// detail validation
	var regContainsLow = regexp.MustCompile("[a-z]+")
	var regContainsUpper = regexp.MustCompile("[A-Z]+")
	var regContainsNumber = regexp.MustCompile("[0-9]+")

	if !regContainsLow.MatchString(createAccountBody.Password) || !regContainsUpper.MatchString(createAccountBody.Password) || !regContainsNumber.MatchString(createAccountBody.Password) {
		return errs.New(errs.PasswordError, http.StatusBadRequest).GatewayResponse()
	}

	// consume token
	tokenService := token_manager.New()
	err = tokenService.ConsumeToken(&createAccountBody.Token)
	if err != nil {
		return errs.New(err.Error(), http.StatusBadRequest).GatewayResponse()
	}
	// create account

	// create profile

	return awsHttp.Ok(CreateAccountResponse{Message: ""}, http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
