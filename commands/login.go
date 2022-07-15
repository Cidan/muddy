package commands

import (
	"context"

	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
	"github.com/Cidan/muddy/interp"
	"github.com/rs/zerolog/log"
)

type Login struct{}

func (l *Login) Process(ctx context.Context, player interp.Player, input string, args ...string) error {
	// TODO(lobato): Login process here, state machine, etc.
	player.Send("test :)")
	return nil
}

func (l *Login) Register(c *interp.Command) {}

func init() {
	log.Debug().Msg("registering login interp")
	interp.Get().Set(playerv1.Player_INTERP_TYPE_LOGIN, &Login{})
}
