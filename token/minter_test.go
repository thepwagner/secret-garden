package token_test

import (
	"context"
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
	appID := os.Getenv("TPW_APP_ID")
	installIDRaw := os.Getenv("TPW_INSTALL_ID")
	pkPath := os.Getenv("TPW_APP_PK")
	if appID == "" || pkPath == "" || installIDRaw == "" {
		t.Skip()
	}
	installID, err := strconv.ParseInt(installIDRaw, 10, 64)
	require.NoError(t, err)
	pkEncodedBytes, err := ioutil.ReadFile(pkPath)
	require.NoError(t, err)
	pkBytes, _ := pem.Decode(pkEncodedBytes)
	require.NotNil(t, pkBytes)
	pk, err := x509.ParsePKCS1PrivateKey(pkBytes.Bytes)
	require.NoError(t, err)

	gh, err := token.NewAppClient(appID, pk)
	require.NoError(t, err)

	m := token.NewMinter(gh, installID)
	repos := []string{
		"thepwagner/dependabot-test-docker",
		"thepwagner/dependabot-test-npm",
	}
	tok, err := m.Mint(context.Background(), repos, &github.InstallationPermissions{
		Contents: github.String("read"),
	})
	require.NoError(t, err)
	t.Log(tok)
}
