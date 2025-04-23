package tapol

import (
	"errors"
	"fmt"
	"io"
)

type (
	HTTPHeader string
	Header     map[HTTPHeader]string
)

var (
	Host          HTTPHeader = "Host"
	Connection    HTTPHeader = "Connection"
	Accept        HTTPHeader = "Accept"
	ContentLength HTTPHeader = "Content-Length"
	ContentType   HTTPHeader = "Content-Type"
	Location      HTTPHeader = "Location"
)

func (h *Header) Add(key HTTPHeader, value string) (err error) {
	if value == "" {
		return errors.New("the value cannot be null")
	}

	if h.exists(key) {
		return errors.New("the key already exists")
	}

	(*h)[key] = value

	return nil
}

func (h *Header) change(key HTTPHeader, value string) (err error) {
	if value == "" {
		return errors.New("the value cannot be null")
	}

	(*h)[key] = value

	return nil
}

func (h *Header) exists(key HTTPHeader) (exists bool) {
	return (*h)[key] != ""
}

func (h *Header) Get(key HTTPHeader) string {
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
