package response

import (
	"fmt"
	"io"

	"github.com/Supasiti/prac-go-http-protocol/internal/headers"
)

type WriterState int

const (
	WriterStateStatusLine WriterState = iota + 1
	WriterStateHeaders
	WriterStateBody
)

type Writer struct {
	writer io.Writer
	state  WriterState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
		state:  WriterStateStatusLine,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != WriterStateStatusLine {
		return fmt.Errorf("writing response out of order: %d", w.state)
	}
	defer func() { w.state = WriterStateHeaders }()

	line := statusLine(statusCode)
	_, err := w.writer.Write([]byte(line))
	return err
}

func (w *Writer) WriteHeaders(headers *headers.Headers) error {
	if w.state != WriterStateHeaders {
		return fmt.Errorf("writing response out of order: %d", w.state)
	}
	defer func() { w.state = WriterStateBody }()

	for key, value := range headers.All() {
		line := fmt.Sprintf("%s: %s\r\n", key, value)
		if _, err := w.writer.Write([]byte(line)); err != nil {
			return err
		}
	}

	_, err := w.writer.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != WriterStateBody {
		return 0, fmt.Errorf("writing response out of order: %d", w.state)
	}

	return w.writer.Write(p)
}
