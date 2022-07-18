package commands

import (
	"context"
	"strings"

	"github.com/Cidan/muddy/atlas"
	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
	"github.com/Cidan/muddy/interp"
	"github.com/rs/zerolog/log"
)

type Login struct {
}

func (l *Login) Process(ctx context.Context, player atlas.Player, input string, args ...string) error {
	switch player.LoginState() {
	case "ASK_NAME":
		l.AskName(ctx, player, input)
	case "CONFIRM_NAME":
		l.ConfirmName(ctx, player, input)
	case "NEW_PASSWORD":
		l.NewPassword(ctx, player, input)
	case "CONFIRM_PASSWORD":
		l.ConfirmPassword(ctx, player, input)
	default:
		log.Warn().Str("state", player.LoginState()).Msg("player login state is invalid and not handled.")
		player.Send("You're in a weird, invalid state. Contact the MUD administrators.")
	}
	return nil
}

func (l *Login) Register(c *interp.Command) {}

func (l *Login) AskName(ctx context.Context, player atlas.Player, input string) {
	player.SetName(input)
	player.Send("Are you sure you want your name to be %s?", input)
	player.SetLoginState("CONFIRM_NAME")
}

func (l *Login) ConfirmName(ctx context.Context, player atlas.Player, input string) {
	if strings.HasPrefix(strings.ToLower(input), "y") {
		player.Send("Welcome, %s!", player.Name())
		player.Send("Please choose a password: ")
		player.SetLoginState("NEW_PASSWORD")
		return
	}

	player.Send("Okay, what would you like your name to be?")
	player.SetName("")
	player.SetLoginState("ASK_NAME")
}

func (l *Login) NewPassword(ctx context.Context, player atlas.Player, input string) {
	player.SetPassword(input)
	player.Send("Type your password again to confirm: ")
	player.SetLoginState("CONFIRM_PASSWORD")
}

func (l *Login) ConfirmPassword(ctx context.Context, player atlas.Player, input string) {
	if !player.CheckPassword(input) {
		player.Send("Password's do not match, please re-enter your password: ")
		player.SetLoginState("NEW_PASSWORD")
		return
	}
	player.Send("Password confirmed, welcome!")
	player.SetLoginState("")
	player.SetInterp(playerv1.Player_INTERP_TYPE_PLAYING)
}

func init() {
	log.Debug().Msg("registering login interp")
	l := &Login{}
	interp.Get().Set(playerv1.Player_INTERP_TYPE_LOGIN, l)
}
