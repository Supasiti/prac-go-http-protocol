package main

import (
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

func handler(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		handle400(w, req)
		return
	case "/myproblem":
		handle500(w, req)
		return
	default:
		handle200(w, req)
	}
}

func handle200(w *response.Writer, _ *request.Request) {
	bodyBytes := []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`)
	headers := response.GetDefaultHeaders(len(bodyBytes))
	headers.Set("Content-Type", "text/html")

	w.WriteStatusLine(response.StatusOk)
	w.WriteHeaders(headers)
	w.WriteBody(bodyBytes)
}

func handle400(w *response.Writer, _ *request.Request) {
	bodyBytes := []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`)
	headers := response.GetDefaultHeaders(len(bodyBytes))
	headers.Set("Content-Type", "text/html")

	w.WriteStatusLine(response.StatusBadRequest)
	w.WriteHeaders(headers)
	w.WriteBody(bodyBytes)
}

func handle500(w *response.Writer, req *request.Request) {
	bodyBytes := []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`)
	headers := response.GetDefaultHeaders(len(bodyBytes))
	headers.Set("Content-Type", "text/html")

	w.WriteStatusLine(response.StatusInternalServerError)
	w.WriteHeaders(headers)
	w.WriteBody(bodyBytes)
}
