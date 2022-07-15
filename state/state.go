package state

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type EventCallback func(context.Context, string) error

type Event struct {
	Name string
	Fn   EventCallback
}

// State machine for handling work flows.
type State struct {
	current string
	states  map[string]*Event
	mutex   sync.RWMutex
}

func New(initial string) *State {
	return &State{
		states:  make(map[string]*Event),
		current: initial,
		mutex:   sync.RWMutex{},
	}
}

func (s *State) Add(e *Event) *State {
	s.states[e.Name] = e
	return s
}

// SetState will set the state of this machine. You must ensure
// you only call SetState from a Process function to avoid a race.
func (s *State) SetState(name string) error {
	if s.states[name] == nil {
		return fmt.Errorf("no such state defined, %s", name)
	}
	s.current = name
	return nil
}

func (s *State) Process(ctx context.Context, text string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.states[s.current] == nil {
		return errors.New("invalid state set")
	}
	return s.states[s.current].Fn(ctx, text)
}
