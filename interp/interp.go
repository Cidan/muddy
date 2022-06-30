package interp

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

type Player interface {
	Send(string, ...interface{}) error
}

type Command struct {
	Player Player
	Text   string
}

type CommandRef struct {
	name  string
	alias []string
	state []string
	fn    commandCallback
}

type commandCallback func(context.Context, ...string) error

type Interp struct {
	input    chan *Command
	commands map[string]*CommandRef
}

var (
	once         sync.Once
	interp       *Interp
	ErrNoCommand = errors.New("command not found")
)

func Get() *Interp {
	once.Do(func() {
		interp = &Interp{
			input: make(chan *Command),
		}
	})

	return interp
}

func (i *Interp) Serve(ctx context.Context) error {
	log.Info().Msg("Starting Interpreter")
	for {
		select {
		case <-ctx.Done():
			return suture.ErrDoNotRestart
		case c := <-i.input:
			all := strings.SplitN(c.Text, " ", 2)
			log.Debug().Strs("input", all).Msg("got command")
			err := i.Process(ctx, all[0], all[1:]...)
			switch err {
			case ErrNoCommand:
				c.Player.Send("huh?")
			}
		}
	}
}

func (i *Interp) Do(c *Command) {
	i.input <- c
}

func (i *Interp) Register(c *CommandRef) {
	for _, alias := range c.alias {
		i.commands[alias] = c
	}
	i.commands[c.name] = c
}

func (i *Interp) Process(ctx context.Context, command string, input ...string) error {
	if i.commands[command] != nil {
		return i.commands[command].fn(ctx, input...)
	}
	return ErrNoCommand
}
