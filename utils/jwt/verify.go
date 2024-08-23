package jwt

import (
	"strings"

	"errors"

	"github.com/aws/aws-lambda-go/events"
)

func VerifyRequest(request events.APIGatewayV2HTTPRequest) (*MyClaims, error) {

	authStr := request.Headers["Authorization"]

	if authStr == "" {
		authStr = request.Headers["authorization"]
	}

	if authStr == "" {

		for _, cookie := range request.Cookies {

			parts := strings.Split(cookie, "=")

			if parts[0] == "Authorization" {
				authStr = parts[1]
				break
			}
		}
	}

	if authStr == "" {
		return nil, errors.New("authorization is not found in request")
	}

	myClaims, err := VerifyToken(authStr)

	if err != nil {
		return nil, err
	}

	return myClaims, nil
}
