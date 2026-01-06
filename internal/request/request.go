package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine *RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(b), "\r\n")
	if len(lines) == 0 {
		return nil, err
	}

	// the first line is  request line
	reqLine, err := parseRequestLine(lines[0])
	if err != nil {
		return nil, err
	}

	return &Request{RequestLine: reqLine}, nil
}

func parseRequestLine(raw string) (*RequestLine, error) {
	parts := strings.Split(raw, " ")

	if len(parts) != 3 {
		return nil, errors.New("expect request line to have 3 parts separated by space")
	}

	if err := validateMethod(parts[0]); err != nil {
		return nil, err
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 {
		return nil, errors.New("expect http-version to be of format: HTTP-name '/' DIGIT '.' DIGIT")
	}

	if err := validateHttpVersion(httpParts[1]); err != nil {
		return nil, err
	}

	return &RequestLine{
		HttpVersion:   httpParts[1],
		Method:        parts[0],
		RequestTarget: parts[1],
	}, nil
}

func validateMethod(raw string) error {
	for _, r := range raw {
		if r < 'A' || r > 'Z' {
			return errors.New("Request method must use capital alphabets")
		}
	}

	return nil
}

func validateHttpVersion(raw string) error {
	if raw != "1.1" {
		return errors.New("Invalid HTTP version")
	}
	return nil
}
