package headers

import (
	"bytes"
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"
)

type Headers struct {
	data map[string]string
}

func NewHeaders() *Headers {
	return &Headers{data: make(map[string]string)}
}

const CRLF = "\r\n"
const keyValueSep = ":"

func (h *Headers) Add(key, value string) {
	if cur := h.Get(key); cur != "" {
		h.Set(key, fmt.Sprintf("%s, %s", cur, value))
	} else {
		h.Set(key, value)
	}
}

func (h *Headers) All() iter.Seq2[string, string] {
	return maps.All(h.data)
}

func (h *Headers) Get(key string) string {
	if value, ok := h.data[strings.ToLower(key)]; ok {
		return value
	} else {
		return ""
	}
}

func (h *Headers) Remove(key string) {
	delete(h.data, strings.ToLower(key))
}

func (h *Headers) Set(key, value string) {
	h.data[strings.ToLower(key)] = value
}

func (h *Headers) Parse(data []byte) (n int, done bool, err error) {
	eol := bytes.Index(data, []byte(CRLF))
	if eol < 0 {
		return 0, false, nil
	}
	if eol == 0 {
		// starting with \r\n
		return eol + 2, true, nil
	}

	clean := bytes.TrimSpace(data[:eol]) // Host: localhost:42069
	sepPos := bytes.Index(clean, []byte(keyValueSep))

	// Key
	key := strings.ToLower(string(clean[:sepPos]))
	if !validKeyTokens(key) {
		return 0, false, fmt.Errorf("Field name coltains invalid character: %s", key)
	}

	// value
	value := string(bytes.TrimSpace(clean[sepPos+1:]))

	h.Add(key, value)
	return eol + 2, false, nil
}

var tokenChars = []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func validKeyTokens(key string) bool {
	for _, k := range key {
		if (k >= 'a' && k <= 'z') || (k >= '0' && k <= '9') {
			continue
		}

		if slices.Contains(tokenChars, k) {
			continue
		}
		return false
	}
	return true
}
