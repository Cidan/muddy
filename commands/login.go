package commands

import (
	"context"
	"strings"

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
	case "CONFIRM_NAME":
		l.ConfirmName(ctx, player, input)
	}
	return nil
}

func (l *Login) Register(c *interp.Command) {}

func (l *Login) AskName(ctx context.Context, player interp.Player, input string) {
	player.SetName(input)
	player.Send("Are you sure you want your name to be %s?", input)
	player.SetLoginState("CONFIRM_NAME")
}

func (l *Login) ConfirmName(ctx context.Context, player interp.Player, input string) {
	if strings.HasPrefix(strings.ToLower(input), "y") {
		player.Send("Welcome, %s!", player.Name())
		player.SetInterp(playerv1.Player_INTERP_TYPE_PLAYING)
		return
	}

	player.Send("Okay, what would you like your name to be?")
	player.SetName("")
	player.SetLoginState("ASK_NAME")
}

func init() {
	log.Debug().Msg("registering login interp")
	l := &Login{}
	interp.Get().Set(playerv1.Player_INTERP_TYPE_LOGIN, l)
}
