package oauth

import (
	"github.com/joho/godotenv"
)

const (
	googleTokenURL = "https://accounts.google.com/o/oauth2/token"
	googleAuthURL  = "https://accounts.google.com/o/oauth2/auth"
)

type Oauth struct {
	Google *GoogleOauth
}

type User struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Email       string `json:"email"`
    Picture     string `json:"picture"`
}

func New(envFiles ...string) (*Oauth, error) {
	err := godotenv.Load(envFiles...)
	if err != nil {
		return nil, err
	}

	return &Oauth{
		Google: newGoogleOauth(),
	}, nil
}
