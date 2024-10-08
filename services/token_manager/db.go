package token_manager

import (
	"crypto/rand"
	"fmt"
	"time"

	db "utils/dynamodb"

	"errors"
	errs "utils/errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/oklog/ulid/v2"
)

type TokenDB struct {
	client *db.Client
}

type TokenMethod interface {
	CreateToken(email *string, expireAt ...time.Time) (string, error)
	ConsumeToken(tokenId *string) (Token, error)
	SearchToken(tokenId *string) (Token, error)
}

const tableName = "token"

func New() TokenMethod {
	client := db.New(tableName)

	return TokenDB{client: &client}
}

// CreateToken 创建一个新的 token 并存储到 DynamoDB
func (t TokenDB) CreateToken(email *string, expireAt ...time.Time) (string, error) {
	tokenId := generateULID()

	// 设置默认过期时间 (1天)
	expiration := time.Now().Add(24 * time.Hour)
	if len(expireAt) > 0 {
		// 如果提供了 expireAt 参数，则使用该参数
		expiration = expireAt[0]
	}

	expireAtStr := fmt.Sprintf("%d", expiration.Unix()) // 过期时间为 Unix 时间戳格式（秒）

	token := &Token{
		TokenId:    tokenId,
		ExpireAt:   expireAtStr,
		ConsumedBy: *email,
		IsConsumed: "false",
	}

	return tokenId, t.client.CreateOrUpdate(token)
}

func generateULID() string {
	// 使用当前时间作为熵的起点
	t := time.Now().UTC()
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}

func (t TokenDB) ConsumeToken(tokenId *string) (Token, error) {

	token, err := t.SearchToken(tokenId)
	if err != nil {
		return token, err
	}

	if token.IsConsumed != "false" {
		return token, errors.New(errs.TokenConsumed)
	}

	token.IsConsumed = "true"

	err = t.client.CreateOrUpdate(token)
	if err != nil {
		return token, errors.New(errs.DBProcessError)
	}

	return token, nil
}

func (t TokenDB) SearchToken(tokenId *string) (Token, error) {
	item, err := t.client.FindById("tokenId", *tokenId)
	var token = Token{}
	if err != nil {
		return token, errors.New(errs.TokenNotFound + *tokenId)
	}

	err = attributevalue.UnmarshalMap(item, &token)
	if err != nil {
		return token, errors.New(errs.UnmarshalError)
	}

	return token, nil
}
