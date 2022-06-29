package server

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/Cidan/muddy/construct"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

var (
	once   sync.Once
	server *Server
)

type Server struct {
	supervisor *suture.Supervisor
}

func Get() *Server {
	once.Do(func() {
		server = &Server{
			supervisor: suture.NewSimple("server"),
		}
	})
	return server
}

func (s *Server) Serve(ctx context.Context) error {
	log.Info().Msg("Starting Server")
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 4000))

	if err != nil {
		log.Err(err).Msg("unable to listen to port")
		return suture.ErrTerminateSupervisorTree
	}

	for {
		c, err := l.Accept()
		if err != nil {
			return err
		} else {
			s.handleConnection(c)
			continue
		}
	}
}

func (s *Server) Supervisor() *suture.Supervisor {
	return s.supervisor
}

func (s *Server) handleConnection(c net.Conn) {
	log.Info().
		Str("address", c.RemoteAddr().String()).
		Msg("New remote player connection")
	p := construct.NewPlayer()
	s.supervisor.Add(p)

	p.SetConnection(c)
	p.Send("Hello, by which name would you like to be known?")
}
