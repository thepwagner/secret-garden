package action

import (
	"context"

	"github.com/google/go-github/v32/github"
	"github.com/thepwagner/action-update/actions"
)

func NewHandlers(env *Environment) *actions.Handlers {
	return &actions.Handlers{
		WorkflowDispatch: func(ctx context.Context) error {
			gh, err := env.NewTokensClient()
			if err != nil {
				return err
			}

			// FIXME: config from input
			token, err := gh.Mint(ctx, []string{"thepwagner-org/private"}, &github.InstallationPermissions{
				Contents: github.String("read"),
			})
			if err != nil {
				return err
			}
			return gh.StoreRepo(ctx, "thepwagner-org", "secret-garden", "TEST_TOKEN", token)
		},
	}
}
