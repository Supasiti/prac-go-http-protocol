package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Supasiti/prac-go-http-protocol/internal/headers"
)

type parserState int

const (
	parserStateInitialised parserState = iota + 1 // Start from 1 instead of 0
	parserStateHeaders
	parserStateDone
)

const bufferSize = 8
const CRLF = "\r\n"

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine *RequestLine
	Headers     *headers.Headers
	state       parserState
}

func newRequest() *Request {
	return &Request{
		RequestLine: &RequestLine{},
		state:       parserStateInitialised,
		Headers:     headers.NewHeaders(),
	}
}

func (r *Request) parse(p []byte) (int, error) {
	parsedN := 0
	for r.state != parserStateDone {
		n, err := r.parseLine(p[parsedN:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			break
		}
		parsedN += n
	}
	return parsedN, nil
}

func (r *Request) parseLine(p []byte) (int, error) {
	// parse a single line
	switch r.state {
	case parserStateInitialised:
		reqLine, n, err := parseRequestLine(p)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}

		r.RequestLine = reqLine
		r.state = parserStateHeaders
		return n, nil
	case parserStateHeaders:
		n, done, err := r.Headers.Parse(p)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = parserStateDone // finished with scanning headers
		}
		return n, nil
	case parserStateDone:
		return 0, nil
	default:
		return 0, errors.New("Unknown parser state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readIdx := 0
	req := newRequest()

	for req.state != parserStateDone {
		// check if we need to increase the buffer size
		if readIdx >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		// read into buf
		n, err := reader.Read(buf[readIdx:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.state != parserStateDone {
					return nil, fmt.Errorf("Incomplete request in state %d", req.state)
				}
				break
			}
			return nil, err
		}
		readIdx += n

		// parsed the read data
		parsedN, err := req.parse(buf[:readIdx])
		if err != nil {
			return nil, err
		}

		// shift the back in buffer
		copy(buf, buf[parsedN:])
		readIdx -= parsedN
	}

	return req, nil
}

func parseRequestLine(raw []byte) (*RequestLine, int, error) {
	eol := bytes.Index(raw, []byte(CRLF))
	if eol < 0 {
		return nil, 0, nil
	}

	parts := strings.Split(string(raw[:eol]), " ")
	if len(parts) != 3 {
		return nil, 0, errors.New("expect request line to have 3 parts separated by space")
	}

	method := parts[0]
	if err := validateMethod(method); err != nil {
		return nil, 0, err
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 {
		return nil, 0, errors.New("expect http-version to be of format: HTTP-name '/' DIGIT '.' DIGIT")
	}

	httpVersion := httpParts[1]
	if err := validateHttpVersion(httpVersion); err != nil {
		return nil, 0, err
	}

	r := &RequestLine{
		HttpVersion:   httpVersion,
		Method:        method,
		RequestTarget: parts[1],
	}

	return r, eol + 2, nil
}

func validateMethod(method string) error {
	for _, r := range method {
		if r < 'A' || r > 'Z' {
			return errors.New("Request method must use capital alphabets")
		}
	}

	return nil
}

func validateHttpVersion(version string) error {
	if version != "1.1" {
		return errors.New("Invalid HTTP version")
	}
	return nil
}
