package main

import (
	"encoding/json"
	"net/http"
	"services/account"
	"services/profile"
	errs "utils/errors"
	awsHttp "utils/http"

	"utils/jwt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type LoginBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type LoginResponse struct {
	Authorization string `json:"Authorization"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var loginBody LoginBody
	err := json.Unmarshal([]byte(request.Body), &loginBody)
	if err != nil {
		return errs.New(errs.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	errStruct := validate.Struct(loginBody)

	if errStruct != nil {
		firstErr := errStruct.(validator.ValidationErrors)[0]
		var errMessage string
		switch t := firstErr.StructField(); t {
		case "Email":
			errMessage = errs.NotValidEmail
		case "Password":
			errMessage = errs.PasswordError
		}

		return errs.New(errMessage, http.StatusBadRequest).GatewayResponse()
	}

	// verify login
	authService := account.New()
	err = authService.VerifyLogin(&loginBody.Email, &loginBody.Password)
	if err != nil {
		return errs.New(errs.PasswordIncorrect, http.StatusBadRequest).GatewayResponse()
	}

	// search from profile
	profileService := profile.New()

	u, err := profileService.FindByEmail(&loginBody.Email)

	if err != nil {
		return errs.New(errs.UserProfileNotFound, http.StatusBadRequest).GatewayResponse()
	}

	// login to system using the email and password
	jwtObj := map[string]interface{}{
		"email":    loginBody.Email,
		"avatar":   u.Avatar,
		"userName": u.UserName,
	}
	jwt_token, err := jwt.CreateToken(jwtObj)

	if err != nil {
		return errs.New(err.Error(), http.StatusInternalServerError).GatewayResponse()
	}

	return awsHttp.Ok(LoginResponse{Authorization: jwt_token}, http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
