package interp

import (
	"context"
	"strings"
	"sync"

	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

type Player interface {
	Send(string, ...interface{}) error
	Interp() playerv1.Player_InterpType
}

type Command struct {
	Player Player
	Text   string
}

type Handler struct {
	input   chan *Command
	interps map[playerv1.Player_InterpType]*Interp
}

var (
	once   sync.Once
	interp *Handler
)

func Get() *Handler {
	once.Do(func() {
		interp = &Handler{
			input:   make(chan *Command),
			interps: make(map[playerv1.Player_InterpType]*Interp),
		}
		interp.interps[playerv1.Player_INTERP_TYPE_PLAYING] = newInterp()
	})
	return interp
}

func (i *Handler) Serve(ctx context.Context) error {
	log.Info().Msg("Starting Interpreter")
	for {
		select {
		case <-ctx.Done():
			return suture.ErrDoNotRestart
		case c := <-i.input:
			all := strings.SplitN(c.Text, " ", 2)
			log.Debug().Strs("input", all).Msg("got command")
			err := i.interps[c.Player.Interp()].Process(ctx, c.Player, all[0], all[1:]...)
			switch err {
			case ErrNoCommand:
				c.Player.Send("huh?")
			}
		}
	}
}

func (i *Handler) Do(c *Command) {
	i.input <- c
}

func (i *Handler) Register(r *CommandRef) {
	i.interps[r.interp].Register(r)
}
