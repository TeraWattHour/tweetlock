package twitter

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string
}

type UserResponse struct {
	Data UserData
}

type UserData struct {
	ID     string
	Name   string
	Handle string `json:"username"`
}

func GetTokensFromCode(code string) (*TokenResponse, error) {
	params := url.Values{}
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")
	params.Add("client_id", os.Getenv("TWITTER_CLIENT_ID"))
	params.Add("client_secret", os.Getenv("TWITTER_CLIENT_SECRET"))
	params.Add("redirect_uri", os.Getenv("TWITTER_REDIRECT"))
	params.Add("code_verifier", "challenge")

	credentials := fmt.Sprintf("%s:%s", os.Getenv("TWITTER_CLIENT_ID"), os.Getenv("TWITTER_CLIENT_SECRET"))
	creds := base64.StdEncoding.EncodeToString([]byte(credentials))

	req, err := http.NewRequest("POST", "https://api.twitter.com/2/oauth2/token?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", creds))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	by, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	tokenResponse := TokenResponse{}
	if err := json.Unmarshal(by, &tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func GetUserData(accessToken string) (*UserData, error) {
	req, err := http.NewRequest("GET", "https://api.twitter.com/2/users/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, errors.New("failed request")
	}

	by, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	userResponse := UserResponse{}
	if err := json.Unmarshal(by, &userResponse); err != nil {
		return nil, err
	}

	data := userResponse.Data

	if data.ID == "" || data.Name == "" {
		return nil, errors.New("incomplete user response")
	}

	return &data, nil
}
