package action_test

import (
	"context"
	"testing"

	"github.com/google/go-github/v32/github"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/secret-garden/action"
)

func TestProcessConfig(t *testing.T) {
	p := action.NewConfigProcessor("thepwagner-org", &mockMinter{}, &mockStorer{})

	err := p.ProcessConfig(context.Background(), "../config")
	require.NoError(t, err)
}

type mockMinter struct{}

func (m *mockMinter) Mint(_ context.Context, repoFullNames []string, perms *github.InstallationPermissions) (string, error) {
	logrus.WithFields(logrus.Fields{
		"repos":    repoFullNames,
		"contents": perms.GetContents(),
	}).Info("minting token")
	return "awesome-token", nil
}

type mockStorer struct{}

func (m mockStorer) StoreRepo(_ context.Context, owner, repo, name, token string) error {
	logrus.WithFields(logrus.Fields{
		"owner": owner,
		"repo":  repo,
		"name":  name,
	}).Info("storing token to repo")
	return nil
}

func (m mockStorer) StoreOrg(ctx context.Context, owner, name, token string, consumers []string) error {
	logrus.WithFields(logrus.Fields{
		"owner":     owner,
		"name":      name,
		"consumers": consumers,
	}).Info("storing token to org")
	return nil
}
