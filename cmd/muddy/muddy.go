package main

import (
	"context"
	"os"

	_ "github.com/Cidan/muddy/commands"
	"github.com/Cidan/muddy/interp"
	"github.com/Cidan/muddy/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	Start(context.Background())
}

func Start(ctx context.Context) {
	// Enable pretty logging while in dev.
	// TODO(lobato): figure out why this races when testing
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	go server.Get().Serve(ctx)
	go interp.Get().Serve(ctx)
	<-ctx.Done()
	// TODO(lobato): wait for sigint code

}
