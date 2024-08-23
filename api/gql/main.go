package main

import (
	"encoding/json"

	"net/http"

	"log"

	"api.it-t.xyz/utils/errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/graphql-go/graphql"
)


type QueryResponse struct {
	Query    string `json:"query"`
}


func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	bodyStr := request.Body

	var response QueryResponse
	err = json.Unmarshal([]byte(bodyStr), &response)

	if err != nil {
		return errors.New("Failed to create user", http.StatusInternalServerError).GatewayResponse()
	}

	params := graphql.Params{Schema: schema, RequestString: response.Query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	return events.APIGatewayProxyResponse{
		Body:      string(rJSON),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}