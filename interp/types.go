package interp

import (
	"context"
	"errors"

	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
)

var (
	ErrNoCommand = errors.New("no such command")
)

type commandCallback func(context.Context, Player, ...string) error

type Command struct {
	Name  string
	Alias []string
	Fn    commandCallback
}

type Input struct {
	Player Player
	Text   string
}

type Player interface {
	Send(string, ...interface{}) error
	LoginState() string
	SetLoginState(string)
	SetInterp(playerv1.Player_InterpType)
	Interp() playerv1.Player_InterpType
	SetName(string)
	Name() string
	SetPassword(string)
	CheckPassword(string) bool
	Password() string
}

type Interp interface {
	Process(context.Context, Player, string, ...string) error
	Register(*Command)
}
