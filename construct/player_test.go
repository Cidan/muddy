package construct

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func FuzzHealth(f *testing.F) {
	p := NewPlayer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	go p.Serve(ctx)

	for i := 0; i < 10; i++ {
		f.Add(int32(i), int32(i))
	}

	f.Fuzz(func(t *testing.T, current int32, max int32) {
		p.SetHealth(current)
		p.SetMaxHealth(max)
		c, m := p.GetHealth()
		assert.EqualValues(t, current, c)
		assert.EqualValues(t, max, m)
	})
}

func FuzzMana(f *testing.F) {
	p := NewPlayer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	go p.Serve(ctx)

	for i := 0; i < 10; i++ {
		f.Add(int32(i), int32(i))
	}

	f.Fuzz(func(t *testing.T, current int32, max int32) {
		p.SetMana(current)
		p.SetMaxMana(max)
		c, m := p.GetMana()
		assert.EqualValues(t, current, c)
		assert.EqualValues(t, max, m)
	})
}

func FuzzMove(f *testing.F) {
	p := NewPlayer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	go p.Serve(ctx)

	for i := 0; i < 10; i++ {
		f.Add(int32(i), int32(i))
	}

	f.Fuzz(func(t *testing.T, current int32, max int32) {
		p.SetMove(current)
		p.SetMaxMove(max)
		c, m := p.GetMove()
		assert.EqualValues(t, current, c)
		assert.EqualValues(t, max, m)
	})
}
