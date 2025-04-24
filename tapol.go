package tapol

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func newChunkedReader(r *bufio.Reader) *chunkedReader {
	return &chunkedReader{
		r: r,
	}
}

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

func buildHeader(header *Header, host string, body *strings.Reader) *Header {
	header.change("Host", host)
	header.change("Connection", "close")

	if body != nil {
		header.change("Content-Length", fmt.Sprint(body.Len()))
	}

	return header
}

func parseRespHeader(reader *bufio.Reader) (resp Response, err error) {
	resp = Response{}
	resp.Header = Header{}

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

	fmt.Fprintln(conn, request.Build())

	reader := bufio.NewReader(conn)
	resp, err = parseRespHeader(reader)
	if err != nil {
		return resp, err
	}

	// A redirection of any status 3xx will lead in another request for now
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		return makeRequest(method, resp.Header.Get("Location"), header, body)
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
			bodyreader = strings.NewReader("")
			resp.Body = io.NopCloser(bodyreader)
			return resp, nil
		}
		bodyreader = io.LimitReader(reader, int64(bodylength))
	}
	resp.Body = io.NopCloser(bodyreader)

	return resp, nil
}

func Get(url string, header *Header) (resp Response, err error) {
	return makeRequest("GET", url, header, nil)
}

func Post(url string, header *Header, body *strings.Reader) (resp Response, err error) {
	return makeRequest("POST", url, header, body)
}

func Put(url string, header *Header, body *strings.Reader) (resp Response, err error) {
	return makeRequest("PUT", url, header, body)
}

func Delete(url string, header *Header) (resp Response, err error) {
	return makeRequest("DELETE", url, header, nil)
}

func connectTLS(host string) (net.Conn, error) {
	conn, err := tls.Dial("tcp", host, &tls.Config{})
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func connect(host string) (net.Conn, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
