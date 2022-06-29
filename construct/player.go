package construct

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

type Player struct {
	connection net.Conn
	input      chan string //*bufio.Reader
	output     chan *playerv1.Output
	config     chan net.Conn
	update     chan *playerv1.Update
	data       *playerv1.Player
	lock       sync.RWMutex
	ticker     *time.Ticker
	textBuffer string
}

// NewPlayer creates a new player object
func NewPlayer() *Player {
	return &Player{
		input:  make(chan string),
		output: make(chan *playerv1.Output),
		config: make(chan net.Conn),
		update: make(chan *playerv1.Update),
		lock:   sync.RWMutex{},
		ticker: time.NewTicker(time.Second),
		data:   &playerv1.Player{},
	}
}

// Serve starts this player's server. This server context must
// be the only way in which a player is written to/modified.
func (p *Player) Serve(ctx context.Context) error {
	log.Info().Msg("Starting Player")
	for {
		select {
		case <-p.ticker.C:
			// TODO(lobato): One second ticker here.
			log.Debug().Msg("tick")
		case u := <-p.update:
			log.Debug().Msg("updating player")
			p.handlePlayerUpdate(ctx, u)
		// Set the player connection
		case c := <-p.config:
			log.Debug().Msg("setting connection on player")
			p.connection = c
			s := bufio.NewScanner(c)
			// Wrap the connection reader in a channel and launch it
			// in a goroutine. This goroutine will exit automatically
			// in the event of the network socket closing.
			go func(s *bufio.Scanner) {
				for {
					if !s.Scan() {
						break
					}
					p.input <- s.Text()
				}
				log.Debug().Msg("player bufio scanner exiting")
			}(s)

		// Process output to the player.
		case output := <-p.output:
			switch output.Type {
			// Direct output of text to the player.
			case playerv1.Output_OUTPUT_TYPE_DIRECT:
				if err := p.write(output.GetText()); err != nil {
					log.Error().Msg("error sending to user")
				}
				// Buffered write, does not actually write.
			case playerv1.Output_OUTPUT_TYPE_BUFFER:
				p.textBuffer += output.GetText()
				// Flush the buffered output for this user.
			case playerv1.Output_OUTPUT_TYPE_FLUSH:
				p.write(p.textBuffer)
				p.textBuffer = ""
			}
		case input := <-p.input:
			p.lock.RLock()
			log.Debug().Str("name", p.data.Name).Str("input", input).Msg("Got input from player")
			p.lock.RUnlock()
			// TODO(lobato): process input via interp
			_ = input
		case <-ctx.Done():
			// TODO(lobato): cleanup, server is shutting down.
			p.connection.Close()
			return suture.ErrDoNotRestart
		}
	}
}

// HandlePlayerUpdate handles an update call to this player
func (p *Player) handlePlayerUpdate(ctx context.Context, u *playerv1.Update) {
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

func (p *Player) write(text string) error {
	// TODO(lobato): color and prompt
	if p.connection != nil {
		// TODO(lobato): catch error
		_, err := fmt.Fprintf(p.connection, "%s\r\xff\xf9", text)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Player) writeRaw(text string, args ...interface{}) {
	if p.connection != nil {
		// TODO(lobato): catch error
		fmt.Fprintf(p.connection, text, args...)
	}
}

// Send will send the given text to the player.
func (p *Player) Send(text string, args ...interface{}) {
	p.output <- &playerv1.Output{
		Type: playerv1.Output_OUTPUT_TYPE_DIRECT,
		Text: fmt.Sprintf(text, args...),
	}
}

// Buffer will buffer the given text for a player, and send
// that text upon calling Flush(). Multiple calls to this will
// append new text to the buffer.
func (p *Player) Buffer(text string, args ...interface{}) {
	p.output <- &playerv1.Output{
		Type: playerv1.Output_OUTPUT_TYPE_BUFFER,
		Text: fmt.Sprintf(text, args...),
	}
}

// Flush will send the buffered text that was called from Buffer()
// to the player. The buffer will be cleared after the text has been sent.
func (p *Player) Flush() {
	p.output <- &playerv1.Output{
		Type: playerv1.Output_OUTPUT_TYPE_FLUSH,
	}
}

// SetConnection sets the user's connection to the given net.Conn.
func (p *Player) SetConnection(c net.Conn) {
	p.config <- c
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
