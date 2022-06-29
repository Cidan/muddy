package interp

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Interp struct {
}

var (
	once   sync.Once
	interp *Interp
)

func Get() *Interp {
	once.Do(func() {
		interp = &Interp{}
	})

	return interp
}

func (i *Interp) Serve(ctx context.Context) error {
	log.Info().Msg("Starting Interpreter")
	for {
		time.Sleep(time.Second)
	}
}
