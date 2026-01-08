package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("unable to accept the connection: %s\n", err)
			continue
		}

		go s.handle(conn)
	}
}

const defaultResponse = `HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 13

Hello World!
`

func (s *Server) handle(conn net.Conn) {
	log.Printf("Connection %s has been accepted\n", conn.LocalAddr())
	defer conn.Close()

	conn.Write([]byte(defaultResponse))
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("Fail to create listener: %s", err)
	}

	server := &Server{
		listener: listener,
	}

	go server.listen()

	return server, nil
}
