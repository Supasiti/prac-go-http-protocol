package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, "localhost:42069", headers.Get("host"))
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with whitespace
	headers = NewHeaders()
	data = []byte("   Host: localhost:42069   \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 29, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing header
	headers = NewHeaders()
	data = []byte("Accept: text/html\r\n")
	_, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "text/html", headers.Get("Accept"))
	assert.False(t, done)
	data = []byte("Accept: application/xhtml+xml\r\n")
	_, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "text/html, application/xhtml+xml", headers.Get("Accept"))
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Valid single header with special characters
	headers = NewHeaders()
	data = []byte("Host1234567890!#$%&'*+-.^_`|~: localhost:42069   \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host1234567890!#$%&'*+-.^_`|~"))
	assert.Equal(t, 51, n)
	assert.False(t, done)

	// Test: Invalid spacing header - header end in whitespace
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid character in header
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

}
