package tapol

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// Returns a new chukedReader
func newChunkedReader(r *bufio.Reader) *chunkedReader {
	return &chunkedReader{
		r: r,
	}
}

// Returns a stablished conn with the server, using or not using TLS
func createConn(schema, url string) (conn net.Conn, err error) {
	if schema == "http" {
		conn, err = connect(url)
	} else {
		conn, err = connectTLS(url)
	}

	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Extract the necessary information from the URL
func parseURL(rawURL string) (schema, host, path string, err error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", "", "", err
	}

	schema = parsedURL.Scheme
	host = parsedURL.Host

	if !strings.Contains(host, ":") {
		switch schema {
		case "http":
			host += ":80"
		default:
			host += ":443"
		}
	}

	path = parsedURL.Path
	if path == "" {
		path = "/"
	}

	return schema, host, path, nil
}

// Add the necessary headers to the request [Header], subscribing if the headers are already set
func buildHeader(header *Header, host string, body *strings.Reader) *Header {
	header.change("Host", host)
	header.change("Connection", "close")

	if body != nil {
		header.change("Content-Length", fmt.Sprint(body.Len()))
	}

	return header
}

// Reads the headers content and parses it to a [Header] map
func parseRespHeader(reader *bufio.Reader) (resp Response, err error) {
	resp = Response{}
	resp.Header = &Header{}

	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	status := strings.Fields(strings.TrimSpace(line))[1:]
	resp.Status = strings.Join(status, " ")
	code, _ := strconv.Atoi(status[0])
	resp.StatusCode = int16(code)

	if resp.StatusCode == 200 {
		resp.Ok = true
	}

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			return Response{}, err
		}

		line = strings.TrimSpace(line)

		if line == "" {
			break
		}

		headerline := strings.SplitN(line, ": ", 2)
		if len(headerline) > 1 {
			resp.Header.Add(headerline[0], headerline[1])
		}
	}

	return resp, nil
}

// The HTTP client entry point and logic holder
//
// Makes the request and return the response from the server
func makeRequest(method string, rawURL string, header *Header, body *strings.Reader) (resp Response, err error) {
	schema, host, path, err := parseURL(rawURL)
	if err != nil {
		return Response{}, err
	}

	conn, err := createConn(schema, host)
	if err != nil {
		return Response{}, err
	}
	// defer conn.Close()

	request := Request{
		Host:   host,
		Header: buildHeader(header, host, body),
		Method: method,
		Path:   path,
		Body:   body,
	}

	// Sending the request
	fmt.Fprintln(conn, request.Build())

	reader := bufio.NewReader(conn)
	resp, err = parseRespHeader(reader)
	if err != nil {
		return Response{}, err
	}

	// Any status of 3xx will lead into other request being made to the Location header
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		location := resp.Header.Get("Location")

		if location == "" {
			return resp, errors.New("the request was redirected, but no valid \"Location\" header were provided")
		}

		return makeRequest(method, location, header, body)
	}

	var bodyreader io.Reader = reader

	contentlength := resp.Header.Get("Content-Length")
	transferencoding := resp.Header.Get("Transfer-Encoding")

	switch transferencoding {
	case "chunked":
		bodyreader = newChunkedReader(reader)
	}

	if contentlength == "" && transferencoding == "" {
		bodyreader = strings.NewReader("")
	} else if contentlength != "" && transferencoding == "" {
		bodylength, err := strconv.Atoi(contentlength)
		if err != nil {
			resp.Body = io.NopCloser(strings.NewReader(""))
			return resp, errors.New("the server didn't provided a valid \"Content-Length\" header")
		}

		bodyreader = io.LimitReader(reader, int64(bodylength))
	}
	resp.Body = io.NopCloser(bodyreader)

	return resp, nil
}

// GET request
func Get(url string, header *Header) (resp Response, err error) {
	return makeRequest("GET", url, header, nil)
}

// POST request
func Post(url string, header *Header, body *strings.Reader) (resp Response, err error) {
	return makeRequest("POST", url, header, body)
}

// PUT request
func Put(url string, header *Header, body *strings.Reader) (resp Response, err error) {
	return makeRequest("PUT", url, header, body)
}

// DELETE request
func Delete(url string, header *Header) (resp Response, err error) {
	return makeRequest("DELETE", url, header, nil)
}

// Create a [net.Conn] using the [tls] package for the
// TLS handshake
// 
// Used in HTTPS connections
func connectTLS(host string) (net.Conn, error) {
	conn, err := tls.Dial("tcp", host, &tls.Config{})
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Create a [net.Conn] using only the TCP handshake
//
// Used in HTTP connections
func connect(host string) (net.Conn, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
