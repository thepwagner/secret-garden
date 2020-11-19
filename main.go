package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/thepwagner/secret-garden/action"
)

func main() {
	var cfg action.Environment
	handlers := action.NewHandlers(&cfg)
	ctx := context.Background()
	if err := handlers.ParseAndHandle(ctx, &cfg); err != nil {
		logrus.WithError(err).Fatal("failed")
	}
}
