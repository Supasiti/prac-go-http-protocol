package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	t := req.RequestLine.RequestTarget
	if strings.HasPrefix(t, "/yourproblem") {
		handle400(w, req)
		return
	}

	if strings.HasPrefix(t, "/myproblem") {
		handle500(w, req)
		return
	}

	if strings.HasPrefix(t, "/httpbin/") {
		handleProxy(w, req)
		return
	}

	handle200(w, req)
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

func handleProxy(w *response.Writer, req *request.Request) {
	path := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	target := fmt.Sprintf("https://httpbin.org/%s", path)

	log.Printf("Proxy target: %s", target)

	// proxy to https://httpbin.org/x
	resp, err := http.Get(target)
	if err != nil {
		log.Printf("Error making GET request: %s", err)
		handle500(w, req)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Request failed with status: %s", resp.Status)
		handle500(w, req)
		return
	}

	// response
	headers := response.GetDefaultHeaders(0)
	headers.Remove("Content-Length")
	headers.Set("Transfer-Encoding", "chunked")

	w.WriteStatusLine(response.StatusOk)
	w.WriteHeaders(headers)

	const maxChunkSize = 1024
	buf := make([]byte, maxChunkSize)

	// write to response body in chunk
	for {
		n, err := resp.Body.Read(buf)
		log.Printf("Read %d bytes", n)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("Reach the end of file")
			} else {
				log.Printf("Error reading response from target: %s", err)
			}
			break
		}

		// write chunk to body
		_, err = w.WriteChunkBody(buf[:n])
		if err != nil {
			log.Printf("Error writing chunk body: %s", err)
			break
		}
	}

	// write closing chunk to body
	_, err = w.WriteChunkBodyDone()
	if err != nil {
		log.Printf("Error writing body done: %s", err)
	}
}
