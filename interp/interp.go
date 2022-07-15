package interp

import (
	"context"
	"strings"
	"sync"

	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

type Handler struct {
	input   chan *Input
	interps map[playerv1.Player_InterpType]Interp
}

var (
	once   sync.Once
	interp *Handler
)

func Get() *Handler {
	once.Do(func() {
		interp = &Handler{
			input:   make(chan *Input),
			interps: make(map[playerv1.Player_InterpType]Interp),
		}
		//interp.interps[playerv1.Player_INTERP_TYPE_PLAYING] = newInterp()
	})
	return interp
}

func (i *Handler) Serve(ctx context.Context) error {
	log.Info().Msg("Starting Interpreter")
	for {
		select {
		case <-ctx.Done():
			return suture.ErrDoNotRestart
		case in := <-i.input:
			all := strings.SplitN(in.Text, " ", 2)
			log.Debug().Strs("input", all).Msg("got command")
			err := i.interps[in.Player.Interp()].Process(ctx, in.Player, all[0], all[1:]...)
			switch err {
			case ErrNoCommand:
				in.Player.Send("huh?")
			}
		}
	}
}

func (i *Handler) Do(in *Input) {
	i.input <- in
}

func (i *Handler) Register(r *Command) {
	i.interps[r.interp].Register(r)
}
