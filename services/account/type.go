package account

type Account struct {
	Email    string `json:"email" dynamodbav:"email"`
	Password []byte `json:"password" dynamodbav:"password"`
	Status   string `json:"status" dynamodbav:"status"`
}
