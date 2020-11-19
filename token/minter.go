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

type Minter struct {
	// TODO: interface
	gh *github.Client
}

func NewAppClient(appID string, key *rsa.PrivateKey) (*github.Client, error) {
	const jwtDuration = 5 * time.Minute
	now := time.Now()
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, &jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(jwtDuration).Unix(),
		Issuer:    appID,
	})
	signed, err := tok.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("signing JWT: %w", err)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: signed})
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return github.NewClient(oauth2.NewClient(ctx, ts)), nil
}

func NewMinter(gh *github.Client) *Minter {
	return &Minter{gh: gh}
}

func (m *Minter) Mint(ctx context.Context, repoIDs []int64, perms *github.InstallationPermissions) (string, error) {
	// repoIDs
	// - must not contain public repos
	// - must be part of this installation
	const installID = 123
	token, _, err := m.gh.Apps.CreateInstallationToken(ctx, installID, &github.InstallationTokenOptions{
		RepositoryIDs: repoIDs,
		Permissions:   perms,
	})
	if err != nil {
		return "", err
	}
	return token.GetToken(), nil
}
