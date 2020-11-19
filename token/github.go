package token

import (
	"context"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func NewAppClient(appID int64, key *rsa.PrivateKey) (*github.Client, error) {
	const jwtDuration = 5 * time.Minute
	now := time.Now()
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, &jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(jwtDuration).Unix(),
		Issuer:    fmt.Sprintf("%d", appID),
	})
	signed, err := tok.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("signing JWT: %w", err)
	}

	return newGitHubClient(signed)
}

func newGitHubClient(token string) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return github.NewClient(oauth2.NewClient(ctx, ts)), nil
}
