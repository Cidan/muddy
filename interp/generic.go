package interp

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNoCommand = errors.New("command not found")
)

type GenericInterp struct {
	commands map[string]*Command
	lock     sync.RWMutex
}

func NewGenericInterp() Interp {
	return &GenericInterp{
		commands: make(map[string]*Command),
		lock:     sync.RWMutex{},
	}
}

func (g *GenericInterp) Process(ctx context.Context, player Player, command string, input ...string) error {
	g.lock.RLock()
	c, ok := g.commands[command]
	g.lock.RUnlock()

	if !ok {
		return ErrNoCommand
	} else {
		return c.Fn(ctx, player, input...)
	}
}

func (g *GenericInterp) Register(c *Command) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.commands[c.Name] = c
}
