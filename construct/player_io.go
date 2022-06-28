package construct

import (
	"context"
	"fmt"
	"net"

	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
	"github.com/thejerf/suture/v4"
)

type PlayerInputOutput struct {
	connection net.Conn
	input      chan string //*bufio.Reader
	output     chan *playerv1.Output
	config     chan net.Conn
	textBuffer string
}

func (p *PlayerInputOutput) Serve(ctx context.Context) error {
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

func (p *PlayerInputOutput) write(text string, args ...interface{}) {
	// TODO(lobato): color and prompt
	if p.connection != nil {
		// TODO(lobato): catch error
		fmt.Fprintf(p.connection, "%s\r\xff\xf9", args...)
	}
}

func (p *PlayerInputOutput) writeRaw(text string, args ...interface{}) {
	if p.connection != nil {
		// TODO(lobato): catch error
		fmt.Fprintf(p.connection, text, args...)
	}
}
