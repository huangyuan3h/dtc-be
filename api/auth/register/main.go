package main

import (
	"encoding/json"
	"net/http"
	"services/email_service"
	"services/token_manager"
	errs "utils/errors"
	awsHttp "utils/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type RegisterBody struct {
	Email string `json:"email" validate:"required,email"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var registerBody RegisterBody
	err := json.Unmarshal([]byte(request.Body), &registerBody)

	if err != nil {
		return errs.New(errs.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	errStruct := validate.Struct(registerBody)

	if errStruct != nil {
		firstErr := errStruct.(validator.ValidationErrors)[0]
		var errMessage string
		switch t := firstErr.StructField(); t {
		case "Email":
			errMessage = errs.NotValidEmail
		}

		return errs.New(errMessage, http.StatusBadRequest).GatewayResponse()
	}

	// create a token
	tokenService := token_manager.New()

	token, err := tokenService.CreateToken(&registerBody.Email)
	if err != nil {
		return errs.New(err.Error(), http.StatusBadRequest).GatewayResponse()
	}
	// send email

	err = email_service.SendEmailWithResend(email_service.SendMessageBody{
		Subject:   "register to it and tea",
		Content:   "please click the link below: " + token,
		ToEmail:   registerBody.Email,
		FromEmail: "admin@it-t.xyz",
	})

	return awsHttp.Ok(RegisterResponse{Token: token}, http.StatusCreated)
}

func main() {
	lambda.Start(Handler)
}
