package interp

import (
	"context"
	"strings"

	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
)

func init() {
	i := Get()
	i.Register(&CommandRef{
		name:   "say",
		interp: playerv1.Player_INTERP_TYPE_PLAYING,
		fn:     DoSay,
	})
}

func DoSay(ctx context.Context, p Player, s ...string) error {
	return p.Send("You say, '%s'", strings.Join(s, " "))
}
