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
	lock    sync.RWMutex
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

			i.lock.RLock()
			interp, ok := i.interps[in.Player.Interp()]
			i.lock.RUnlock()

			if !ok {
				log.Warn().Msg("player sent a command with an invalid interp")
				in.Player.Send("huh?")
				break
			}

			err := interp.Process(ctx, in.Player, all[0], all[1:]...)
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

func (i *Handler) Add(interpType playerv1.Player_InterpType, interp Interp) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.interps[interpType] = interp
}
