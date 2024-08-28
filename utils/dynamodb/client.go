package dynamodb

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const STAGE_KEY = "SST_STAGE"
const PROJECT_STR = "-dtc-common-be-"

type Client struct {
	TableName *string
	Client    *dynamodb.Client
}

type ClientBaseMethod interface {
	CreateOrUpdate(in interface{}) error
	FindById(keyName, id string) (map[string]types.AttributeValue, error)
	DeleteById(keyName, id string) error
	ExecuteStatement(statement *string, parameters []types.AttributeValue)
}

func (c Client) CreateOrUpdate(in interface{}) error {
	item, err := attributevalue.MarshalMap(in)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: c.TableName,
	}

	_, err = c.Client.PutItem(context.TODO(), input)

	if err != nil {
		return err
	}

	return nil
}

func (c Client) FindById(keyName, id string) (map[string]types.AttributeValue, error) {
	if keyName == "" {
		return nil, errors.New("keyname is required")
	}
	val, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, err
	}

	key := map[string]types.AttributeValue{
		keyName: val,
	}

	return GetItem(c.Client, *c.TableName, key)
}

func (c Client) DeleteById(keyName, id string) error {
	if keyName == "" {
		return errors.New("keyname is required")
	}

	val, err := attributevalue.Marshal(id)
	if err != nil {
		return err
	}

	key := map[string]types.AttributeValue{
		keyName: val,
	}

	return DeleteItem(c.Client, *c.TableName, key)
}

func GetItem(svc *dynamodb.Client, tableName string, key map[string]types.AttributeValue) (map[string]types.AttributeValue, error) {

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}

	result, err := svc.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return result.Item, nil
}

func DeleteItem(svc *dynamodb.Client, tableName string, key map[string]types.AttributeValue) error {

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}

	_, err := svc.DeleteItem(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}

func New(tableName string) Client {
	client := initDynamo()
	t := getTableName(tableName)
	return Client{TableName: &t, Client: client}
}

func initDynamo() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	return client
}

func getTableName(tableName string) string {
	var stage string
	if os.Getenv(STAGE_KEY) != "" {
		stage = os.Getenv(STAGE_KEY)
	} else {
		stage = "dev"
	}

	return stage + PROJECT_STR + tableName
}

func (c Client) ExecuteStatement(executeInput *dynamodb.ExecuteStatementInput) (*dynamodb.ExecuteStatementOutput, error) {
	return c.Client.ExecuteStatement(context.TODO(), executeInput)
}
