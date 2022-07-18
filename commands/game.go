package commands

import (
	"context"
	"strings"
	"sync"

	"github.com/Cidan/muddy/atlas"
	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
	"github.com/Cidan/muddy/interp"
	"github.com/rs/zerolog/log"
)

type Game struct {
	commands map[string]*interp.Command
	lock     sync.RWMutex
}

func (g *Game) Process(ctx context.Context, player atlas.Player, command string, input ...string) error {
	g.lock.RLock()
	c, ok := g.commands[command]
	g.lock.RUnlock()
	if !ok {
		return interp.ErrNoCommand
	} else {
		return c.Fn(ctx, player, input...)
	}
}

func (g *Game) Register(c *interp.Command) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.commands[c.Name] = c
}

func (g *Game) DoSay(ctx context.Context, player atlas.Player, input ...string) error {
	text := strings.Join(input, " ")
	player.Send("You say '%s'", text)
	// player.SendToRoom("%s says, '%s'", player.Name(), text)
	return nil
}

func init() {
	log.Debug().Msg("registering game interp")
	g := &Game{
		commands: make(map[string]*interp.Command),
	}
	g.Register(&interp.Command{
		Name: "say",
		Fn:   g.DoSay,
	})
	interp.Get().Set(playerv1.Player_INTERP_TYPE_PLAYING, g)
}
