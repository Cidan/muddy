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

type Handler struct {
	input   chan *Input
	interps map[playerv1.Player_InterpType]Interp
	lock    sync.RWMutex
}

var (
	once         sync.Once
	interp       *Handler
	ErrNoInterp  = errors.New("interp not set for this interp type")
	ErrNoCommand = errors.New("no such command")
)

func Get() *Handler {
	once.Do(func() {
		interp = &Handler{
			input:   make(chan *Input),
			interps: make(map[playerv1.Player_InterpType]Interp),
			lock:    sync.RWMutex{},
		}
	})
	return interp
}

func (h *Handler) Serve(ctx context.Context) error {
	// TODO(lobato): potentially move this to a worker pool so
	// that multiple commands can run concurrently.
	log.Info().Msg("Starting Interpreter")
	for {
		select {
		case <-ctx.Done():
			return suture.ErrDoNotRestart
		case in := <-h.input:
			all := strings.SplitN(in.Text, " ", 2)

			h.lock.RLock()
			interp, ok := h.interps[in.Player.Interp()]
			h.lock.RUnlock()

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

func (h *Handler) Do(in *Input) {
	h.input <- in
}

func (h *Handler) Set(interpType playerv1.Player_InterpType, interp Interp) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.interps[interpType] = interp
}
