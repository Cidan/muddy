package main

import (
	"context"
	"os"

	"github.com/Cidan/muddy/atlas"
	_ "github.com/Cidan/muddy/commands"
	"github.com/Cidan/muddy/construct"
	roomv1 "github.com/Cidan/muddy/gen/proto/go/room/v1"
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
	// Probably because of init side effects in commands import
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	makeWorld()
	go server.Get().Serve(ctx)
	go interp.Get().Serve(ctx)
	<-ctx.Done()
	// TODO(lobato): wait for sigint code

}

func makeWorld() {
	log.Debug().Msg("making starting world")
	r1 := construct.NewRoom()
	r1.SetCoordinates(&roomv1.Room_Coordinates{
		X: 0,
		Y: 0,
		Z: 0,
	})
	r1.SetName("The Alpha")

	r2 := construct.NewRoom()
	r2.SetCoordinates(&roomv1.Room_Coordinates{
		X: 1,
		Y: 0,
		Z: 0,
	})
	r2.SetName("The Omega")

	atlas.Get().AddRoom(r1)
	atlas.Get().AddRoom(r2)
}
