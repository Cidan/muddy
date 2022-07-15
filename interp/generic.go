package interp

import (
	"context"
	"errors"
)

var (
	ErrNoCommand = errors.New("command not found")
)

type GenericInterp struct {
	commands map[string]*Command
}

func NewGenericInterp() Interp {
	return &GenericInterp{
		commands: make(map[string]*Command),
	}
}

func (g *GenericInterp) Process(ctx context.Context, player Player, command string, input ...string) error {
	if c, ok := g.commands[command]; !ok {
		return ErrNoCommand
	} else {
		return c.fn(ctx, player, input...)
	}
}

func (g *GenericInterp) Register(c *Command) {
	g.commands[c.name] = c
}

/*

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
*/
