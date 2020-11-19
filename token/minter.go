package token

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v32/github"
)

type Minter struct {
	gh             *github.Client
	installationID int64
}

func NewMinter(gh *github.Client, installationID int64) *Minter {
	return &Minter{gh: gh, installationID: installationID}
}

func (m *Minter) Mint(ctx context.Context, repoFullNames []string, perms *github.InstallationPermissions) (string, error) {
	repoIDs, err := m.resolveRepoIDs(ctx, repoFullNames)
	if err != nil {
		return "", err
	}

	// Return token scoped to the target repositories and permissions:
	token, _, err := m.gh.Apps.CreateInstallationToken(ctx, m.installationID, &github.InstallationTokenOptions{
		RepositoryIDs: repoIDs,
		Permissions:   perms,
	})
	if err != nil {
		return "", err
	}
	return token.GetToken(), nil
}

func (m *Minter) resolveRepoIDs(ctx context.Context, repoFullNames []string) ([]int64, error) {
	repos, err := m.listInstallationRepos(ctx)
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

func (m *Minter) listInstallationRepos(ctx context.Context) (map[string]int64, error) {
	listToken, _, err := m.gh.Apps.CreateInstallationToken(ctx, m.installationID, &github.InstallationTokenOptions{})
	if err != nil {
		return nil, fmt.Errorf("generating token for repo list: %w", err)
	}
	listClient, err := newGitHubClient(listToken.GetToken())
	if err != nil {
		return nil, fmt.Errorf("creating client for repo list: %w", err)
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
