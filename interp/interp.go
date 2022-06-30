package interp

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

type Player interface{}

type Command struct {
	Player Player
	Text   string
}

type Interp struct {
	input chan *Command
}

var (
	once   sync.Once
	interp *Interp
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
			log.Debug().Str("input", c.Text).Msg("got command")
		}
	}
}

func (i *Interp) Do(c *Command) {
	i.input <- c
}
