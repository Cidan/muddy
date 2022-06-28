package construct

import (
	"context"
	"fmt"
	"net"

	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

type Player struct {
	connection net.Conn
	supervisor *suture.Supervisor
	input      chan string //*bufio.Reader
	output     chan *playerv1.Output
	config     chan net.Conn
	Data       *playerv1.Player
	textBuffer string
}

// NewPlayer creates a new player object
func NewPlayer() *Player {
	return &Player{
		supervisor: suture.NewSimple("player"),
		Data:       &playerv1.Player{},
	}
}

// Serve starts this player's server. This server context must
// be the only way in which a player is written to/modified.
func (p *Player) Serve(ctx context.Context) error {
	log.Info().Msg("Starting Player")
	for {
		select {
		// Set the player connection
		case c := <-p.config:
			p.connection = c
		// Process output to the player.
		case output := <-p.output:
			switch output.Type {
			// Direct output of text to the player.
			case playerv1.Output_OUTPUT_TYPE_DIRECT:
				p.write(output.GetText())
				// Buffered write, does not actually write.
			case playerv1.Output_OUTPUT_TYPE_BUFFER:
				p.textBuffer += output.GetText()
				// Flush the buffered output for this user.
			case playerv1.Output_OUTPUT_TYPE_FLUSH:
				p.write(p.textBuffer)
				p.textBuffer = ""
			}
		case input := <-p.input:
			// TODO(lobato): process input via interp
			_ = input
		case <-ctx.Done():
			p.connection.Close()
			return suture.ErrDoNotRestart
			// TODO(lobato): break here
		}
	}
}

func (p *Player) write(text string, args ...interface{}) {
	// TODO(lobato): color and prompt
	if p.connection != nil {
		// TODO(lobato): catch error
		fmt.Fprintf(p.connection, "%s\r\xff\xf9", args...)
	}
}

func (p *Player) writeRaw(text string, args ...interface{}) {
	if p.connection != nil {
		// TODO(lobato): catch error
		fmt.Fprintf(p.connection, text, args...)
	}
}

func (p *Player) Buffer(text string, args ...interface{}) {
	p.output <- &playerv1.Output{
		Type: playerv1.Output_OUTPUT_TYPE_BUFFER,
		Text: fmt.Sprintf(text, args...),
	}
}

func (p *Player) Send(text string, args ...interface{}) {
	p.output <- &playerv1.Output{
		Type: playerv1.Output_OUTPUT_TYPE_DIRECT,
		Text: fmt.Sprintf(text, args...),
	}
}

func (p *Player) Flush() {
	p.output <- &playerv1.Output{
		Type: playerv1.Output_OUTPUT_TYPE_FLUSH,
	}
}
