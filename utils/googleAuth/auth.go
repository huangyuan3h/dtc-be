package googleAuth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type GoogleTokenInfo struct {
	Aud              string `json:"aud"`
	Exp              string `json:"exp"`
	Iss              string `json:"iss"`
	Sub              string `json:"sub"`
	Email            string `json:"email"`
	Azp              string `json:"azp"`
	Name             string `json:"name"`
	Picture          string `json:"picture"`
	EmailVerified    string `json:"email_verified"`
	GivenName        string `json:"given_name"`
	FamilyName       string `json:"family_name"`
	Locale           string `json:"locale"`
	Hd               string `json:"hd"`
	ErrorDescription string `json:"error_description"`
}

func VerifyGoogleToken(ctx context.Context, token string) (*GoogleTokenInfo, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%s", token)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google api returned non-200 status code: %d", resp.StatusCode)
	}

	var tokenInfo GoogleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	exp, err := strconv.ParseInt(tokenInfo.Exp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("token expired field error")
	}

	if time.Now().Unix() > exp {
		return nil, fmt.Errorf("token is expired")
	}

	if tokenInfo.Aud != os.Getenv("GOOGLE_CLIENT_ID") {
		return nil, fmt.Errorf("token is not issued for this app")
	}

	return &tokenInfo, nil
}
