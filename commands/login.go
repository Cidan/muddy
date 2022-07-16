package commands

import (
	"context"

	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
	"github.com/Cidan/muddy/interp"
	"github.com/rs/zerolog/log"
)

type Login struct {
}

func (l *Login) Process(ctx context.Context, player interp.Player, input string, args ...string) error {
	switch player.LoginState() {
	case "ASK_NAME":
		l.AskName(ctx, player, input)
	}
	return nil
}

func (l *Login) Register(c *interp.Command) {}

func (l *Login) AskName(ctx context.Context, player interp.Player, input string) {
	player.Send("Are you sure you want your name to be %s?", input)
}

func init() {
	log.Debug().Msg("registering login interp")
	l := &Login{}
	interp.Get().Set(playerv1.Player_INTERP_TYPE_LOGIN, l)
}
