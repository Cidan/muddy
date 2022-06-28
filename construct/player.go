package construct

import (
	"context"
	"fmt"
	"net"
	"sync"

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
	update     chan *playerv1.Update
	data       *playerv1.Player
	lock       sync.RWMutex
	textBuffer string
}

// NewPlayer creates a new player object
func NewPlayer() *Player {
	return &Player{
		supervisor: suture.NewSimple("player"),
		input:      make(chan string),
		output:     make(chan *playerv1.Output),
		config:     make(chan net.Conn),
		update:     make(chan *playerv1.Update),
		lock:       sync.RWMutex{},
		data:       &playerv1.Player{},
	}
}

// Serve starts this player's server. This server context must
// be the only way in which a player is written to/modified.
func (p *Player) Serve(ctx context.Context) error {
	log.Info().Msg("Starting Player")
	for {
		select {
		case u := <-p.update:
			p.HandlePlayerUpdate(ctx, u)
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

// HandlePlayerUpdate handles an update call to this player
func (p *Player) HandlePlayerUpdate(ctx context.Context, u *playerv1.Update) {
	p.lock.Lock()
	defer p.lock.Unlock()
	switch u.Type {
	case playerv1.Update_UPDATE_TYPE_NAME:
		p.data.Name = u.GetName()
	case playerv1.Update_UPDATE_TYPE_HMV:
		p.data.Health += u.GetHealth()
		p.data.MaxHealth += u.GetMaxHealth()
		p.data.Mana += u.GetMana()
		p.data.MaxMana += u.GetMaxMana()
		p.data.Move += u.GetMove()
		p.data.MaxMove += u.GetMaxMove()
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

// Health changes the player's health by the given amounts.
// Current will change the player's current health (i.e. take damage, heal)
// whereas max will change the player's max health permanently.
func (p *Player) Health(current, max int32) {
	p.update <- &playerv1.Update{
		Health:    current,
		MaxHealth: max,
	}
}

// Mana changes the player's mana by the given amounts.
// Current will change the player's current mana
// whereas max will change the player's max mana permanently.
func (p *Player) Mana(current, max int32) {
	p.update <- &playerv1.Update{
		Mana:    current,
		MaxMana: max,
	}
}

// Move changes the player's move by the given amounts.
// Current will change the player's current move
// whereas max will change the player's max move permanently.
func (p *Player) Move(current, max int32) {
	p.update <- &playerv1.Update{
		Move:    current,
		MaxMove: max,
	}
}

// GetHealth returns the player's current and maximum health.
func (p *Player) GetHealth() (int32, int32) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.data.Health, p.data.MaxHealth
}

// GetMana returns the player's current and maximum mana.
func (p *Player) GetMana() (int32, int32) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.data.Mana, p.data.MaxMana
}

// GetMove returns the player's current and maximum move.
func (p *Player) GetMove() (int32, int32) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.data.Move, p.data.MaxMove
}
