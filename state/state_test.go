package state

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndToEnd(t *testing.T) {
	s := New("START")

	s.Add(&Event{
		Name: "START",
		Fn: func(ctx context.Context, str string) error {
			assert.Equal(t, "help", str)
			return nil
		},
	}).Add(&Event{
		Name: "SET_OUTSIDE",
		Fn: func(ctx context.Context, str string) error {
			assert.Equal(t, "set outside", str)
			s.SetState("SET_INSIDE")
			return nil
		},
	}).Add(&Event{
		Name: "SET_INSIDE",
		Fn: func(ctx context.Context, str string) error {
			assert.Equal(t, "set inside", str)
			return nil
		},
	})

	ctx := context.Background()
	assert.Nil(t, s.Process(ctx, "help"))
	s.SetState("SET_OUTSIDE")
	assert.Nil(t, s.Process(ctx, "set outside"))
	assert.Nil(t, s.Process(ctx, "set inside"))
	assert.Error(t, s.SetState("invalid_state"))
}
