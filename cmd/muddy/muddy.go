package main

import (
	"context"

	"github.com/Cidan/muddy/construct"
	"github.com/Cidan/muddy/server"
	"github.com/thejerf/suture/v4"
)

func main() {
	sup := suture.New("muddy", suture.Spec{
		PassThroughPanics: true,
	})

	sup.Add(server.Get())
	sup.Add(construct.NewPlayer())
	if err := sup.Serve(context.Background()); err != nil {
		panic(err)
	}

}
