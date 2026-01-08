package response

import (
	"fmt"
	"io"

	"github.com/Supasiti/prac-go-http-protocol/internal/headers"
)

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, headers *headers.Headers) error {
	for key, value := range headers.All() {
		line := fmt.Sprintf("%s: %s\r\n", key, value)
		if _, err := w.Write([]byte(line)); err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\r\n"))
	return err
}
