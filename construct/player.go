package construct

import playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"

type Player struct {
	data *playerv1.Player
}

func NewPlayer() *Player {
	return &Player{
		data: &playerv1.Player{},
	}
}
