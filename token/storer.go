package token

import (
	"context"
	"encoding/base64"
	"fmt"

	sodium "github.com/GoKillers/libsodium-go/cryptobox"
	"github.com/google/go-github/v32/github"
)

type Storer interface {
	StoreRepo(ctx context.Context, owner, repo, name, token string) error
	StoreOrg(ctx context.Context, owner, name, token string, consumers []string) error
}

var _ Storer = (*TokensClient)(nil)

func (t *TokensClient) StoreRepo(ctx context.Context, owner, repo, name, token string) error {
	gh, err := t.newClient(ctx)
	if err != nil {
		return fmt.Errorf("preparing repo client: %w", err)
	}

	key, _, err := gh.Actions.GetRepoPublicKey(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("fetching repo public key: %w", err)
	}

	encrypted, err := newEncryptedSecret(key, name, token)
	if err != nil {
		return err
	}

	_, err = gh.Actions.CreateOrUpdateRepoSecret(ctx, owner, repo, encrypted)
	if err != nil {
		return fmt.Errorf("storing repo secret: %w", err)
	}
	return nil
}

func (t *TokensClient) StoreOrg(ctx context.Context, owner, name, token string, consumers []string) error {
	gh, err := t.newClient(ctx)
	if err != nil {
		return fmt.Errorf("preparing org client: %w", err)
	}

	key, _, err := gh.Actions.GetOrgPublicKey(ctx, owner)
	if err != nil {
		return fmt.Errorf("fetching org public key: %w", err)
	}

	encrypted, err := newEncryptedSecret(key, name, token)
	if err != nil {
		return err
	}

	if len(consumers) == 0 {
		encrypted.Visibility = "private"
	} else {
		encrypted.Visibility = "selected"

		selectedRepoIDs, err := t.resolveRepoIDs(ctx, consumers)
		if err != nil {
			return fmt.Errorf("resolving consumer repos: %w", err)
		}
		encrypted.SelectedRepositoryIDs = selectedRepoIDs
	}

	_, err = gh.Actions.CreateOrUpdateOrgSecret(ctx, owner, encrypted)
	if err != nil {
		return fmt.Errorf("storing org secret: %w", err)
	}
	return nil
}

func newEncryptedSecret(key *github.PublicKey, name, token string) (*github.EncryptedSecret, error) {
	decoded, err := base64.StdEncoding.DecodeString(key.GetKey())
	if err != nil {
		return nil, fmt.Errorf("decoding public key: %v", err)
	}
	encrypted, exit := sodium.CryptoBoxSeal([]byte(token), decoded)
	if exit != 0 {
		return nil, fmt.Errorf("encrypting secret failed")
	}
	return &github.EncryptedSecret{
		Name:           name,
		KeyID:          key.GetKeyID(),
		EncryptedValue: base64.StdEncoding.EncodeToString(encrypted),
	}, nil
}
