package action

import (
	"context"
	"strings"

	"github.com/thepwagner/action-update/actions"
)

func NewHandlers(env *Environment) *actions.Handlers {
	return &actions.Handlers{
		WorkflowDispatch: func(ctx context.Context) error {
			tc, err := env.NewTokensClient()
			if err != nil {
				return err
			}

			repoSplit := strings.SplitN(env.GitHubRepository, "/", 2)
			org := repoSplit[0]

			c := NewConfigProcessor(org, tc, tc)
			return c.ProcessConfig(ctx, "config")
		},
	}
}
