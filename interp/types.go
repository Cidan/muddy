package interp

import (
	"context"

	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
)

type commandCallback func(context.Context, Player, ...string) error

type Command struct {
	name   string
	interp playerv1.Player_InterpType
	alias  []string
	fn     commandCallback
}

type Input struct {
	Player Player
	Text   string
}

type Player interface {
	Send(string, ...interface{}) error
	Interp() playerv1.Player_InterpType
}

type Interp interface {
	Process(context.Context, Player, string, ...string) error
	Register(*Command)
}
