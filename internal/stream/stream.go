package stream

import (
	"context"
	"reflect"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/drewbailey/nomad-deploy-notifier/internal/bot"
	"github.com/hashicorp/nomad/api"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
}

type Stream struct {
	nomad *api.Client
}

func NewStream(cfg Config) *Stream {
	client, _ := api.NewClient(&api.Config{})
	return &Stream{
		nomad: client,
	}
}

type DeployMsg struct {
	Deployment api.Deployment `mapstructure:"Deployment"`
}

func (s *Stream) Subscribe(ctx context.Context, slack *bot.Bot) {

	events := s.nomad.EventStream()

	topics := map[api.Topic][]string{
		api.Topic("Deployment"): {"*"},
	}
	eventCh, errCh := events.Stream(ctx, topics, 0, &api.QueryOptions{})

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-errCh:
			spew.Dump(err)
			return
		case event := <-eventCh:
			if event.IsHeartBeat() {
				continue
			}

			for _, e := range event.Events {

				// decode fully to ensure we received expected out
				var out DeployMsg
				var md mapstructure.Metadata
				cfg := &mapstructure.DecoderConfig{
					DecodeHook: mapstructure.ComposeDecodeHookFunc(
						ToTimeHookFunc(),
					),
					Metadata: &md,
					Result:   &out,
				}
				dec, err := mapstructure.NewDecoder(cfg)
				if err != nil {
					spew.Dump(err)
					return
				}
				if err := dec.Decode(e.Payload); err != nil {
					spew.Dump("ERROR", err)
					return
				}

				if err := slack.UpsertDeployMsg(out.Deployment); err != nil {
					spew.Dump("ERROR", err)
					return
				}
			}
		}
	}
}

func ToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
		// Convert it by parsing
	}
}
