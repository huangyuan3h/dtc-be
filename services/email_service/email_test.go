package email_service

import (
	"testing"
)

func TestCreateNew(t *testing.T) {

	err := SendEmailWithResend(SendMessageBody{
		Subject: "Subject",
		Content: "Content ",
		ToEmail: "huangyuan3h@gmail.com",
	})

	if err != nil {
		t.Error(err)
	}

}
