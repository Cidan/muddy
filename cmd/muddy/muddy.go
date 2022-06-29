package main

import (
	"context"
	"os"

	"github.com/Cidan/muddy/interp"
	"github.com/Cidan/muddy/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

func main() {
	// Enable pretty logging while in dev.
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	sup := suture.New("muddy", suture.Spec{
		PassThroughPanics: true,
	})

	sup.Add(server.Get())
	sup.Add(server.Get().Supervisor())
	sup.Add(interp.Get())

	if err := sup.Serve(context.Background()); err != nil {
		panic(err)
	}

}
