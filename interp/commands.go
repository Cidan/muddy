package interp

import (
	"context"
	"strings"
)

func init() {
	i := Get()
	i.Register(&CommandRef{
		name: "say",
		fn:   DoSay,
	})
}

func DoSay(ctx context.Context, p Player, s ...string) error {
	return p.Send("You say, '%s'", strings.Join(s, " "))
}
