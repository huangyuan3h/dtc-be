package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"context"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"path/filepath"

	"os"

	"services/profile"
	"utils/errors"
	googleAuth "utils/googleAuth"
	awsHttp "utils/http"
	myJWT "utils/jwt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

type GoogleLoginRequest struct {
	Credential string `json:"credential"`
}

type LoginResponse struct {
	Authorization string `json:"Authorization"`
}

func downloadImage(ctx context.Context, url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	limitedReader := io.LimitReader(resp.Body, 5*1024*1024) // limit only 5mb allowed

	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	if _, err := io.Copy(io.Discard, limitedReader); err != nil {
		return nil, fmt.Errorf("image size exceeds 5MB limit")
	}
	return data, nil
}

func uploadToS3(ctx context.Context, data []byte, bucket string, key string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}
	client := s3.NewFromConfig(cfg)

	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to decode image config: %w", err)
	}
	if format != "jpeg" && format != "png" {
		return fmt.Errorf("unsupported image format: %s", format)
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(data),
		ContentType:   aws.String(fmt.Sprintf("image/%s", format)),
		ContentLength: aws.Int64(int64(len(data))),
		ACL:           types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return fmt.Errorf("failed to upload image to S3: %w", err)
	}

	return nil
}

func handleGoogleLogin(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	wg := sync.WaitGroup{}

	var loginReq GoogleLoginRequest
	err := json.Unmarshal([]byte(request.Body), &loginReq)
	if err != nil {
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	// verify credential
	userInfo, err := googleAuth.VerifyGoogleToken(ctx, loginReq.Credential)

	if err != nil {
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	userRepo := profile.New()
	u, err := userRepo.FindByEmail(&userInfo.Email)
	if err != nil {
		return errors.New(err.Error(), http.StatusInternalServerError).GatewayResponse()
	}

	if u.Email == "" { // not found
		bucketName := os.Getenv("AVATAR_BUCKET_NAME")
		fileName := fmt.Sprintf("%s-%d%s", uuid.New(), time.Now().Unix(), filepath.Ext(userInfo.Picture))

		data, err := downloadImage(ctx, userInfo.Picture)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to download image: %s", err.Error()), http.StatusInternalServerError).GatewayResponse()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err = uploadToS3(ctx, data, bucketName, fileName)
			if err != nil {
				log.Printf("Failed to upload image to S3: %v", err)
			}
		}()

		if err != nil {
			return errors.New(fmt.Sprintf("failed to upload image to S3: %s", err.Error()), http.StatusInternalServerError).GatewayResponse()
		}

		u = &profile.User{
			Email:    userInfo.Email,
			UserName: userInfo.Name,
			Avatar:   fmt.Sprintf("https://%s.s3.us-east-1.amazonaws.com/%s", bucketName, fileName),
		}

		if err := userRepo.CreateNew(u); err != nil {
			return errors.New("Failed to create user", http.StatusInternalServerError).GatewayResponse()
		}
	}

	jwtObj := map[string]interface{}{
		"email":    u.Email,
		"avatar":   u.Avatar,
		"userName": u.UserName,
	}
	jwtToken, err := myJWT.CreateToken(jwtObj)
	if err != nil {
		return errors.New(err.Error(), http.StatusInternalServerError).GatewayResponse()
	}
	wg.Wait()

	return awsHttp.Auth(jwtToken)
}

func main() {
	lambda.Start(handleGoogleLogin)
}
