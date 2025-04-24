package tapol

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type (
	// String map representing the header key and value
	Header map[string]string
)

// Add a new key and value to the [Header] map if the key does't already exists
func (h *Header) Add(key string, value string) (err error) {
	if value == "" {
		return errors.New("the value cannot be null")
	}

	if h.exists(key) {
		return errors.New("the key already exists")
	}

	(*h)[key] = value

	return nil
}

// Add a new key to the [Header] map, but it changes the value if the key  exists
func (h *Header) change(key string, value string) (err error) {
	if value == "" {
		return errors.New("the value cannot be null")
	}

	(*h)[key] = value

	return nil
}

// Returns a bool if the key exists in the [Header] map
func (h *Header) exists(key string) (exists bool) {
	return (*h)[key] != ""
}

// Returns the value of a [Header] map key
func (h *Header) Get(key string) string {
	return (*h)[key]
}

func (h Header) String() string {
	str := ""

	for k, v := range h {
		str += fmt.Sprintf("%s: %s\n", k, v)
	}

	return str
}

var HTTPversion = "1.1"

// Represents the request that is going to be made to the server
type Request struct {
	// For now the body is just a [strings.Reader]
	// 
	// Represents the body to be sent,
	// there are no streaming on the request body
	Body *strings.Reader

	// The url host with the post number
	Host string // localhost:3000

	// A representation of the header key and value in a map
	Header *Header

	// The path you are making the request
	Path string

	// The method of you request
	Method string
}

// Reads the [Request] values and build the string message to send to the server
func (r *Request) Build() (request string) {
	HTTPFormedRequest := ""

	// Adding the HTTP header
	HTTPFormedRequest += fmt.Sprintf("%s %s HTTP/%s\r\n", r.Method, r.Path, HTTPversion)

	// Adding the request headers
	for k, v := range *r.Header {
		HTTPFormedRequest += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	if r.Body == nil {
		return HTTPFormedRequest
	}

	b, _ := io.ReadAll(r.Body)

	HTTPFormedRequest += fmt.Sprintf("\r\n%s\r\n", string(b))

	return HTTPFormedRequest
}

func (r Request) String() (str string) {
	return fmt.Sprintf("Host: %s; Path: %s; Method: %s", r.Host, r.Path, r.Method)
}

// Represents the response sent by the server
type Response struct {
	// The HTTP status code
	StatusCode int16 // 200

	// The HTTP status code with the status message
	Status string // "201 Created"

	// A representation of the header key and value in a map
	Header *Header

	// The body of the request
	//
	// This representatoin allows the body to be streamed from the server if
	// needed. Its you responsability to close the reader
	//
	// The returned body will never be nil, even if the the body is not present
	// in the response. If it is not present in the response it with be just a
	// [strings.Reader] from a empty string
	Body io.ReadCloser

	// Ok will be set to true if the response status code is 200, otherwise it
	// will be false
	Ok bool
}

// The the [io.Reader] to read the body content when the "Transfer-Encoding"="chunked" header is present
type chunkedReader struct {
	// A [bufio.Reader] to represent the reader we are wrapping up
	r      *bufio.Reader

	// How much bytes we need to read from the chunk-data
	remain int

	// Flag set to true when the "0\r\n" chunk-data size is received,
	// indicating the reading end
	done   bool
}

func (c *chunkedReader) Read(p []byte) (n int, err error) {
	if c.done {
		return 0, io.EOF
	}

	// Read the next chunk-data length, only if the remain is 0
	// other wise will keep reading the content of the current chunk-data
	if c.remain == 0 {
		line, err := c.r.ReadString('\n')
		if err != nil {
			return 0, err
		}

		line = strings.TrimSpace(line)

		// HEX to decimal
		length, err := strconv.ParseInt(line, 16, 64)
		if err != nil {
			return 0, err
		}

		if length == 0 {
			c.done = true

			// Reading the last chunk-data
			c.r.ReadString('\n')
			return 0, io.EOF
		}

		c.remain = int(length)
	}

	toRead := min(len(p), c.remain)

	n, err = c.r.Read(p[:toRead])
	if err != nil {
		return 0, err
	}

	c.remain -= n

	if c.remain == 0 {
		// Go to the last part of the chunk-data
		// cleaning the way for the next chunk-data length reading
		_, err := c.r.ReadString('\n')
		if err != nil {
			return 0, err
		}
	}

	return n, nil
}
