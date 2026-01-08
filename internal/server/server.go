package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/Supasiti/prac-go-http-protocol/internal/request"
	"github.com/Supasiti/prac-go-http-protocol/internal/response"
)

const handlerBufSize = 1024

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (h *HandlerError) Write(w io.Writer) {
	msgByte := []byte(h.Message)
	headers := response.GetDefaultHeaders(len(msgByte))

	response.WriteStatusLine(w, h.StatusCode)
	response.WriteHeaders(w, headers)
	w.Write(msgByte)
}

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

	// Parse the request
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Bad request: %s\n", err)
		hErr := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    err.Error(),
		}
		hErr.Write(conn)
		return
	}
	log.Printf("Received %s request on %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget)

	// Calling handler
	buf := bytes.NewBuffer([]byte{})
	if hErr := s.handler(buf, req); hErr != nil {
		log.Printf("%d error caught in handler: %s\n", hErr.StatusCode, hErr.Message)
		hErr.Write(conn)
		return
	}

	// Writing response
	log.Printf("Writing response...\n")
	headers := response.GetDefaultHeaders(buf.Len())

	response.WriteStatusLine(conn, response.StatusOk)
	response.WriteHeaders(conn, headers)
	conn.Write(buf.Bytes())

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
