package server

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	once   sync.Once
	server *Server
)

type Server struct{}

func Get() *Server {
	once.Do(func() {
		server = &Server{}
	})
	return server
}

func (s *Server) Serve(ctx context.Context) error {
	log.Info().Msg("Starting Server")
	// TODO(lobato): implement server
	for {
		time.Sleep(1 * time.Hour)
	}
}
