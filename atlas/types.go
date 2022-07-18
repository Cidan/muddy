package atlas

import (
	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
	roomv1 "github.com/Cidan/muddy/gen/proto/go/room/v1"
)

type Player interface {
	Send(string, ...interface{}) error
	LoginState() string
	SetLoginState(string)
	SetInterp(playerv1.Player_InterpType)
	Interp() playerv1.Player_InterpType
	SetName(string)
	Name() string
	SetPassword(string)
	CheckPassword(string) bool
	Password() string
}

type Room interface {
	Coordinates() *roomv1.Room_Coordinates
}
