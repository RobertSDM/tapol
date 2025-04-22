package tapol

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func createConn(url *url.URL) (conn net.Conn, err error) {
	if url.Scheme == "http" {
		conn, err = connect((*url).Host)
	} else {
		conn, err = connectTLS((*url).Host)
	}

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func buildHeader(header Header, host string, body []byte) Header {
	header.Change(Host, host)
	header.Change(Connection, "close")

	if body != nil {
		bodylength := len(body)
		header.Change(ContentLength, fmt.Sprint(bodylength))
	}

	return header
}

func Get(url *url.URL, header Header) (resp Response, err error) {
	conn, err := createConn(url)
	if err != nil {
		return Response{}, err
	}
	defer conn.Close()

	request := Request{}
	request.Host = url.Host
	request.Header = buildHeader(header, url.Host, nil)
	request.Method = "GET"

	request.Path = url.Path
	if url.Path == "" {
		request.Path = "/"
	}

	HTTPRequest := request.Build()
	fmt.Fprintln(conn, HTTPRequest)

	resp = Response{}
	scanner := bufio.NewScanner(conn)
	responseHeader := false
	respHeader := Header{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		if !responseHeader {
			status := strings.Split(line, " ")[1:]
			resp.Status = strings.Join(status, " ")
			resp.StatusCode = status[0]

			if resp.StatusCode == "200" {
				resp.Ok = true
			}

			responseHeader = true
		} else {
			headerLine := strings.Split(line, ": ")
			respHeader.Add(HTTPHeader(headerLine[0]), headerLine[1])
		}
	}

	resp.Header = respHeader

	contentLength, err := strconv.Atoi(resp.Header.Get(ContentLength))
	if err == nil && contentLength > 0 {
		respBody := ""
		for scanner.Scan() {
			respBody += scanner.Text()
		}

		resp.Body = respBody
	}

	return resp, nil
}

func Post(url *url.URL, header Header, body []byte) (resp Response, err error) {
	conn, err := createConn(url)
	if err != nil {
		return Response{}, err
	}
	defer conn.Close()

	request := Request{}
	request.Host = url.Host
	request.Header = buildHeader(header, url.Host, body)
	request.Body = body
	request.Method = "POST"

	request.Path = url.Path
	if url.Path == "" {
		request.Path = "/"
	}

	HTTPRequest := request.Build()
	fmt.Fprintln(conn, HTTPRequest)

	resp = Response{}
	scanner := bufio.NewScanner(conn)
	responseHeader := false
	respHeader := Header{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		if !responseHeader {
			status := strings.Split(line, " ")[1:]
			resp.Status = strings.Join(status, " ")
			resp.StatusCode = status[0]

			if resp.StatusCode == "200" {
				resp.Ok = true
			}

			responseHeader = true
		} else {
			headerLine := strings.Split(line, ": ")
			respHeader.Add(HTTPHeader(headerLine[0]), headerLine[1])
		}
	}

	resp.Header = respHeader

	contentLength, err := strconv.Atoi(resp.Header.Get(ContentLength))
	if err == nil && contentLength > 0 {
		respBody := ""
		for scanner.Scan() {
			respBody += scanner.Text()
		}

		resp.Body = respBody
	}

	return resp, nil
}

func Put(url *url.URL, header Header, body []byte) (resp Response, err error) {
	conn, err := createConn(url)
	if err != nil {
		return Response{}, err
	}
	defer conn.Close()

	request := Request{}
	request.Host = url.Host
	request.Header = buildHeader(header, url.Host, body)
	request.Body = body
	request.Method = "PUT"

	request.Path = url.Path
	if url.Path == "" {
		request.Path = "/"
	}

	HTTPRequest := request.Build()
	fmt.Fprintln(conn, HTTPRequest)

	resp = Response{}
	scanner := bufio.NewScanner(conn)
	responseHeader := false
	respHeader := Header{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		if !responseHeader {
			status := strings.Split(line, " ")[1:]
			resp.Status = strings.Join(status, " ")
			resp.StatusCode = status[0]

			if resp.StatusCode == "200" {
				resp.Ok = true
			}

			responseHeader = true
		} else {
			headerLine := strings.Split(line, ": ")
			respHeader.Add(HTTPHeader(headerLine[0]), headerLine[1])
		}
	}

	resp.Header = respHeader

	contentLength, err := strconv.Atoi(resp.Header.Get(ContentLength))
	if err == nil && contentLength > 0 {
		respBody := ""
		for scanner.Scan() {
			respBody += scanner.Text()
		}

		resp.Body = respBody
	}

	return resp, nil
}

func Delete(url *url.URL, header Header) (resp Response, err error) {
	conn, err := createConn(url)
	if err != nil {
		return Response{}, err
	}
	defer conn.Close()

	request := Request{}
	request.Host = url.Host
	request.Header = buildHeader(header, url.Host, nil)
	request.Method = "DELETE"

	request.Path = url.Path
	if url.Path == "" {
		request.Path = "/"
	}

	HTTPRequest := request.Build()
	fmt.Fprintln(conn, HTTPRequest)

	resp = Response{}
	scanner := bufio.NewScanner(conn)
	responseHeader := false
	respHeader := Header{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		if !responseHeader {
			status := strings.Split(line, " ")[1:]
			resp.Status = strings.Join(status, " ")
			resp.StatusCode = status[0]

			if resp.StatusCode == "200" {
				resp.Ok = true
			}

			responseHeader = true
		} else {
			headerLine := strings.Split(line, ": ")
			respHeader.Add(HTTPHeader(headerLine[0]), headerLine[1])
		}
	}

	resp.Header = respHeader

	contentLength, err := strconv.Atoi(resp.Header.Get(ContentLength))
	if err == nil && contentLength > 0 {
		respBody := ""
		for scanner.Scan() {
			respBody += scanner.Text()
		}

		resp.Body = respBody
	}

	return resp, nil
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
