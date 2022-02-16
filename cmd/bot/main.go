package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/drewbailey/nomad-deploy-notifier/internal/bot"
	"github.com/drewbailey/nomad-deploy-notifier/internal/stream"
)

func main() {
	os.Exit(realMain(os.Args))
}

func realMain(args []string) int {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	token := os.Getenv("SLACK_TOKEN")
	toChannel := os.Getenv("SLACK_CHANNEL")

	slackCfg := bot.Config{
		Token:   token,
		Channel: toChannel,
	}

	stream := stream.NewStream()

	slackBot, err := bot.NewBot(slackCfg)
	if err != nil {
		panic(err)
	}

	stream.Subscribe(ctx, slackBot)

	return 0
}
