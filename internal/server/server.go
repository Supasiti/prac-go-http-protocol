package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/Supasiti/prac-go-http-protocol/internal/request"
	"github.com/Supasiti/prac-go-http-protocol/internal/response"
)

const handlerBufSize = 1024

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	handler  Handler
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

func (s *Server) handle(conn net.Conn) {
	log.Printf("Connection %s has been accepted\n", conn.LocalAddr())
	defer conn.Close()

	res := response.NewWriter(conn)

	// Parse the request
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Bad request: %s\n", err)
		msg := []byte(err.Error())

		res.WriteStatusLine(response.StatusBadRequest)
		res.WriteHeaders(response.GetDefaultHeaders(len(msg)))
		res.WriteBody(msg)
		return
	}
	log.Printf("Received %s request on %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget)

	// Calling handler
	s.handler(res, req)

	log.Printf("Successfully wrote response\n")
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("Fail to create listener: %s", err)
	}

	server := &Server{
		handler:  handler,
		listener: listener,
	}

	go server.listen()

	return server, nil
}
