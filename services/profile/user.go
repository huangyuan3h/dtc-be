package profile

import (
	db "utils/dynamodb"

	"errors"
	errs "utils/errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

const tableName = "user"

type Profile struct {
	User
	client *db.Client
}

type ProfileMethod interface {
	CreateNew(user *User) error
	FindByEmail(email *string) (*User, error)
}

func New() ProfileMethod {
	client := db.New(tableName)

	return Profile{client: &client}
}

func (u Profile) CreateNew(user *User) error {

	return u.client.CreateOrUpdate(user)
}

func (u Profile) FindByEmail(email *string) (*User, error) {

	item, err := u.client.FindById("email", *email)

	if err != nil {
		return nil, err
	}

	user := User{}

	err = attributevalue.UnmarshalMap(item, &user)
	if err != nil {
		return nil, errors.New(errs.UnmarshalError)
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
