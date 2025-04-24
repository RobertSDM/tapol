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
	Header map[string]string
)

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

func (h *Header) change(key string, value string) (err error) {
	if value == "" {
		return errors.New("the value cannot be null")
	}

	(*h)[key] = value

	return nil
}

func (h *Header) exists(key string) (exists bool) {
	return (*h)[key] != ""
}

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

type Request struct {
	Body   []byte
	Host   string
	Header *Header
	Path   string
	Method string
}

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

	HTTPFormedRequest += fmt.Sprintf("\r\n%s\r\n", r.Body)

	return HTTPFormedRequest
}

func (r Request) String() (str string) {
	return fmt.Sprintf("Host: %s; Path: %s; Method: %s", r.Host, r.Path, r.Method)
}

type Response struct {
	StatusCode int16
	Status     string
	Header     Header
	Body       io.ReadCloser
	Ok         bool
}

type chunkedReader struct {
	r      *bufio.Reader
	remain int
	done   bool
}

func (c *chunkedReader) Read(p []byte) (n int, err error) {
	if c.done {
		return 0, io.EOF
	}

	if c.remain == 0 {
		line, err := c.r.ReadString('\n')
		if err != nil {
			return 0, err
		}

		line = strings.TrimSpace(line)

		length, err := strconv.ParseInt(line, 16, 64)
		if err != nil {
			return 0, err
		}

		if length == 0 {
			c.done = true
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
		_, err := c.r.ReadString('\n')
		if err != nil {
			return 0, err
		}
	}

	return n, nil
}
