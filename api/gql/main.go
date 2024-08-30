package main

import (
	"encoding/json"

	"net/http"

	"log"

	"utils/errors"

	"api.it-t.xyz/gql/schema"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/graphql-go/graphql"
)

type QueryResponse struct {
	Query string `json:"query"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	schema, err := schema.CreateSchema()
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// 处理请求
	bodyStr := request.Body
	var response QueryResponse
	err = json.Unmarshal([]byte(bodyStr), &response)

	if err != nil {
		return errors.New("request is invalid", http.StatusInternalServerError).GatewayResponse()
	}

	params := graphql.Params{Schema: schema, RequestString: response.Query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
		return errors.New("invalid query", http.StatusInternalServerError).GatewayResponse()
	}

	rJSON, _ := json.Marshal(r)
	return events.APIGatewayProxyResponse{
		Body:       string(rJSON),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
