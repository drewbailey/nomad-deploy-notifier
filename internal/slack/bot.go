package bot

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/slack-go/slack"
)

type Config struct {
	Token string
	Channel string
}

type Bot struct {
	chanName string
	api      *slack.Client
}

func NewBot(cfg Config) (*Bot, error) {
	if cfg.Token == "" {
		return nil, fmt.Errorf("no token provided")
	}

	api := slack.New(cfg.Token)

	bot := &Bot{
		api: api,
	}

	return bot, nil
}

func (b *Bot) Run(ctx context.Context, L hclog.Logger) error {
	rtm := b.api.NewRTM()

	go rtm.ManageConnection()
	defer rtm.Disconnect()

	channel, err := rtm.CreateConversationContext(ctx, b.chanName, false)
	if err != nil {
		return fmt.Errorf("creating conversation rtm: %w", err)
	}

	respCh, timestamp, err := rtm.PostMessage(
		channel.ID,
		slack.MsgOptionIconEmoji(":nomad-loading:"),
		slack.MsgOptionTS("")
	if err != nil {
		return fmt.Errorf("sending message %w", err)
	}

}
