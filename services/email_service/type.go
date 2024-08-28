package email_service

type SendMessageBody struct {
	Subject string `json:"subject" validate:"required,min=6,max=50"`
	Content string `json:"content" validate:"required,min=6,max=500"`
	ToEmail string `json:"toEmail" validate:"required,email"`
}
