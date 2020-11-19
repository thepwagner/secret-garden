package token_test

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-github/v32/github"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/secret-garden/token"
)

func TestMinter_Mint(t *testing.T) {
	appID := os.Getenv("TPW_APP_ID")
	pkPath := os.Getenv("TPW_APP_PK")
	if appID == "" || pkPath == "" {
		t.Skip()
	}
	pkEncodedBytes, err := ioutil.ReadFile(pkPath)
	require.NoError(t, err)
	pkBytes, _ := pem.Decode(pkEncodedBytes)
	require.NotNil(t, pkBytes)
	pk, err := x509.ParsePKCS1PrivateKey(pkBytes.Bytes)
	require.NoError(t, err)

	gh, err := token.NewAppClient(appID, pk)
	require.NoError(t, err)
	m := token.NewMinter(gh)

	var repoIDs []int64
	tok, err := m.Mint(context.Background(), repoIDs, &github.InstallationPermissions{
		Contents: github.String("read"),
	})
	require.NoError(t, err)
	t.Log(tok)
}
