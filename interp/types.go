package interp

import (
	"context"

	"github.com/Cidan/muddy/atlas"
)

type commandCallback func(context.Context, atlas.Player, ...string) error

type Command struct {
	Name  string
	Alias []string
	Fn    commandCallback
}

type Input struct {
	Player atlas.Player
	Text   string
}

type Interp interface {
	Process(context.Context, atlas.Player, string, ...string) error
	Register(*Command)
}
