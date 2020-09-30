package main

import (
	"context"
	"os"

	"github.com/drewbailey/nomad-deploy-notifier/bot"
	"github.com/hashicorp/go-hclog"
)

func main() {

	token := os.Getenv("SLACK_TOKEN")
	toChannel := os.Getenv("SLACK_CHANNEL")

	slackCfg := bot.Config{
		Token:   token,
		Channel: toChannel,
	}

	bot, err := NewBot(slackCfg)
	if err != nil {
		panic(err)
	}

	err := bot.Run(context.Background(), hclog.Default())

}
