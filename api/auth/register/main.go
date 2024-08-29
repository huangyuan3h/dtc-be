package main

import (
	"encoding/json"
	"net/http"
	"services/email_service"
	"services/token_manager"
	errs "utils/errors"
	awsHttp "utils/http"

	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type RegisterBody struct {
	Email  string `json:"email" validate:"required,email"`
	Locale string `json:"locale" validate:"required"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type Config struct {
	En map[string]string `json:"en"` // English translations
	Fr map[string]string `json:"fr"` // French translations
	Zh map[string]string `json:"zh"` // Chinese translations
}

var EmailConfig = Config{
	En: map[string]string{
		"subject": "Welcome to IT&TEA",
		"content": "<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\" /><meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\" /><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" /><title>Register for IT&TEA</title><style>body {font-family: Arial, sans-serif;  background-color: #f4f4f4;  color: #333;  padding: 20px;  text-align: center;}.container {max-width: 600px;  margin: 0 auto;  background-color: #fff;  padding: 20px;  border-radius: 8px;  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);}h1 {color: #298717;}p {font-size: 16px;  line-height: 1.5;}.button {display: inline-block;  background-color: #298717;  color: white;  text-decoration: none;  padding: 10px 20px;  border-radius: 5px;  font-size: 18px;  margin-top: 20px;}.button:hover {background-color: #22ad07;}</style></head><body><div class=\"container\"><h1>Welcome to IT&amp;TEA</h1><p>We're excited to have you join us! Please click the button below to register.</p><a href=\"https://www.it-t.xyz/create-account/{id}\" class=\"button\" style=\"color:white !important\">Register Now</a><p>If you did not request this email, please ignore it.</p></div></body></html>",
	},
	Fr: map[string]string{
		"subject": "Bienvenue chez IT&TEA",
		"content": "<!DOCTYPE html><html lang=\"fr\"><head><meta charset=\"UTF-8\" /><meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\" /><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" /><title>Register for IT&TEA</title><style>body {font-family: Arial, sans-serif;  background-color: #f4f4f4;  color: #333;  padding: 20px;  text-align: center;}.container {max-width: 600px;  margin: 0 auto;  background-color: #fff;  padding: 20px;  border-radius: 8px;  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);}h1 {color: #298717;}p {font-size: 16px;  line-height: 1.5;}.button {display: inline-block;  background-color: #298717;  color: white;  text-decoration: none;  padding: 10px 20px;  border-radius: 5px;  font-size: 18px;  margin-top: 20px;}.button:hover {background-color: #22ad07;}</style></head><body><div class=\"container\"><h1>Bienvenue chez IT&TEA</h1><p>Nous sommes ravis de vous accueillir parmi nous ! Veuillez cliquer sur le bouton ci-dessous pour vous inscrire.</p><a href=\"https://www.it-t.xyz/create-account/{id}\" class=\"button\" style=\"color:white !important\">Inscrivez-vous maintenant</a><p>Si vous n’avez pas demandé cet e-mail, veuillez l’ignorer.</p></div></body></html>",
	},
	Zh: map[string]string{
		"subject": "欢迎加入IT&TEA",
		"content": "<!DOCTYPE html><html lang=\"zh\"><head><meta charset=\"UTF-8\" /><meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\" /><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" /><title>Register for IT&TEA</title><style>body {font-family: Arial, sans-serif;  background-color: #f4f4f4;  color: #333;  padding: 20px;  text-align: center;}.container {max-width: 600px;  margin: 0 auto;  background-color: #fff;  padding: 20px;  border-radius: 8px;  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);}h1 {color: #298717;}p {font-size: 16px;  line-height: 1.5;}.button {display: inline-block;  background-color: #298717;  color: white;  text-decoration: none;  padding: 10px 20px;  border-radius: 5px;  font-size: 18px;  margin-top: 20px;}.button:hover {background-color: #22ad07;}</style></head><body><div class=\"container\"><h1>欢迎加入 IT&TEA</h1><p>我们很高兴你能加入我们！请点击下面的按钮进行注册。</p><a href=\"https://www.it-t.xyz/create-account/{id}\" class=\"button\" style=\"color:white !important\">立即注册</a><p>如果您没有请求此邮件，请忽略它。</p></div></body></html>",
	},
}

func GetTitleContent(locale string, id string) (string, string, error) {

	title := ""
	content := ""

	switch locale {
	case "fr":
		title = EmailConfig.Fr["subject"]
		content = EmailConfig.Fr["content"]
	case "zh":
		title = EmailConfig.Zh["subject"]
		content = EmailConfig.Zh["content"]
	default:
		title = EmailConfig.En["subject"]
		content = EmailConfig.En["content"]
	}
	content = strings.ReplaceAll(content, "{id}", id)
	return title, content, nil
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
	title, content, err := GetTitleContent(registerBody.Locale, token)
	if err != nil {
		return errs.New(err.Error(), http.StatusBadRequest).GatewayResponse()
	}

	err = email_service.SendEmailWithResend(email_service.SendMessageBody{
		Subject: title,
		Content: content,
		ToEmail: registerBody.Email,
	})

	if err != nil {
		return errs.New(err.Error(), http.StatusBadRequest).GatewayResponse()
	}

	return awsHttp.Ok(RegisterResponse{Token: token}, http.StatusCreated)
}

func main() {
	lambda.Start(Handler)
}
