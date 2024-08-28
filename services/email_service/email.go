package email_service

import (
	"errors"
	"os"

	errs "utils/errors"

	"github.com/go-playground/validator/v10"
	"github.com/resend/resend-go/v2"
)

func SendEmailWithResend(data SendMessageBody) error {

	validate := validator.New(validator.WithRequiredStructEnabled())
	errStruct := validate.Struct(data)

	// validation

	if errStruct != nil {
		firstErr := errStruct.(validator.ValidationErrors)[0]
		var errMessage string
		switch t := firstErr.StructField(); t {
		case "FromEmail":
			errMessage = errs.NotValidEmail
		case "Subject":
			errMessage = errs.SubjectInvalid
		case "Content":
			errMessage = errs.ContentInvalid
		}

		return errors.New(errMessage)
	}

	apiKey := os.Getenv("EmailToken")

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "admin@it-t.xyz",
		To:      []string{data.ToEmail},
		Subject: data.Subject,
		Html:    data.Content,
	}

	_, err := client.Emails.Send(params)

	return err
}
