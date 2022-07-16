package main

import (
	"context"
	"os"
	"time"

	_ "github.com/Cidan/muddy/commands"
	"github.com/Cidan/muddy/interp"
	"github.com/Cidan/muddy/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	Start()
}

func Start() {
	// Enable pretty logging while in dev.
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	ctx := context.Background()
	go server.Get().Serve(ctx)
	go interp.Get().Serve(ctx)
	time.Sleep(time.Second * 60)
}
