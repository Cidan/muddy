package construct

import (
	"context"
	"fmt"
	"net"

	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
	"github.com/thejerf/suture/v4"
)

type Player struct {
	supervisor *suture.Supervisor
	io         *PlayerInputOutput
	Data       *playerv1.Player
}

// NewPlayer creates a new player object
func NewPlayer() *Player {
	return &Player{
		supervisor: suture.NewSimple("player"),
		Data:       &playerv1.Player{},
	}
}

func (p *Player) Serve(ctx context.Context) error {
	io := &PlayerInputOutput{
		input:  make(chan string),
		output: make(chan *playerv1.Output),
		config: make(chan net.Conn),
	}
	p.supervisor.Add(io)

	return p.supervisor.Serve(ctx)
}

func (p *Player) Buffer(text string, args ...interface{}) {
	p.io.output <- &playerv1.Output{
		Type: playerv1.Output_OUTPUT_TYPE_BUFFER,
		Text: fmt.Sprintf(text, args...),
	}
}

func (p *Player) Send(text string, args ...interface{}) {
	p.io.output <- &playerv1.Output{
		Type: playerv1.Output_OUTPUT_TYPE_DIRECT,
		Text: fmt.Sprintf(text, args...),
	}
}

func (p *Player) Flush() {
	p.io.output <- &playerv1.Output{
		Type: playerv1.Output_OUTPUT_TYPE_FLUSH,
	}
}
