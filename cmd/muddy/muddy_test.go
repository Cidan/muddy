package main

import (
	"context"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	Start(ctx)
}
