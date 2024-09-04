package profile

type User struct {
	Email    string `json:"email" dynamodbav:"email"`
	Avatar   string `json:"avatar" dynamodbav:"avatar"`
	UserName string `json:"userName" dynamodbav:"userName"`
	Bio      string `json:"bio" dynamodbav:"bio"`
}
