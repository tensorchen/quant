package discord

import (
	"context"
	"os"
	"strconv"

	"github.com/tensorchen/quant/entity"
	"github.com/tensorchen/quant/env"
	"github.com/tensorchen/quant/notify"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
)

var _ notify.Notifier = (*Notifier)(nil)

type Notifier struct {
	cli webhook.Client
}

func (n *Notifier) Notify(ctx context.Context, information entity.Information) error {
	_, err := n.cli.CreateEmbeds([]discord.Embed{discord.Embed(information)})
	return err
}

func New() (*Notifier, error) {
	id, err := strconv.Atoi(os.Getenv(env.DiscordIDKey))
	if err != nil {
		return nil, err
	}
	cli := webhook.New(snowflake.ID(id), os.Getenv(env.DiscordTokenKey))
	return &Notifier{cli: cli}, nil
}
