package interp

import (
	"context"
	"errors"
	"strings"
	"sync"

	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
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
	state []playerv1.Player_StateType
	fn    commandCallback
}

type commandCallback func(context.Context, Player, ...string) error

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
			input:    make(chan *Command),
			commands: make(map[string]*CommandRef),
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
			err := i.Process(ctx, c.Player, all[0], all[1:]...)
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
	i.commands[c.name] = c
}

func (i *Interp) Process(ctx context.Context, player Player, command string, input ...string) error {
	if i.commands[command] != nil {
		return i.commands[command].fn(ctx, player, input...)
	}
	return ErrNoCommand
}
