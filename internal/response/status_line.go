package response

import (
	"fmt"
	"io"
)

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func statusLine(statusCode StatusCode) string {
	reasonPharse := ""
	switch statusCode {
	case StatusOk:
		reasonPharse = "OK"
	case StatusBadRequest:
		reasonPharse = "Bad Request"
	case StatusInternalServerError:
		reasonPharse = "Internal Server Error"
	}

	return fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPharse)
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	line := statusLine(statusCode)
	_, err := w.Write([]byte(line))
	return err
}
