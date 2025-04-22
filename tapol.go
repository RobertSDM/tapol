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
		conn, err = connect(fmt.Sprintf("%s:%d", (*url).Host, 80))
	} else {
		conn, err = connectTLS(fmt.Sprintf("%s:%d", (*url).Host, 443))
	}

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func buildHeader(header Header, url *url.URL, body []byte) Header {
	header.change(Host, url.Host)
	header.change(Connection, "close")

	if body != nil {
		header.change(ContentLength, fmt.Sprint(len(body)))
	}

	return header
}

func NewRequest(method string, ogurl string, header *Header, body []byte) (resp Response, err error) {
	url, _ := url.Parse(ogurl)
	if err != nil {
		return Response{}, err
	}

	conn, err := createConn(url)
	if err != nil {
		return Response{}, err
	}
	defer conn.Close()

	request := Request{}
	request.Host = url.Host
	request.Header = buildHeader(*header, url, nil)
	request.Method = method

	request.Path = url.Path
	if url.Path == "" {
		request.Path = "/"
	}

	req := request.Build()
	fmt.Fprintln(conn, req)

	scanner := bufio.NewScanner(conn)

	resp = Response{}
	httpheader := false
	respheader := Header{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		if !httpheader {
			status := strings.Split(line, " ")[1:]
			resp.Status = strings.Join(status, " ")
			code, _ := strconv.Atoi(status[0])
			resp.StatusCode = int16(code)

			if resp.StatusCode == 200 {
				resp.Ok = true
			}

			httpheader = true
		} else {
			headerline := strings.Split(line, ": ")
			respheader.Add(HTTPHeader(headerline[0]), headerline[1])
		}
	}

	resp.Header = respheader

	// Redirecting the request
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		return NewRequest(method, respheader.Get(Location), header, body)
	}

	bodylength, err := strconv.Atoi(resp.Header.Get(ContentLength))
	if err == nil && bodylength > 0 {
		respBody := ""
		for scanner.Scan() {
			respBody += scanner.Text()
		}

		resp.Body = respBody
	}

	return resp, nil
}

func Get(url string, header *Header) (resp Response, err error) {
	return NewRequest("GET", url, header, nil)
}

func Post(url string, header *Header, body []byte) (resp Response, err error) {
	return NewRequest("POST", url, header, body)
}

func Put(url string, header *Header, body []byte) (resp Response, err error) {
	return NewRequest("PUT", url, header, body)
}

func Delete(url string, header *Header) (resp Response, err error) {
	return NewRequest("DELETE", url, header, nil)
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
