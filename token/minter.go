package token

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v32/github"
	"github.com/sirupsen/logrus"
)

type Minter interface {
	Mint(ctx context.Context, repoFullNames []string, perms *github.InstallationPermissions) (string, error)
}

var _ Minter = (*TokensClient)(nil)

type TokensClient struct {
	gh             *github.Client
	installationID int64
}

func NewTokensClient(gh *github.Client, installationID int64) *TokensClient {
	return &TokensClient{gh: gh, installationID: installationID}
}

func (t *TokensClient) Mint(ctx context.Context, repoFullNames []string, perms *github.InstallationPermissions) (string, error) {
	repoIDs, err := t.resolveRepoIDs(ctx, repoFullNames)
	if err != nil {
		return "", err
	}

	// Return token scoped to the target repositories and permissions:
	token, _, err := t.gh.Apps.CreateInstallationToken(ctx, t.installationID, &github.InstallationTokenOptions{
		RepositoryIDs: repoIDs,
		Permissions:   perms,
	})
	if err != nil {
		return "", err
	}
	logrus.WithFields(logrus.Fields{
		"repo_ids":   repoIDs,
		"expires_at": token.GetExpiresAt(),
	}).Info("issued token")

	token.GetExpiresAt()
	return token.GetToken(), nil
}

func (t *TokensClient) resolveRepoIDs(ctx context.Context, repoFullNames []string) ([]int64, error) {
	if len(repoFullNames) == 0 {
		return nil, nil
	}

	repos, err := t.listInstallationRepos(ctx)
	if err != nil {
		return nil, err
	}

	repoIDs := make([]int64, 0, len(repoFullNames))
	var missing []string
	for _, fullName := range repoFullNames {
		repoID, ok := repos[fullName]
		if ok {
			repoIDs = append(repoIDs, repoID)
		} else {
			missing = append(missing, fullName)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("installation can not access repositories: %s", strings.Join(missing, ", "))
	}
	return repoIDs, nil
}

func (t *TokensClient) listInstallationRepos(ctx context.Context) (map[string]int64, error) {
	listClient, err := t.newClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("preparing client for list: %w", err)
	}
	// TODO: mo pages, _less_ problems
	repos, _, err := listClient.Apps.ListRepos(ctx, &github.ListOptions{
		PerPage: 100,
	})
	if err != nil {
		return nil, fmt.Errorf("listing installation repos: %w", err)
	}

	index := make(map[string]int64, len(repos))
	for _, r := range repos {
		index[r.GetFullName()] = r.GetID()
	}
	return index, nil
}

func (t *TokensClient) newClient(ctx context.Context) (*github.Client, error) {
	token, _, err := t.gh.Apps.CreateInstallationToken(ctx, t.installationID, &github.InstallationTokenOptions{})
	if err != nil {
		return nil, fmt.Errorf("generating token: %w", err)
	}
	client, err := newGitHubClient(token.GetToken())
	if err != nil {
		return nil, fmt.Errorf("creating client: %w", err)
	}
	return client, nil
}
