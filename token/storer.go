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
	// TODO: org
}

var _ Storer = (*TokensClient)(nil)

func (t *TokensClient) StoreRepo(ctx context.Context, owner, repo, name, token string) error {
	repoClient, err := t.newClient(ctx)
	if err != nil {
		return fmt.Errorf("preparing client for list: %w", err)
	}

	repoKey, _, err := repoClient.Actions.GetRepoPublicKey(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("fetching repo public repoKey: %w", err)
	}

	key, err := base64.StdEncoding.DecodeString(repoKey.GetKey())
	if err != nil {
		return fmt.Errorf("decoding public repoKey: %v", err)
	}
	encrypted, exit := sodium.CryptoBoxSeal([]byte(token), key)
	if exit != 0 {
		return fmt.Errorf("encrypting secret failed")
	}

	_, err = repoClient.Actions.CreateOrUpdateRepoSecret(ctx, owner, repo, &github.EncryptedSecret{
		Name:           name,
		KeyID:          repoKey.GetKeyID(),
		EncryptedValue: base64.StdEncoding.EncodeToString(encrypted),
	})
	if err != nil {
		return fmt.Errorf("storing repo secret: %w", err)
	}
	return nil
}
