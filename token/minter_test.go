package token_test

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/google/go-github/v32/github"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/secret-garden/token"
)

func TestMinter_Mint(t *testing.T) {
	tc := newTokensClient(t)
	repos := []string{
		"thepwagner/dependabot-test-docker",
		"thepwagner/dependabot-test-npm",
	}
	tok, err := tc.Mint(context.Background(), repos, &github.InstallationPermissions{
		Contents: github.String("read"),
	})
	require.NoError(t, err)
	t.Log(tok)
}

func newTokensClient(t *testing.T) *token.TokensClient {
	appID, pk, installID := envAppCredentials(t)
	gh, err := token.NewAppClient(appID, pk)
	require.NoError(t, err)
	return token.NewTokensClient(gh, installID)
}

func envAppCredentials(t *testing.T) (appID string, appPK *rsa.PrivateKey, installationID int64) {
	appID = os.Getenv("TPW_APP_ID")
	installIDRaw := os.Getenv("TPW_INSTALL_ID")
	pkPath := os.Getenv("TPW_APP_PK")
	if appID == "" || pkPath == "" || installIDRaw == "" {
		t.Skip()
	}

	installationID, err := strconv.ParseInt(installIDRaw, 10, 64)
	require.NoError(t, err)

	pkEncodedBytes, err := ioutil.ReadFile(pkPath)
	require.NoError(t, err)
	pkBytes, _ := pem.Decode(pkEncodedBytes)
	require.NotNil(t, pkBytes)
	appPK, err = x509.ParsePKCS1PrivateKey(pkBytes.Bytes)
	require.NoError(t, err)
	return
}
