# Helios

<img align="right" src="helios.svg" alt="Helios Logo" width="100px"/>

![Build Status](https://travis-ci.com/yonasadiel/helios.svg?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/yonasadiel/helios/badge.svg?branch=master)](https://coveralls.io/github/yonasadiel/helios?branch=master)

Helios is experimental web apps backbone written in golang.

## Installation

Helios is made with golang 1.13. Install it by typing:

```sh
$ go get -u github.com/yonasadiel/helios
```

## Usage

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/yonasadiel/helios"
)

type SimpleJSONResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

func SimpleJSONHttpHandler(req helios.Request) {
    req.SendJSON(SimpleJSONResponse{Code: "success", Message: "ok"}, http.StatusOK)
}

func main() {
    http.HandleFunc("/", helios.Handle(SimpleJSONHttpHandler))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Middleware

You can define your own middlewares.

```go
middleware1 := func(f HTTPHandler) HTTPHandler {
    return func(req Request) {
        req.SetContextData("some-data", "some-value")
        req.SetContextData("some-data-overwritten", "old-value")
        f(req)
    }
}
middleware2 := func(f HTTPHandler) HTTPHandler {
    return func(req Request) {
        req.SetContextData("some-data-overwritten", "new-value")
        f(req)
    }
}
handler := func(req Request) {
    req.GetContextData("some-data") // "some-value"
    req.GetContextData("some-data-overwritten") // "new-value"
}
http.HandleFunc("/", WithMiddleware(handler, []Middleware{middleware1, middleware2}))

```

## Future of Helios

What I am going to do with Helios:
- Remove gorilla/sessions. It should be the user's choice how to use sessions.
- Remove gorilla/mux. It should be the user's choice how to parse query.
- As you see, this library is not yet ready for usage. The database is still hardcoded and only support sqlite.
- Multiple database connection. For example, an handler may query using read-only database connection, but other handlers use full access database connection
- At this point, I wonder, what if we just store the map of database, and let user define what is the database.
