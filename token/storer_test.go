package token_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorer_StoreRepo(t *testing.T) {
	tc := newTokensClient(t)
	const (
		owner = "thepwagner"
		repo  = "dependabot-test-npm"
		name  = "THEODOSIA"
		token = "dear"
	)
	err := tc.StoreRepo(context.Background(), owner, repo, name, token)
	require.NoError(t, err)
}
