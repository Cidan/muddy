package construct

import (
	"context"
	"testing"
	"time"
)

func TestServe(t *testing.T) {
	p := NewPlayer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

	go p.Serve(ctx)

	p.SetName("doot")
	p.Health(10, 100)
	p.Mana(10, 100)
	p.Move(10, 100)
	time.Sleep(time.Second * 2)
	cancel()
}
