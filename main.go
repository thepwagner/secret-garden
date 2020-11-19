package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/thepwagner/action-update/actions"
	"github.com/thepwagner/secret-garden/action"
)

func main() {
	var cfg actions.Environment
	handlers := action.NewHandlers(&cfg)
	ctx := context.Background()
	if err := handlers.ParseAndHandle(ctx, &cfg); err != nil {
		logrus.WithError(err).Fatal("failed")
	}
}
