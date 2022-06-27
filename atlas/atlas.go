package atlas

import (
	"context"
	"sync"
	"time"

	ppb "github.com/Cidan/muddy/gen/proto/go/player/v1"
	"github.com/rs/zerolog/log"
)

var (
	once  sync.Once
	atlas *Atlas
)

func Get() *Atlas {
	once.Do(func() {
		atlas = &Atlas{}
	})
	return atlas
}

type Atlas struct{}

func (a *Atlas) Serve(ctx context.Context) error {
	p := &ppb.Player{}
	_ = p
	log.Info().Msg("Starting Atlas")
	// TODO(lobato): Do something real with atlas here.
	for {
		time.Sleep(1 * time.Hour)
	}
}
