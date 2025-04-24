# **Tapol**

<img src="imgs/logo.jpg" alt="logo" width="200" style="margin: 0 auto;display: block"/>
<p  style="text-align: center;display: block;font-size: 1.5rem">HTTP Client</p>

## About

Tapol is a HTTP Client library made for the programming language GO \
It uses the net package sockets to do the requests, without any use of the `net/http` package

## Features

-   [ x ] Requests to HTTP and HTTPs protocols
-   [ x ] Streaming responses
-   [ x ] Usable version (the v0.1.2 version is usable)
-   [ x ] Support for `Transfer-Encoding: chunked` header
-   [ x ] Support to HTTP/1.1
-   [ x ] Redirect `3xx` support
-   [ ] Support to `Keep-Alive` and `Connection: keep-alive`
-   [ ] Support to cache
-   [ ] Support to connection pooling
