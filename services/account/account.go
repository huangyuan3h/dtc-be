package account

import (
	db "utils/dynamodb"

	"errors"
	errs "utils/errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const tableName = "auth"

type Auth struct {
	Account
	client *db.Client
}

type AuthMethod interface {
	CreateAccount(email, password *string) error
	VerifyLogin(email, password *string) error
}

func New() Auth {
	client := db.New(tableName)

	return Auth{client: &client}
}

func (a Auth) GetKey() map[string]types.AttributeValue {
	email, err := attributevalue.Marshal(a.Email)
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"email": email}
}

func (a Auth) CreateAccount(email, password *string) error {

	item, err := a.client.FindById("email", *email)

	if err != nil {
		return err
	}

	if item != nil {
		return errors.New(errs.AccountAlreadyExists)
	}

	encryptedPassword, err := encrypt([]byte(*password))

	if err != nil {
		return err
	}

	auth := Account{
		Email:    *email,
		Password: encryptedPassword,
		Status:   "actived",
	}

	return a.client.CreateOrUpdate(auth)
}

func (a Auth) VerifyLogin(email, password *string) error {
	item, err := a.client.FindById("email", *email)
	if err != nil {
		return err
	}

	if item == nil {
		return errors.New(errs.AccountNotExist)
	}

	account := Account{}

	err = attributevalue.UnmarshalMap(item, &account)

	if err != nil {
		return err
	}

	if account.Status != "actived" {
		return errors.New(errs.AccountIsNotActive)
	}

	err = attributevalue.UnmarshalMap(item, &a)
	if err != nil {
		return errors.New(errs.UnmarshalError)
	}

	decryptedPassword, err := decrypt(a.Password)

	if err != nil {
		return errors.New(errs.PasswordDecryptedError)
	}

	if string(decryptedPassword) != string(*password) {
		return errors.New(errs.PasswordIncorrect)
	}

	return nil
}
