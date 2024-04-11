package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

)

type GoogleOauth struct {
	clientSecret string
	clientID     string
}

func newGoogleOauth() *GoogleOauth {
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")


	return &GoogleOauth{
		clientID:     googleClientId,
		clientSecret: googleClientSecret,
	}
}

func (g *GoogleOauth) Authorize(redirectUri string, scopes ...string) string {
	url := fmt.Sprintf(
		"%s?scope=%s&access_type=%s&include_granted_scopes=%s&response_type=%s&redirect_uri=%s&client_id=%s",
		googleAuthURL,
		url.QueryEscape("https://www.googleapis.com/auth/userinfo.profile"), // scope
		url.QueryEscape("online"),    // access_type
		url.QueryEscape("true"),      // include_granted_scopes
		url.QueryEscape("code"),      // response_type
		url.QueryEscape(redirectUri), // redirect_uri
		url.QueryEscape(g.clientID),  // client_id
	)

	return url
}

type GoogleAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

func (g *GoogleOauth) CompleteAuth(r *http.Request, redirectUri string) (GoogleAccessTokenResponse, error) {
	code := r.URL.Query().Get("code")

	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", g.clientID)
	data.Set("client_secret", g.clientSecret)
	data.Set("redirect_uri", redirectUri)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(googleTokenURL, data)
	if err != nil {
		return GoogleAccessTokenResponse{}, err
	}
	defer resp.Body.Close()

	var tokenResp GoogleAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return GoogleAccessTokenResponse{}, err
	}

	return tokenResp, nil
}

func (g *GoogleOauth) FetchUser(token string) (*User, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user data: %s", resp.Status)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user data: %s", err)
	}

	return &user, nil
}
