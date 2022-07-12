package interp

import (
	"context"
	"errors"

	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
)

var (
	ErrNoCommand = errors.New("command not found")
)

type commandCallback func(context.Context, Player, ...string) error

type CommandRef struct {
	name   string
	interp playerv1.Player_InterpType
	alias  []string
	fn     commandCallback
}

type Interp struct {
	commands map[string]*CommandRef
}

func newInterp() *Interp {
	return &Interp{
		commands: make(map[string]*CommandRef),
	}
}

func (i *Interp) Register(r *CommandRef) {
	i.commands[r.name] = r
}

func (i *Interp) Process(ctx context.Context, player Player, command string, input ...string) error {
	if i.commands[command] != nil {
		return i.commands[command].fn(ctx, player, input...)
	}
	return ErrNoCommand
}
