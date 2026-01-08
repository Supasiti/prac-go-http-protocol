package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Supasiti/prac-go-http-protocol/internal/request"
	"github.com/Supasiti/prac-go-http-protocol/internal/response"
	"github.com/Supasiti/prac-go-http-protocol/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    "Your problem is not my problem\n",
		}
	case "/myproblem":
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    "Woopsie, my bad\n",
		}
	default:
		if _, err := w.Write([]byte("All good, frfr\n")); err != nil {
			return &server.HandlerError{
				StatusCode: response.StatusInternalServerError,
				Message:    "Internal server error\n",
			}
		}
		return nil
	}
}
