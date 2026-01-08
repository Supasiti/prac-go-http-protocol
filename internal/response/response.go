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
	WriterStateTrailers
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

func (w *Writer) WriteChunkBody(p []byte) (int, error) {
	if w.state != WriterStateBody {
		return 0, fmt.Errorf("writing response out of order: %d", w.state)
	}

	nTotal := 0
	n, err := fmt.Fprintf(w.writer, "%X\r\n", len(p))
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	body := fmt.Appendf(p, "\r\n")
	n, err = w.writer.Write(body)
	nTotal += n

	return nTotal, err
}

func (w *Writer) WriteChunkBodyDone() (int, error) {
	if w.state != WriterStateBody {
		return 0, fmt.Errorf("writing response out of order: %d", w.state)
	}
	defer func() { w.state = WriterStateTrailers }()

	body := []byte("0\r\n")
	return w.writer.Write(body)
}

func (w *Writer) WriteTrailers(h *headers.Headers) error {
	if w.state != WriterStateTrailers {
		return fmt.Errorf("writing response out of order: %d", w.state)
	}

	if h != nil {
		for key, value := range h.All() {
			line := fmt.Sprintf("%s: %s\r\n", key, value)

			fmt.Println(line)
			_, err := w.writer.Write([]byte(line))
			if err != nil {
				return err
			}
		}
	}

	_, err := w.writer.Write([]byte("\r\n"))
	return err
}
