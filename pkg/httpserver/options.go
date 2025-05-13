package httpserver

import "net"

type Option func(s *Server)

func Port(port string) Option {
	return func(s *Server) {
		s.App.Addr = net.JoinHostPort("", port)
	}
}
