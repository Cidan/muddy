package atlas

import (
	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
	roomv1 "github.com/Cidan/muddy/gen/proto/go/room/v1"
)

type Player interface {
	Uuid() string
	Buffer(string, ...interface{})
	Flush() error
	Send(string, ...interface{}) error
	SendToRoom(string, ...interface{})
	LoginState() string
	SetLoginState(string)
	SetInterp(playerv1.Player_InterpType)
	Interp() playerv1.Player_InterpType
	SetName(string)
	Name() string
	SetPassword(string)
	CheckPassword(string) bool
	Password() string
	ToRoom(Room)
	Room() Room
	Do(string)
}

type Room interface {
	Coordinates() *roomv1.Room_Coordinates
	AddPlayer(string, Player)
	RemovePlayer(string)
	Send(string, Player, ...interface{})
	Name() string
}
