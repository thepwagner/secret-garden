package action

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/thepwagner/action-update/actions"
	"github.com/thepwagner/secret-garden/token"
)

type Environment struct {
	actions.Environment

	AppID            int64  `env:"INPUT_APP_ID"`
	AppPrivateKeyPEM []byte `env:"INPUT_APP_PRIVATE_KEY"`
	InstallationID   int64  `env:"INPUT_INSTALLATION_ID"`
}

func (e *Environment) NewTokensClient() (*token.TokensClient, error) {
	pkBytes, _ := pem.Decode(e.AppPrivateKeyPEM)
	if pkBytes == nil {
		return nil, fmt.Errorf("invalid private key PEM")
	}
	pk, err := x509.ParsePKCS1PrivateKey(pkBytes.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing private key: %w", err)
	}

	gh, err := token.NewAppClient(e.AppID, pk)
	if err != nil {
		return nil, fmt.Errorf("creating app client: %w", err)
	}
	return token.NewTokensClient(gh, e.InstallationID), nil
}
