package main

import (
	"context"
	"testing"
	"time"

	"github.com/Cidan/muddy/interp"
	"github.com/Cidan/muddy/server"
)

func TestAll(t *testing.T) {
	/*
		sup := suture.New("muddy", suture.Spec{
			PassThroughPanics: true,
		})

		sup.Add(server.Get())
		sup.Add(server.Get().Supervisor())
		sup.Add(interp.Get())
	*/
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	go server.Get().Serve(ctx)
	go interp.Get().Serve(ctx)
	<-ctx.Done()
	//go server.Get().Supervisor().
	//sup.Serve(ctx)
}
