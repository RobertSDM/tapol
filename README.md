# **Tapol**

## About

Tapol is a HTTP Client library made for the programming language GO \
It uses the net package sockets to do the requests, without any use of the `net/http` package

## Usage

```go
    url := "https://github.com"
    header := &tapol.Header{}

    resp, err := tapol.Get(url, header)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println(resp.Status)
    fmt.Println(resp.Header)

    b := make([]byte, 1024)

    for {
        n, err := resp.Body.Read(b)
        if err == io.EOF {
            break
        }

        if err != nil {
            panic(err)
        }

        fmt.Print(string(b[:n]))
    }
```

## Features

-   [x] Requests to HTTP and HTTPs protocols
-   [x] Streaming responses
-   [x] Usable version (the v0.1.2 version is usable)
-   [x] Support for `Transfer-Encoding: chunked` header
-   [x] Support to HTTP/1.1
-   [x] Redirect `3xx` support
-   [ ] Support for cache
-   [ ] Support to `Keep-Alive` and `Connection: keep-alive`
-   [ ] Support to cache
-   [ ] Support to connection pooling
