package atlas

import playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"

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
