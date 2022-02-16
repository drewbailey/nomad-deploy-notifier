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
	ctx, closer := CtxWithInterrupt(context.Background())
	defer closer()

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

func CtxWithInterrupt(ctx context.Context) (context.Context, func()) {

	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	return ctx, func() {
		signal.Stop(ch)
		cancel()
	}
}
