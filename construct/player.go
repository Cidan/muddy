package construct

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	playerv1 "github.com/Cidan/muddy/gen/proto/go/player/v1"
	"github.com/Cidan/muddy/interp"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

type Player struct {
	connection net.Conn
	input      chan string //*bufio.Reader
	disconnect chan bool
	data       *playerv1.Player
	lock       sync.RWMutex
	ticker     *time.Ticker
	interp     playerv1.Player_InterpType
	loginState string
	textBuffer string
}

// NewPlayer creates a new player object
func NewPlayer() *Player {
	return &Player{
		input:      make(chan string),
		disconnect: make(chan bool),
		lock:       sync.RWMutex{},
		ticker:     time.NewTicker(time.Second),
		interp:     playerv1.Player_INTERP_TYPE_LOGIN,
		loginState: "ASK_NAME",
		data:       &playerv1.Player{},
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
		case input := <-p.input:
			// TODO(lobato): sanitize player input
			p.lock.RLock()
			log.Debug().Str("name", p.data.Name).Str("input", input).Msg("Got input from player")
			p.lock.RUnlock()
			interp.Get().Do(&interp.Input{
				Player: p,
				Text:   input,
			})
			// Slow the user down so they can't spam the game.
			// TODO(lobato): make this adjustable as part of a "haste" mechanism
			time.Sleep(time.Millisecond * 50)
		case <-p.disconnect:
			if p.interp == playerv1.Player_INTERP_TYPE_LOGIN {
				p.interp = playerv1.Player_INTERP_TYPE_UNSPECIFIED
				log.Debug().Msg("player has disconnected during login, cleaning up")
				p.cleanup()
				return suture.ErrDoNotRestart
			}
		case <-ctx.Done():
			p.cleanup()
			// TODO(lobato): cleanup, server is shutting down.
			return suture.ErrDoNotRestart
		}
	}
}

func (p *Player) cleanup() {
	p.ticker.Stop()
	close(p.input)
	if p.connection != nil {
		p.connection.Close()
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
func (p *Player) Send(text string, args ...interface{}) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.write(fmt.Sprintf(text, args...))
}

// Buffer will buffer the given text for a player, and send
// that text upon calling Flush(). Multiple calls to this will
// append new text to the buffer.
func (p *Player) Buffer(text string, args ...interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.textBuffer += fmt.Sprintf(text, args...)
}

// Flush will send the buffered text that was called from Buffer()
// to the player. The buffer will be cleared after the text has been sent.
func (p *Player) Flush() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	err := p.write(p.textBuffer)
	p.textBuffer = ""
	return err
}

// SetConnection sets the user's connection to the given net.Conn.
func (p *Player) SetConnection(c net.Conn) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.connection != nil {
		return fmt.Errorf("you must disconnect the player before you set their connection")
	}

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
		p.disconnect <- true
		log.Debug().Msg("player bufio scanner exiting")
	}(s)
	return nil
}

func (p *Player) Disconnect() {
	p.disconnect <- true
}

// SetName sets the name of this player.
func (p *Player) SetName(name string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.data.Name = name
	// TODO(lobato): validate name
}

// Health changes the player's health by the given amounts.
// Current will change the player's current health (i.e. take damage, heal)
// whereas max will change the player's max health permanently.
func (p *Player) AddHealth(current, max int32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.data.Health += current
	p.data.MaxHealth += max
}

// Mana changes the player's mana by the given amounts.
// Current will change the player's current mana
// whereas max will change the player's max mana permanently.
func (p *Player) AddMana(current, max int32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.data.Mana += current
	p.data.MaxMana += max
}

// Move changes the player's move by the given amounts.
// Current will change the player's current move
// whereas max will change the player's max move permanently.
func (p *Player) AddMove(current, max int32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.data.Move += current
	p.data.MaxMove += max
}

// Health changes the player's health by the given amounts.
// Current will change the player's current health (i.e. take damage, heal)
// whereas max will change the player's max health permanently.
func (p *Player) SetHealth(current int32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.data.Health = current
}

// Mana changes the player's mana by the given amounts.
// Current will change the player's current mana
// whereas max will change the player's max mana permanently.
func (p *Player) SetMana(current int32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.data.Mana = current
}

// Move changes the player's move by the given amounts.
// Current will change the player's current move
// whereas max will change the player's max move permanently.
func (p *Player) SetMove(current int32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.data.Move = current
}

// Health changes the player's health by the given amounts.
// Current will change the player's current health (i.e. take damage, heal)
// whereas max will change the player's max health permanently.
func (p *Player) SetMaxHealth(max int32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.data.MaxHealth = max
}

// Mana changes the player's mana by the given amounts.
// Current will change the player's current mana
// whereas max will change the player's max mana permanently.
func (p *Player) SetMaxMana(max int32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.data.MaxMana = max
}

// Move changes the player's move by the given amounts.
// Current will change the player's current move
// whereas max will change the player's max move permanently.
func (p *Player) SetMaxMove(max int32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.data.MaxMove = max
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

func (p *Player) Interp() playerv1.Player_InterpType {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.interp
}

func (p *Player) SetInterp(interp playerv1.Player_InterpType) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.interp = interp
}

func (p *Player) LoginState() string {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.loginState
}

func (p *Player) SetLoginState(state string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.loginState = state
}

func (p *Player) Name() string {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.data.Name
}
